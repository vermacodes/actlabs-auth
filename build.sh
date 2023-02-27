#!/bin/bash

# This script starts the web app and the server. Both server and the webapp needs to be exposed to the world outside.
#
# WebApp runs on port 3000
# Server runs on port 8080.

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


docker build -t actlab.azurecr.io/actlabs-auth .

az acr login --name actlab
docker push actlab.azurecr.io/actlabs-auth:latest