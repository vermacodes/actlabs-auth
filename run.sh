#!/bin/bash

export ROOT_DIR=$(pwd)

if [[ "${SAS_TOKEN}" == "" ]]; then
    echo "SAS TOKEN missing"
    exit 1
fi

if [[ "${STORAGE_ACCOUNT_NAME}" == "" ]]; then
    echo "STORAGE ACCOUNT NAME missing"
    exit 1
fi

echo "Storage Account -> ${STORAGE_ACCOUNT_NAME}"

go build -ldflags "-X 'actlabs-auth/entity.SasToken=$SAS_TOKEN' -X 'actlabs-auth/entity.StorageAccountName=$STORAGE_ACCOUNT_NAME'"

redis-cli flushall && export LOG_LEVEL="-4" && export PORT="8882" && ./actlabs-auth