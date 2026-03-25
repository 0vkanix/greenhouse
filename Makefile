.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed 'e'

## test-full: run all tests with coverage and no caching
.PHONY: test-full
test-full:
	go test -v -count=1 -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

.PHONY: test
test:
	go test ./cmd/api/ -coverprofile=coverage.out

.PHONY: build-api
build-api:
	docker build -t greenlight-api -f build/package/Dockerfile .

.PHONY: up
up:
	docker compose up --build

.PHONY: down
down:
	docker compose down

.PHONY: reset
reset:
	docker compose down -v
	docker compose up --build

.PHONY: migrate-new
migrate-new:
	@read -p "Migration name: " name; \
	migrate create -seq -ext=.sql -dir=./migrations $$name
