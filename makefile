# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=password-manager
DEP=dep
GOLINT=golint

all: clean deps lint unit-test build


build:
		$(GOBUILD) -o $(BINARY_NAME) -v
unit-test:
		$(GOTEST) -v ./...
lint:
		$(GOGET) -u golang.org/x/lint/golint
		$(GOLINT)  cmd/... pkg/...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
deps:
		$(DEP) ensure