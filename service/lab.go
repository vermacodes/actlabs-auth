package service

import (
	"encoding/json"

	"actlabs-auth/entity"
	"actlabs-auth/helper"

	"golang.org/x/exp/slog"
)

type labService struct {
	labRepository entity.LabRepository
}

func NewLabService(repo entity.LabRepository) entity.LabService {
	return &labService{
		labRepository: repo,
	}
}

func (l *labService) GetPublicLabs(typeOfLab string) ([]entity.LabType, error) {
	labs := []entity.LabType{}

	er, err := l.labRepository.GetEnumerationResults(typeOfLab)
	if err != nil {
		slog.Error("Not able to get list of blobs", err)
		return labs, err
	}

	for _, element := range er.Blobs.Blob {
		lab, err := l.labRepository.GetLab(element.Name, typeOfLab)
		if err != nil {
			slog.Error("not able to get blob from given url", err)
			continue
		}
		labs = append(labs, lab)
	}

	return labs, nil
}

func (l *labService) AddPublicLab(lab entity.LabType) error {
	// If lab Id is not yet generated Generate
	if lab.Id == "" {
		lab.Id = helper.Generate(20)
	}

	val, err := json.Marshal(lab)
	if err != nil {
		slog.Error("not able to convert object to string", err)
		return err
	}

	if err := l.labRepository.AddLab(lab.Id, string(val), lab.Type); err != nil {
		slog.Error("not able to save lab", err)
		return err
	}

	return nil
}

func (l *labService) DeletePublicLab(lab entity.LabType) error {
	if err := l.labRepository.DeleteLab(lab.Id, lab.Type); err != nil {
		slog.Error("not able to delete lab", err)
		return err
	}
	return nil
}
