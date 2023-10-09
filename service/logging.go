package service

import (
	"actlabs-auth/entity"

	"golang.org/x/exp/slog"
)

type loggingService struct {
	loggingRepository entity.LoggingRepository
}

func NewLoggingService(loggingRepository entity.LoggingRepository) entity.LoggingService {
	return &loggingService{
		loggingRepository: loggingRepository,
	}
}

func (s *loggingService) OperationRecord(operation entity.Operation, userPrincipal string) error {
	slog.Info("Logging operation: " + operation.OperationId + " for user: " + userPrincipal)
	return s.loggingRepository.OperationRecord(operation, userPrincipal)
}
