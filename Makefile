GOCMD		= go
GOLINT		= golangci-lint
SRC         = ./cmd/service

build:
	$(GOCMD) build -o worker_pool $(SRC)

test:
	$(GOCMD) test -v -race ./...

lint:
	$(GOLINT) run -v
