#!/bin/bash

# This script is for local testing. It starts both server and UI in one go.

# gather input parameters
# if flag -d is set the set LOG_LEVEL to -4 else to 0

while getopts ":d" opt; do
    case $opt in
    d)
        LOG_LEVEL="-4"
        ;;
    \?)
        echo "Invalid option -$OPTARG" >&2
        ;;
    esac
done

if [ -z "${LOG_LEVEL}" ]; then
    LOG_LEVEL="0"
fi

echo "LOG_LEVEL = ${LOG_LEVEL}"

export ROOT_DIR=$(pwd)

if [[ "${SAS_TOKEN}" == "" ]]; then
    echo "SAS TOKEN missing"
    exit 1
fi

if [[ "${STORAGE_ACCOUNT_NAME}" == "" ]]; then
    echo "STORAGE ACCOUNT NAME missing"
    exit 1
fi

# Service principal opject ID
if [[ "${AUTH_TOKEN_AUD}" == "" ]]; then
    echo "AUTH_TOKEN_AUD missing"
    exit 1
fi

# "https://login.microsoftonline.com/{tenant-id}/v2.0"
if [[ "${AUTH_TOKEN_ISS}" == "" ]]; then
    echo "AUTH_TOKEN_ISS missing"
    exit 1
fi

echo "Storage Account -> ${STORAGE_ACCOUNT_NAME}"

# Remove existing binary.
rm actlabs-auth

go build -ldflags "-X 'actlabs-auth/entity.SasToken=$SAS_TOKEN' -X 'actlabs-auth/entity.StorageAccountName=$STORAGE_ACCOUNT_NAME'"

redis-cli flushall && export LOG_LEVEL="${LOG_LEVEL}" && export PORT="8882" && ./actlabs-auth