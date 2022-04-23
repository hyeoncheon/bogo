
LINTERS			= gocognit,gocyclo,misspell
MORE_LINTERS		= errorlint,goerr113,forcetypeassert
DISABLED_LINTERS	= godot,godox,wsl,varnamelen,nlreturn,testpackage,paralleltest,wrapcheck,exhaustivestruct,gci,gochecknoglobals

default: test lint

lint:
	golangci-lint run -E $(LINTERS)
	golangci-lint run -E $(LINTERS),$(MORE_LINTERS) --tests=false

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
