package service

import (
	"encoding/json"
	"errors"
	"strings"

	"actlabs-auth/entity"
	"actlabs-auth/helper"

	"golang.org/x/exp/slog"
)

type assignmentService struct {
	assignmentRepository entity.AssignmentRepository
	labService           entity.LabService
}

func NewAssignmentService(assignmentRepository entity.AssignmentRepository, labService entity.LabService) entity.AssignmentService {
	return &assignmentService{
		assignmentRepository: assignmentRepository,
		labService:           labService,
	}
}

func (a *assignmentService) GetAssignments() ([]entity.Assigment, error) {
	assignments := []entity.Assigment{}

	ar, err := a.assignmentRepository.GetEnumerationResults()
	if err != nil {
		slog.Error("not able to list assignments", err)
		return assignments, err
	}

	for _, element := range ar.Blobs.Blob {
		assignment, err := a.assignmentRepository.GetAssignment(element.Name)
		if err != nil {
			slog.Error("not able to get assignmet "+assignment.Id, err)
			continue
		}
		assignments = append(assignments, assignment)
	}

	return assignments, nil
}

func (a *assignmentService) GetMyAssignments(userPrincipal string) ([]entity.LabType, error) {
	assignedLabs := []entity.LabType{}

	assignments, err := a.GetAssignments()
	if err != nil {
		slog.Error("not able to get assignments", err)
		return assignedLabs, err
	}

	labs, err := a.labService.GetPublicLabs("labexercises")
	if err != nil {
		slog.Error("not able to get lab exercises", err)
		return assignedLabs, err
	}

	for _, assignment := range assignments {
		slog.Info("Assignment ID : " + assignment.Id)
		for _, lab := range labs {
			slog.Info("Lab ID : " + lab.Name)
			if assignment.LabId == lab.Id {
				if assignment.User == userPrincipal {
					lab.ExtendScript = "redacted"
					assignedLabs = append(assignedLabs, lab)
					break
				}
			}
		}
	}

	return assignedLabs, nil
}

func (a *assignmentService) CreateAssignment(assignment entity.Assigment) error {
	// Gnerate Assignmetn ID
	if assignment.Id == "" {
		assignment.Id = helper.Generate(20)
	}

	if !strings.Contains("@microsoft.com", assignment.User) {
		assignment.User = assignment.User + "@microsoft.com"
	}

	// Validate User ID
	valid, err := a.assignmentRepository.ValidateUser(assignment.User)
	if err != nil {
		slog.Error("not able to validate user id", err)
	}

	if !valid {
		err := errors.New("user id is not valid")
		slog.Error("user id is not valid", err)
		return err
	}

	assignments, err := a.GetAssignments()
	if err != nil {
		slog.Error("not able to list existing assingments", err)
		return err
	}

	for _, element := range assignments {
		if element.User == assignment.User && element.LabId == assignment.LabId {
			slog.Info("assignment already exits")
			return nil
		}
	}

	val, err := json.Marshal(assignment)
	if err != nil {
		slog.Error("not able to convert assignment object to string", err)
		return err
	}

	if err := a.assignmentRepository.CreateAssignment(assignment.Id, string(val)); err != nil {
		slog.Error("not able to create assignment", err)
		return err
	}

	return nil
}

func (a *assignmentService) DeleteAssignment(assignment entity.Assigment) error {
	if err := a.assignmentRepository.DeleteAssignment(assignment.Id); err != nil {
		slog.Error("not able to delete assingment with id "+assignment.Id, err)
		return err
	}
	return nil
}
