APP_NAME=core-transaction-srv
VERSION=0.0.1

ifndef env
env := local
endif

dev-dependencies: dependencies
	go install github.com/cosmtrek/air/...@latest
	go install github.com/rubenv/sql-migrate/...@latest

dependencies:
	export GO111MODULE=on; \
	export GOPRIVATE=github.com/Neutronpay; \
	go mod tidy

.PHONY: dev
## dev: Run dev server with live reload support
dev:
	air

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

.PHONY: help
all: help
# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

dev-srv:
	go run ./main.go -e ${env}

docker-local:
	docker-compose -f docker-compose.local.yml up -d

