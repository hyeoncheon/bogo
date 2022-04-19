
default: test lint

lint:
	golangci-lint run -E gocognit,gocyclo,misspell --tests=false

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
