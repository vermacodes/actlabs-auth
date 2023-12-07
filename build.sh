#!/bin/bash

# gather input parameters
# -t tag

while getopts ":t:" opt; do
    case $opt in
    t)
        TAG="$OPTARG"
        ;;
    \?)
        echo "Invalid option -$OPTARG" >&2
        ;;
    esac
done

if [ -z "${TAG}" ]; then
    TAG="latest"
fi

echo "TAG = ${TAG}"
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


docker build -t actlab.azurecr.io/actlabs-auth:${TAG} .

rm actlabs-auth

az acr login --name actlab --subscription ACT-CSS-Readiness
docker push actlab.azurecr.io/actlabs-auth:${TAG}

docker tag actlab.azurecr.io/actlabs-auth:${TAG} ashishvermapu/actlabs-auth:${TAG}
docker push ashishvermapu/actlabs-auth:${TAG}