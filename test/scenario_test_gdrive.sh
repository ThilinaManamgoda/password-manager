#!/usr/bin/env bash

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

set -e
MASTER_PASSWORD="test1234"
NEW_MASTER_PASSWORD="test12345"
VERSION="v0.9.3"
GDRIVE_CLIENT_ID=""
GDRIVE_CLIENT_SEC=""

pushd ../
 make build-darwin TOOL_VERSION=${VERSION} GDRIVE_CLIENT_ID=${GDRIVE_CLIENT_ID} GDRIVE_CLIENT_SEC=${GDRIVE_CLIENT_SEC}
popd

for e in $(env | grep PM_)
do
  unset "$(echo $e | awk '{split($0,a,"="); print a[1]}')"
done

PM_DIRECTORYPATH="$(pwd)/password-manager-test"
export PM_DIRECTORYPATH
export PM_STORAGE_FILE_ENABLE=false
export PM_STORAGE_GOOGLEDRIVE_ENABLE=true
export PM_STORAGE_GOOGLEDRIVE_DIRECTORY="password-manager-test"
PM_STORAGE_GOOGLEDRIVE_PASSWORDDBFILE="$(uuidgen)"
export PM_STORAGE_GOOGLEDRIVE_PASSWORDDBFILE
export PM_SELECTLISTSIZE=2

test -d ./password-manager-test && rm -rf ./password-manager-test
test -f ./export-data.csv  && rm ./export-data.csv

../target/darwin/${VERSION}/password-manager init -m ${MASTER_PASSWORD}
../target/darwin/${VERSION}/password-manager add test -u test.com -p test12345 -l "fb,gmail" -d "Test description @maanadev" -m ${MASTER_PASSWORD}
../target/darwin/${VERSION}/password-manager get test -m ${MASTER_PASSWORD} -s | grep "test12345"
echo "===Searching password==="
../target/darwin/${VERSION}/password-manager search test -s -m ${MASTER_PASSWORD} | grep "test12345"
echo "===Searching password with a label==="
../target/darwin/${VERSION}/password-manager search -l "fb" -s -m ${MASTER_PASSWORD} | grep "test12345"

echo "===Entering new Master password==="
../target/darwin/${VERSION}/password-manager change-master-password -m ${MASTER_PASSWORD} -n ${NEW_MASTER_PASSWORD}
echo "===Switching master password to new master password==="
MASTER_PASSWORD=${NEW_MASTER_PASSWORD}
../target/darwin/${VERSION}/password-manager search test -s -m ${MASTER_PASSWORD}

echo "===Importing password from a CSV file==="
../target/darwin/${VERSION}/password-manager import --csv-file ./mock-data/data.csv -m ${MASTER_PASSWORD}
../target/darwin/${VERSION}/password-manager get ryendalll@latimes.com -m ${MASTER_PASSWORD} -s | grep "3KVu0V"

echo "===Exporting passwords to a CSV file==="
../target/darwin/${VERSION}/password-manager export --csv-file ./export-data.csv -m ${MASTER_PASSWORD}

echo "===Remove a password==="
../target/darwin/${VERSION}/password-manager remove test -m ${MASTER_PASSWORD}

test -d ./password-manager-test && rm -rf ./password-manager-test
test -f ./export-data.csv  && rm ./export-data.csv
