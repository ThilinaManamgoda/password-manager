# Go parameters

TOOL_VERSION=0.9.1
GENERATE_DOC=false
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=password-manager
DEP=dep
GOLINT=golint
GOFMT=$(GOCMD) fmt

TEST_PKGS=./pkg/...
FMT_PKGS=./cmd/... ./pkg/...
LDFLAGS=-X github.com/ThilinaManamgoda/password-manager/cmd.Version=$(TOOL_VERSION) -X github.com/ThilinaManamgoda/password-manager/cmd.IsGenerateDoc=$(GENERATE_DOC)

all: clean deps lint unit-test build

build:
		$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) -v

unit-test:
		$(GOTEST) -v $(TEST_PKGS)

lint:
		$(GOGET) -u golang.org/x/lint/golint
		$(GOLINT) $(PKGS)

clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)

fmt:
		$(GOFMT) $(FMT_PKGS)

run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)

deps:
		$(DEP) ensure