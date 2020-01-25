# Go parameters

TOOL_VERSION=0.8.0
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
LDFLAGS=-X github.com/ThilinaManamgoda/password-manager/cmd.Version=$(TOOL_VERSION)

all: clean deps lint unit-test build-linux build-darwin build-windows

build-doc:
		$(GOBUILD) -tags doc  -o gen_doc -v

build:
		$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) -v

unit-test:
		$(GOTEST) -v $(TEST_PKGS)
		rm pkg/passwords/testPasswordDB

lint:
		$(GOGET) -u golang.org/x/lint/golint
		$(GOLINT) $(FMT_PKGS)

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

build-linux:
		env GOOS="linux" GOARCH="amd64" $(GOBUILD) -ldflags "$(LDFLAGS)" -o "target/linux/v$(TOOL_VERSION)/$(BINARY_NAME)" -v

build-windows:
		env GOOS="windows" GOARCH="amd64" $(GOBUILD) -ldflags "$(LDFLAGS)" -o "target/windows/v$(TOOL_VERSION)/$(BINARY_NAME).exe" -v

build-darwin:
		env GOOS="darwin" GOARCH="amd64" $(GOBUILD) -ldflags "$(LDFLAGS)" -o "target/darwin/v$(TOOL_VERSION)/$(BINARY_NAME)" -v
