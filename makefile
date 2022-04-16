

lint:
	golangci-lint run -E gocognit,gocyclo,misspell --tests=false

test:
	go test -cover ./...
