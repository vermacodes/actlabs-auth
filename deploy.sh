  # create container app env.
  az containerapp env create --name ${CONTAINERAPP_NAME}-env \
    --resource-group ${RESOURCE_GROUP} \
    --logs-destination none

  # create container app
  az containerapp create --name ${CONTAINERAPP_NAME} \
    --resource-group ${RESOURCE_GROUP} \
    --environment ${CONTAINERAPP_NAME}-env \
    --allow-insecure false \
    --image ${DOCKER_IMAGE} \
    --ingress 'external' \
    --min-replicas 1 \
    --max-replicas 1 \
    --target-port 80 \
    --env-vars "ARM_CLIENT_ID=$ARM_CLIENT_ID" "ARM_CLIENT_SECRET=secretref:arm-client-secret" "ARM_SUBSCRIPTION_ID=$ARM_SUBSCRIPTION_ID" "ARM_TENANT_ID=$ARM_TENANT_ID" "ARM_USER_PRINCIPAL_NAME=$ARM_USER_PRINCIPAL_NAME" "LOG_LEVEL=$LOG_LEVEL" \
    --secrets "arm-client-secret=$ARM_CLIENT_SECRET"