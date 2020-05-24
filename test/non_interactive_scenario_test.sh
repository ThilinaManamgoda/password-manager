#!/usr/bin/env bash
set -e

echo "Running non interactive scenario tests!!!"

MASTER_PASSWORD="test1234"
NEW_MASTER_PASSWORD="test12345"
VERSION=$1
OS="linux"

for e in $(env | grep PM_)
do
  unset "$(echo "$e" | awk '{split($0,a,"="); print a[1]}')"
done


PM_DIRECTORYPATH="$(pwd)/password-manager-test"
export  PM_DIRECTORYPATH
test -d ./password-manager-test && rm -rf ./password-manager-test
test -f ./export-data.csv  && rm ./export-data.csv
test -f ./export-data.html && rm ./export-data.html

../target/${OS}/${VERSION}/password-manager init -m ${MASTER_PASSWORD}
../target/${OS}/${VERSION}/password-manager add test -u test.com -p test12345 -l "fb,gmail" -d "Test description @maanadev" -m ${MASTER_PASSWORD}
../target/${OS}/${VERSION}/password-manager get test -m ${MASTER_PASSWORD} -s | grep "test12345"
echo "===Entering new Master password==="
../target/${OS}/${VERSION}/password-manager change-master-password -m ${MASTER_PASSWORD} -n ${NEW_MASTER_PASSWORD}
echo "===Switching master password to new master password==="
MASTER_PASSWORD=${NEW_MASTER_PASSWORD}
../target/${OS}/${VERSION}/password-manager get test -m ${MASTER_PASSWORD} -s | grep "test12345"

echo "===Importing password from a CSV file==="
../target/${OS}/${VERSION}/password-manager import --csv-file ./mock-data/data.csv -m ${MASTER_PASSWORD}
../target/${OS}/${VERSION}/password-manager get ryendalll@latimes.com -m ${MASTER_PASSWORD} -s | grep "3KVu0V"

echo "===Exporting passwords to a CSV file==="
../target/${OS}/${VERSION}/password-manager export --csv-file ./export-data.csv -m ${MASTER_PASSWORD}
cat ./export-data.csv | grep karmit8@github.io | tr ',' ' '| awk '{print $2}'| grep "Karlik"
cat ./export-data.csv | grep karmit8@github.io | tr ',' ' '| awk '{print $3}'| grep "lfuqz1k"

echo "===Exporting passwords to a HTML file==="
../target/${OS}/${VERSION}/password-manager export --html-file ./export-data.html -m ${MASTER_PASSWORD}
cat ./export-data.html | grep karmit8@github.io

echo "===Remove a password==="
../target/${OS}/${VERSION}/password-manager remove test -m ${MASTER_PASSWORD}

test -d ./password-manager-test && rm -rf ./password-manager-test
test -f ./export-data.csv  && rm ./export-data.csv
test -f ./export-data.html && rm ./export-data.html
