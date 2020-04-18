#!/usr/bin/env bash

set -e

VERSION=$1
WORK_DIR=$(dirname "${BASH_SOURCE[0]}")
INSTALLER_CLONE_DIR="$WORK_DIR/password-manager-installers"

WITHOUT_V_VERSION="$(echo "$VERSION" | sed 's/^.//')"
git clone https://github.com/ThilinaManamgoda/password-manager-installers.git "$INSTALLER_CLONE_DIR" -b "$(echo "$VERSION" | sed 's/.$/x/' | sed 's/^.//')"
cp target/linux/$VERSION/password-manager "$INSTALLER_CLONE_DIR"
docker run -it -v "${PWD}/$INSTALLER_CLONE_DIR":/home/ubuntu maanadev/debian-tools-builder:1.0.0 ${WITHOUT_V_VERSION}
cp "$INSTALLER_CLONE_DIR/password-manager_${WITHOUT_V_VERSION}_amd64.deb" target/linux/$VERSION/
rm -rf "$INSTALLER_CLONE_DIR"
