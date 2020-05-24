# Copyright © 2020 Thilina Manamgoda
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
FROM golang:1.12.1 AS builder

ENV TOOL_VERSION=v0.9.0
ENV GDRIVE_CLIENT_ID=""
ENV GDRIVE_CLIENT_SEC=""
ENV LDFLAGS="-X 'github.com/ThilinaManamgoda/password-manager/cmd.Version=${TOOL_VERSION}' -X 'github.com/ThilinaManamgoda/password-manager/pkg/storage/googledrive.ClientID=${GDRIVE_CLIENT_ID}' -X 'github.com/ThilinaManamgoda/password-manager/pkg/storage/googledrive.ClientSecret=${GDRIVE_CLIENT_SEC}'"

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.5.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

COPY Gopkg.toml Gopkg.lock main.go $GOPATH/src/github.com/ThilinaManamgoda/password-manager/
COPY cmd/ $GOPATH/src/github.com/ThilinaManamgoda/password-manager/cmd/
COPY pkg/ $GOPATH/src/github.com/ThilinaManamgoda/password-manager/pkg/
COPY vendor/ $GOPATH/src/github.com/ThilinaManamgoda/password-manager/vendor/

WORKDIR $GOPATH/src/github.com/ThilinaManamgoda/password-manager
RUN dep ensure

RUN GOOS="linux" GOARCH="amd64" go build -ldflags "${LDFLAGS}" -a -installsuffix nocgo -o /password-manager $GOPATH/src/github.com/ThilinaManamgoda/password-manager

FROM ubuntu:18.04
RUN apt update && apt install -y xclip
COPY --from=builder /password-manager ./

