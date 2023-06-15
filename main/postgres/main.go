package main

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/tiendc/sqlboiler-extensions-demo/db/models"
)

const (
	connStr = "postgres://admin:password@localhost:5432/mydb?sslmode=disable"
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// boil.DebugMode = true
	boil.SetDB(db)

	// Optional: set a specific value for DefaultPageSize
	// models.DefaultPageSize = 11111

	return db, nil
}

func main() {
	ctx := context.Background()

	db, err := initDB()
	if err != nil {
		print("error opening db connection", err)
		return
	}
	defer func() {
		_ = db.Close()
	}()

	// We have 3 tables: users, orders, shipments
	// Relationships: users <--1:M--> orders <--1:1--> shipments

	// Clear all current data
	_, err = models.Shipments().DeleteAll(ctx, db)
	_, err = models.Orders().DeleteAll(ctx, db)
	_, err = models.Users().DeleteAll(ctx, db)

	// Create 1000 users, 100000 orders and shipments
	numUsers := 1000
	numOrders := 100000

	allUsers := models.UserSlice{}
	for i := 1; i <= numUsers; i++ {
		allUsers = append(allUsers, &models.User{
			ID:   i,
			Name: fmt.Sprintf("user-%d", i),
		})
	}

	allOrders := models.OrderSlice{}
	allShipments := models.ShipmentSlice{}
	for i := 1; i <= numOrders; i++ {
		orderID := i
		allOrders = append(allOrders, &models.Order{
			ID:     orderID,
			Name:   fmt.Sprintf("order-%d", i),
			UserID: rand.Intn(numUsers) + 1,
		})
		allShipments = append(allShipments, &models.Shipment{
			ID:      i,
			Status:  models.ShipmentStatusNOT_SHIPPED,
			OrderID: orderID,
		})
	}

	start := time.Now()

	// Test InsertAll
	print("\n>>> Test InsertAll 1,000 users")
	_, err = allUsers.InsertAll(ctx, db, boil.Infer())
	print(">>> Expect OK. Got error:", err)
	expectPassed(err)

	print("\n>>> Test InsertAll 100,000 orders/shipments")
	_, err = allOrders.InsertAll(ctx, db, boil.Infer())
	print(">>> Expect Error \"PostgreSQL only supports 65535 parameters\". Got error:", err)
	expectFailed(err)

	print("\n>>> Test InsertAllByPage 100,000 orders")
	_, err = allOrders.InsertAllByPage(ctx, db, boil.Infer())
	print(">>> Expect OK. Got error:", err)
	expectPassed(err)
	print("\n>>> Test InsertAllByPage 100,000 shipments")
	_, err = allShipments.InsertAllByPage(ctx, db, boil.Infer())
	print(">>> Expect OK. Got error:", err)
	expectPassed(err)

	// Test UpsertAll
	print("\n>>> Test UpsertAll: Add new user-1001 and change user-1 name")
	allUsers[0].Name = "user-1-edited"
	allUsers = append(allUsers, &models.User{
		ID:   1001,
		Name: "user-1001-new",
	})
	_, err = allUsers.UpsertAll(ctx, db, true, nil, boil.Infer(), boil.Infer())
	print(">>> Expect OK. Got error:", err)

	// Verify the change
	user1, _ := models.FindUser(ctx, db, 1)
	user1001, _ := models.FindUser(ctx, db, 1001)
	print(">>> Verification: User-1 name:", user1.Name, ", User-1001 name:", user1001.Name)
	expectPassed(err)

	// Test eager loading
	print("\n>>> Test EagerLoading orders from users")
	allUsers, err = models.Users(
		qm.Load(models.UserRels.Orders),
	).All(ctx, db)
	print(">>> Expect OK. Got error:", err)
	expectPassed(err)

	// Eager loading of 2 levels will fail at the 2nd step.
	// Step 1: Eager loading 100,000 orders for 1,000 users -> OK
	// Step 2: Eager loading 100,000 shipments for 100,000 orders -> FAILED
	print("\n>>> Test EagerLoading orders and shipments from users")
	allUsers, err = models.Users(
		qm.Load(qm.Rels(models.UserRels.Orders, models.OrderRels.Shipments)),
	).All(ctx, db)
	print(">>> Expect Error \"PostgreSQL only supports 65535 parameters\". Got error:", err)
	expectFailed(err)

	// Test eager loading by page
	print("\n>>> Test EagerLoading orders and shipments from users by page")
	allUsers, err = models.Users(
		qm.Load(models.UserRels.Orders),
	).All(ctx, db)
	// Eager loading shipments is done in a separated call
	err = allUsers.GetLoadedOrders().LoadShipmentsByPage(ctx, db)
	print(">>> Expect OK. Got error:", err)
	expectPassed(err)

	// Test DeleteAll
	print("\n>>> Test DeleteAll 100,000 shipments")
	_, err = allShipments.DeleteAll(ctx, db)
	print(">>> Expect Error \"PostgreSQL only supports 65535 parameters\". Got error:", err)
	expectFailed(err)

	print("\n>>> Test DeleteAllByPage 100,000 shipments")
	_, err = allShipments.DeleteAllByPage(ctx, db)
	print(">>> Expect OK. Got error:", err)
	expectPassed(err)

	// Test DeleteAll shipments of 70,000 orders
	orderIDs := []interface{}{}
	for i := 1000; i < 80000; i++ {
		orderIDs = append(orderIDs, uint(i))
	}
	print("\n>>> Test DeleteAll shipments of 70,000 orders")
	_, err = models.Shipments(
		qm.WhereIn("order_id IN ?", orderIDs...), // too many IDs
	).DeleteAll(ctx, db)
	print(">>> Expect Error \"PostgreSQL only supports 65535 parameters\". Got error:", err)
	expectFailed(err)

	print("\n>>> Test DeleteAll shipments of 70,000 orders by page")
	for _, chunk := range models.SplitInChunks(orderIDs) {
		_, err = models.Shipments(
			qm.WhereIn("order_id IN ?", chunk...),
		).DeleteAll(ctx, db)
		if err != nil {
			break
		}
	}
	print(">>> Expect OK. Got error:", err)
	expectPassed(err)

	print("\nTOTAL TEST DURATION:", time.Since(start))
}

func print(args ...interface{}) {
	fmt.Println(args...)
}

func expectPassed(err error) {
	if err == nil {
		print(">>> PASSED")
	} else {
		print(">>> FAILED")
	}
}

func expectFailed(err error) {
	if err != nil {
		print(">>> PASSED")
	} else {
		print(">>> FAILED")
	}
}
