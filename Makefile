GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
BINARY_SERVER := bin/server
BINARY_SEEDER := bin/seeder

all: test build
build:
	$(GOBUILD) -o ./$(BINARY_SERVER) -v ./cmd/server/main.go
	$(GOBUILD) -o ./$(BINARY_SEEDER) -v ./cmd/seeder/main.go
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_SERVER)
server:
	make build
	./$(BINARY_SERVER)
seeder:
	make build
	./$(BINARY_SEEDER)
dep:
	dep ensure

