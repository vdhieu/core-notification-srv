APP_NAME=core-transaction-srv
VERSION=0.0.1

ifndef env
env := local
endif

.IGNORE:

dev-dependencies: dependencies
	go install github.com/cosmtrek/air/...@latest
	go install github.com/rubenv/sql-migrate/...@latest

dependencies:
	export GO111MODULE=on; \
	export GOPRIVATE=github.com/Neutronpay; \
	go mod tidy

.PHONY: migrate-new
## migrate-new: Create new manual migration file
migrate-new:
	atlas migrate new \
	${name}

.PHONY: migrate-hash
migrate-hash:
## migrate-hash: Create checksum hash for manual migration file
	atlas migrate hash

.PHONY: migrate-up
## migrate-up: Run migration
migrate-up:
	atlas migrate apply --env=${env}

.PHONY: migrate-dry-run
migrate-dry-run:
	atlas migrate apply --dry-run --env=${env}

.PHONY: migrate-down
## migrate-down: Revert migration -- need testing with atlas
migrate-down:
	atlas schema apply --env=${env} --to "file;//migrations?version=${to-version}"

DEFAULT_QUERY := CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";
.PHONY: atlas-migration
## atlas-migration: spin up test db for atlas migration
atlas-migration:
	docker run --name atlas-migration -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=test -p 0.0.0.0:55432:5432/tcp -d postgres:13.1
	sleep 1
	docker exec -i atlas-migration psql -U postgres -d test -c "${DEFAULT_QUERY}"

.PHONY: atlas-clean
## atlas-clean: clean up test db
atlas-clean:
	docker stop atlas-migration
	docker rm -f atlas-migration

.PHONY: migrate-generate
## migrate-generate: Generate migration file using gorm AutoMigrate and atlas migrate diff
migrate-generate: .IGNORE atlas-migration # ignore error if migration failed in order to make sure atlas-migration container is stopped
	go run ./cmd/migration/main.go -conn 'postgres://postgres:postgres@localhost:55432/test?sslmode=disable'
	set -x; \
	atlas migrate diff \
  --dir file://migrations \
  --dev-url docker://postgres \
  --to 'postgres://postgres:postgres@localhost:55432/test?sslmode=disable' \
  ${name}
	$(MAKE) atlas-clean # invoke clean target to remove atlas-migration container. Invoke with $(MAKE) to keep the same context


.PHONY: help
all: help
# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

dev:
	go run ./cmd/server/webbook/main.go -e ${env}

docker-local:
	docker-compose -f docker-compose.local.yml up -d

