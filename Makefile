PATH := $(PATH):$(GOPATH)/bin
BOIL_VER := v4.15.0
BOIL_EXT_VER := v0.7.2

prepare:
	@go get -u github.com/tiendc/sqlboiler-extensions@$(BOIL_EXT_VER)
	@go install github.com/volatiletech/sqlboiler/v4@$(BOIL_VER)
	@go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@$(BOIL_VER)
	@go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@$(BOIL_VER)
	@go install github.com/glerchundi/sqlboiler-crdb/v4@latest

gen-models-mysql:
	@sqlboiler mysql -c db/sqlboiler.toml \
 		--templates $(GOPATH)/pkg/mod/github.com/volatiletech/sqlboiler/v4@$(BOIL_VER)/templates/main \
 		--templates $(GOPATH)/pkg/mod/github.com/tiendc/sqlboiler-extensions@$(BOIL_EXT_VER)/templates/boilv4/mysql

gen-models-postgres:
	@sqlboiler psql -c db/sqlboiler.toml \
 		--templates $(GOPATH)/pkg/mod/github.com/volatiletech/sqlboiler/v4@$(BOIL_VER)/templates/main \
 		--templates $(GOPATH)/pkg/mod/github.com/tiendc/sqlboiler-extensions@$(BOIL_EXT_VER)/templates/boilv4/postgres

# for cockroachdb and postgres we can use the same templates
gen-models-crdb:
	@sqlboiler crdb -c db/sqlboiler.toml \
 		--templates $(GOPATH)/pkg/mod/github.com/volatiletech/sqlboiler/v4@$(BOIL_VER)/templates/main \
 		--templates $(GOPATH)/pkg/mod/github.com/tiendc/sqlboiler-extensions@$(BOIL_EXT_VER)/templates/boilv4/postgres


run-test-mysql:
	@go run $(PWD)/main/mysql/...

run-test-postgres:
	@go run $(PWD)/main/postgres/...

run-test-crdb:
	@go run $(PWD)/main/crdb/...
