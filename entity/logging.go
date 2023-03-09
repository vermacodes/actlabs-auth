package entity

import "github.com/gin-gonic/gin"

type OperationRecord struct {
	PartitionKey    string `json:"PartitionKey"`
	RowKey          string `json:"RowKey"`
	UserPrincipal   string `json:"UserPrincipal"`
	OperationId     string `json:"OperationId"`
	OperationStatus string `json:"OperationStatus"`
	OperationType   string `json:"OperationType"`
	LabId           string `json:"LabId"`
	LabName         string `json:"LabName"`
	LabType         string `json:"LabType"`
}

type Operation struct {
	OperationId     string `json:"OperationId"`
	OperationStatus string `json:"OperationStatus"`
	OperationType   string `json:"OperationType"`
	LabId           string `json:"LabId"`
	LabName         string `json:"LabName"`
	LabType         string `json:"LabType"`
}

type LoggingService interface {
	OperationRecord(Operation, string) error
}

type LoggingRepository interface {
	OperationRecord(Operation, string) error
}

type LoggingHandler interface {
	OperationRecord(c *gin.Context)
}
