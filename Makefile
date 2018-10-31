GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
BUILT := bin
BINARY_SERVER := ${BUILT}/cms-server
BINARY_SEEDER := ${BUILT}/cms-seeder

all: test lint build
build-server:
	$(GOBUILD) -o ./$(BINARY_SERVER) -v ./cmd/server/main.go
build-seeder:
	$(GOBUILD) -o ./$(BINARY_SEEDER) -v ./cmd/seeder/main.go
build: build-server build-seeder
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	@rm -f $(BINARY_SERVER)
server: build-server
	./$(BINARY_SERVER)
seeder: build-seeder
	./$(BINARY_SEEDER)
dep:
	dep ensure
lint:
	@echo "gometalinter"
	@gometalinter --vendor ./...
	@echo "gofmt (simplify)"
	@! gofmt -s -l . 2>&1 | grep -v vendor
	@echo "goimports"
	@! goimports -l . | grep -v vendor
