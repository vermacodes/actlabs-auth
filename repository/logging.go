package repository

import (
	"actlabs-auth/entity"
	"context"
	"encoding/json"

	"golang.org/x/exp/slog"
)

type loggingRepository struct{}

func NewLoggingRepository() entity.LoggingRepository {
	return &loggingRepository{}
}

func (r *loggingRepository) OperationRecord(operation entity.Operation, userPrincipal string) error {

	// getServiceClient() can be found in repsitory/auth.go
	serviceClient := getServiceClient().NewClient("Operations")

	operationRecord := entity.OperationRecord{
		PartitionKey:    userPrincipal,
		RowKey:          operation.OperationId + "-" + operation.OperationType + "-" + operation.OperationStatus,
		UserPrincipal:   userPrincipal,
		OperationId:     operation.OperationId,
		OperationStatus: operation.OperationStatus,
		OperationType:   operation.OperationType,
		LabId:           operation.LabId,
		LabName:         operation.LabName,
		LabType:         operation.LabType,
	}

	marshalledOperationRecord, err := json.Marshal(operationRecord)
	if err != nil {
		slog.Error("Error marshalling operation record: ", err)
		return err
	}

	slog.Info("Adding Operation Record: " + string(marshalledOperationRecord))

	_, err = serviceClient.AddEntity(context.TODO(), marshalledOperationRecord, nil)
	if err != nil {
		slog.Error("Error adding entity: ", err)
		return err
	}

	return nil
}
