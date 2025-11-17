-include .env 

BINARY_NAME = app
DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(POSTGRES_DB)?sslmode=${DB_SSLMODE}

.PHONY: all build clean 

run_products_service:
	go run ./cmd/main.go


docker-compose-infra-up:
	@echo Starting infrastructure docker-compose up
	docker-compose -f docker-compose-infra.yml up -d

# https://github.com/golang-migrate/migrate/blob/856ea12df9d230b0145e23d951b7dbd6b86621cb/database/postgres/TUTORIAL.md
# https://github.com/golang-migrate/migrate/blob/856ea12df9d230b0145e23d951b7dbd6b86621cb/GETTING_STARTED.md
# https://github.com/golang-migrate/migrate/blob/856ea12df9d230b0145e23d951b7dbd6b86621cb/MIGRATIONS.md
# https://github.com/golang-migrate/migrate/tree/856ea12df9d230b0145e23d951b7dbd6b86621cb/cmd/migrate#usage
go-migrate:
	@./scripts/go-migrate.sh -p ./migrations/ -c up -o "$(DB_URL)"	

go-migrate-down:
	@./scripts/go-migrate.sh -p ./migrations/ -c down -o "$(DB_URL)"

# --- Main scripts ---
build:
	mkdir -p ./.bin
	go build -o ./.bin/${BINARY_NAME} ./cmd

run: build
	./.bin/${BINARY_NAME}

clean:
	rm -rf .bin

test:
	go test ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy
