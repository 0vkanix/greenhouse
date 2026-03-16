.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed 'e'

## test-full: run all tests with coverage and no caching
.PHONY: test-full
test-full:
	go test -v -count=1 -coverprofile=coverage.out ./cmd/api/...
	go tool cover -func=coverage.out
