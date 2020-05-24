# Copyright © 2019 Thilina Manamgoda
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

TOOL_VERSION=v0.9.3
DATABASE_VERSION="v0.9.1"
GDRIVE_CLIENT_ID=""
GDRIVE_CLIENT_SEC=""
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
LDFLAGS=-X 'github.com/ThilinaManamgoda/password-manager/cmd.Version=$(TOOL_VERSION)' -X 'github.com/ThilinaManamgoda/password-manager/pkg/storage/googledrive.ClientID=$(GDRIVE_CLIENT_ID)' -X 'github.com/ThilinaManamgoda/password-manager/pkg/storage/googledrive.ClientSecret=$(GDRIVE_CLIENT_SEC)' -X 'github.com/ThilinaManamgoda/password-manager/pkg/passwords.DatabaseVersion=$(DATABASE_VERSION)'

all: clean deps lint ineffassign unit-test build-linux build-darwin build-windows non-interactive-sce-test

build-doc:
		$(GOBUILD) -tags doc  -o gen_doc -v

build:
		$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) -v

unit-test:
		$(GOTEST) -v $(TEST_PKGS)
		# This directory is created with the password repo unit tests.
		rm -rf pkg/passwords/password-manager-tmp

lint:
		$(GOGET) -u golang.org/x/lint/golint
		$(GOLINT) $(FMT_PKGS)

clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)

fmt:
		$(GOFMT) $(FMT_PKGS)

run:
		$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)

deps:
		$(DEP) ensure

build-linux:
		env GOOS="linux" GOARCH="amd64" $(GOBUILD) -ldflags "$(LDFLAGS)" -o "target/linux/$(TOOL_VERSION)/$(BINARY_NAME)" -v

build-windows:
		env GOOS="windows" GOARCH="amd64" $(GOBUILD) -ldflags "$(LDFLAGS)" -o "target/windows/$(TOOL_VERSION)/$(BINARY_NAME).exe" -v

build-darwin:
		env GOOS="darwin" GOARCH="amd64" $(GOBUILD) -ldflags "$(LDFLAGS)" -o "target/darwin/$(TOOL_VERSION)/$(BINARY_NAME)" -v

build-darwin-tar:
		cd "target/darwin/$(TOOL_VERSION)"; tar -cvf password-manager-darwin-$(TOOL_VERSION).tar.gz $(BINARY_NAME);

build-linux-deb:
		./resources/release-scripts/build-debian-package.sh $(TOOL_VERSION)

release-artifcats: all build-darwin-tar build-linux-deb

build-docker-image:
	    docker build -t password-manager:$(TOOL_VERSION) ./ --file test/dockerfile/Dockerfile

non-interactive-sce-test:
		cd "./test";./non_interactive_scenario_test.sh $(TOOL_VERSION);

ineffassign:
		$(GOGET) -u github.com/gordonklaus/ineffassign
		ineffassign main.go cmd/* pkg/*
