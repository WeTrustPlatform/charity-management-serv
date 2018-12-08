GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
BUILT := bin
BINARY_SERVER := $(BUILT)/cms-server
BINARY_SEEDER := $(BUILT)/cms-seeder

all: test lint build
build-server:
	$(GOBUILD) -o ./$(BINARY_SERVER) ./cmd/server/main.go
build-seeder:
	$(GOBUILD) -o ./$(BINARY_SEEDER) ./cmd/seeder/main.go
build: build-server build-seeder
test:
	$(GOTEST) -coverprofile cp.out -v ./...
clean:
	$(GOCLEAN)
	rm -r $(BUILT)
dev: build-server
	./live_reload.sh $(BINARY_SERVER) "make build-server"
launch:
	./$(BINARY_SERVER)
seeder: build-seeder
	./$(BINARY_SEEDER)
dep:
	dep ensure
lint:
	@echo "Step: gometalinter"
	@gometalinter --vendor --exclude "defer" ./...
	@echo "Step: gofmt (simplify)"
	@! gofmt -s -l . 2>&1 | grep -v vendor
	@echo "Step: goimports"
	@! goimports -l . | grep -v vendor
