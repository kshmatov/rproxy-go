MAIN=proxy
DSTSRV=dstsrv
CLIENT=tclient
ifeq ($(OS),Windows_NT)
    BIN=./bin/$(MAIN).exe
	SRVBIN=./bin/$(DSTSRV).exe
	CLIBIN=./bin/$(CLIENT).exe
else
    BIN=./bin/$(MAIN)
	SRVBIN=./bin/$(DSTSRV)
	CLIBIN=./bin/$(CLIENT)
endif

.PHONY: ALL
all: deps vet build-dev test-dev

.PHONY: fmt
fmt:
	goimports -w .

.PHONY: vet
vet: fmt
	golangci-lint run  ./...

.PHONY: test-dev
test-dev:
	go test ./...

.PHONY: build-dev
build-proxy:
	go build -o $(BIN) ./cmd/proxy/main.go

build-dst:
	go build -o $(SRVBIN) ./cmd/dst-srv/main.go

build-cli:
	go build -o $(CLIBIN) ./cmd/client/main.go

.PHONY: deps
deps:
	go mod tidy

init:
	go mod init

run:  build-proxy
	$(BIN) -l d  -h :9999 127.0.0.1:8000

run-dst: build-dst
	$(SRVBIN) -n DEF &
	$(SRVBIN) -n DEF -p 8082 &

stop-dst:
	killall $(SRVBIN)

run-cli: build-cli
	$(CLIBIN) -d 127.0.0.1:9999/test

run-direct: build-cli
	$(CLIBIN) -d 127.0.0.1:8081/test