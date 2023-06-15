PATH := $(PATH):$(GOPATH)/bin
BOIL_VER := v4.13.0

prepare-mysql:
	@go install github.com/volatiletech/sqlboiler/v4@$(BOIL_VER)
	@go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@$(BOIL_VER)

gen-models-mysql:
	@sqlboiler mysql -c db/sqlboiler.toml \
 		--templates $(GOPATH)/pkg/mod/github.com/volatiletech/sqlboiler/v4@$(BOIL_VER)/templates/main \
 		--templates ./db/extensions/templates/boilv4/mysql

run-test-mysql:
	@go run $(PWD)/main/mysql/...


prepare-postgres:
	@go install github.com/volatiletech/sqlboiler/v4@$(BOIL_VER)
	@go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@$(BOIL_VER)

gen-models-postgres:
	@sqlboiler psql -c db/sqlboiler.toml \
 		--templates $(GOPATH)/pkg/mod/github.com/volatiletech/sqlboiler/v4@$(BOIL_VER)/templates/main \
 		--templates ./db/extensions/templates/boilv4/postgres

run-test-postgres:
	@go run $(PWD)/main/postgres/...


prepare-crdb:
	@go install github.com/volatiletech/sqlboiler/v4@$(BOIL_VER)
	@go install github.com/glerchundi/sqlboiler-crdb/v4@latest

# cockroachdb and postgres can use the same templates
gen-models-crdb:
	@sqlboiler crdb -c db/sqlboiler.toml \
 		--templates $(GOPATH)/pkg/mod/github.com/volatiletech/sqlboiler/v4@$(BOIL_VER)/templates/main \
 		--templates ./db/extensions/templates/boilv4/postgres

run-test-crdb:
	@go run $(PWD)/main/crdb/...
