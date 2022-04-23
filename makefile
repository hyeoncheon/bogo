
LINT_TARGET		?= ./...
LINTERS			= gocognit,gocyclo,misspell
MORE_LINTERS		= errorlint,goerr113,forcetypeassert,gosec

DISABLED_LINTERS_NEVER	= godox,forbidigo,varnamelen,paralleltest,testpackage
DISABLED_LINTERS	= $(DISABLED_LINTERS_NEVER),exhaustivestruct,ireturn,gci,gochecknoglobals,wrapcheck
DISABLED_LINTERS_TEST	= $(DISABLED_LINTERS),nolintlint,wsl,goerr113,forcetypeassert,errcheck,goconst,scopelint

default: test lint

lint:
	golangci-lint run -E $(LINTERS)
	golangci-lint run -E $(LINTERS),$(MORE_LINTERS) --tests=false

lint-hard:
	go test -cover $(LINT_TARGET)
	golangci-lint run --exclude-use-default=false --enable-all \
		-D $(DISABLED_LINTERS) --tests=false $(LINT_TARGET)
	golangci-lint run --exclude-use-default=false --enable-all \
		-D $(DISABLED_LINTERS_TEST) $(LINT_TARGET)

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
