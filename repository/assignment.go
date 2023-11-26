package repository

import (
	"context"
	"encoding/json"
	"strings"

	"actlabs-auth/entity"

	"golang.org/x/exp/slog"
)

type assignmentRepository struct{}

func NewAssignmentRepository() entity.AssignmentRepository {
	return &assignmentRepository{}
}

// List of all the available assignments.
// func (a *assignmentRepository) GetEnumerationResults() (entity.EnumerationResults, error) {
// 	er := entity.EnumerationResults{}

// 	// URL of the container to list the blobs
// 	url := "https://" + entity.StorageAccountName + ".blob.core.windows.net/repro-project-assignments" + entity.SasToken + "&restype=container&comp=list"

// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return er, err
// 	}

// 	req.Header.Add("Accept", "application/json")
// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return er, err
// 	}

// 	defer res.Body.Close()
// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return er, err
// 	}

// 	if err := xml.Unmarshal(body, &er); err != nil {
// 		return er, err
// 	}

// 	return er, nil
// }

func (a *assignmentRepository) GetAllAssignments() ([]entity.Assignment, error) {
	assignment := entity.Assignment{}
	assignments := []entity.Assignment{}

	// URL of the blob with SAS token.
	serviceClient := getServiceClient().NewClient("ReadinessAssignments")
	pager := serviceClient.NewListEntitiesPager(nil)
	for pager.More() {
		response, err := pager.NextPage(context.Background())
		if err != nil {
			slog.Error("Error getting entities: ", err)
			return assignments, err
		}

		for _, element := range response.Entities {
			//var myEntity aztables.EDMEntity
			if err := json.Unmarshal(element, &assignment); err != nil {
				slog.Error("Error unmarshal entity: ", err)
				return assignments, err
			}
			assignments = append(assignments, assignment)
		}
	}

	return assignments, nil
}

func (a *assignmentRepository) GetAssignmentsByLabId(labId string) ([]entity.Assignment, error) {
	assignment := entity.Assignment{}
	assignments := []entity.Assignment{}

	// URL of the blob with SAS token.
	serviceClient := getServiceClient().NewClient("ReadinessAssignments")
	pager := serviceClient.NewListEntitiesPager(nil)
	for pager.More() {
		response, err := pager.NextPage(context.Background())
		if err != nil {
			slog.Error("Error getting entities: ", err)
			return assignments, err
		}

		for _, element := range response.Entities {
			//var myEntity aztables.EDMEntity
			if err := json.Unmarshal(element, &assignment); err != nil {
				slog.Error("Error unmarshal entity: ", err)
				return assignments, err
			}

			if assignment.LabId == labId {
				assignments = append(assignments, assignment)
			}
		}
	}

	return assignments, nil
}

func (a *assignmentRepository) GetAssignmentsByUserId(userId string) ([]entity.Assignment, error) {
	assignment := entity.Assignment{}
	assignments := []entity.Assignment{}

	// URL of the blob with SAS token.
	serviceClient := getServiceClient().NewClient("ReadinessAssignments")
	pager := serviceClient.NewListEntitiesPager(nil)
	for pager.More() {
		response, err := pager.NextPage(context.Background())
		if err != nil {
			slog.Error("Error getting entities: ", err)
			return assignments, err
		}

		for _, element := range response.Entities {
			//var myEntity aztables.EDMEntity
			if err := json.Unmarshal(element, &assignment); err != nil {
				slog.Error("Error unmarshal entity: ", err)
				return assignments, err
			}

			if assignment.UserId == userId {
				assignments = append(assignments, assignment)
			}
		}
	}

	return assignments, nil
}

func (a *assignmentRepository) DeleteAssignment(assignmentId string) error {

	slog.Debug("Deleting assignment: ", assignmentId)

	userId := assignmentId[:strings.Index(assignmentId, "-")]

	getServiceClient := getServiceClient().NewClient("ReadinessAssignments")
	_, err := getServiceClient.DeleteEntity(context.Background(), userId, assignmentId, nil)
	if err != nil {
		slog.Error("Error deleting assignment record: ", err)
		return err
	}
	slog.Debug("Assignment record deleted successfully")
	return nil
}

func (a *assignmentRepository) UpsertAssignment(assignment entity.Assignment) error {
	serviceClient := getServiceClient().NewClient("ReadinessAssignments")
	assignmentRecord := entity.Assignment{
		PartitionKey: assignment.UserId,
		RowKey:       assignment.AssignmentId,
		AssignmentId: assignment.AssignmentId,
		UserId:       assignment.UserId,
		LabId:        assignment.LabId,
		CreatedBy:    assignment.CreatedBy,
		CreatedOn:    assignment.CreatedOn,
		Status:       assignment.Status,
	}

	val, err := json.Marshal(assignmentRecord)
	if err != nil {
		slog.Error("Error marshalling assignment record: ", err)
		return err
	}

	slog.Debug("Assignment record: ", string(val))

	_, err = serviceClient.UpsertEntity(context.TODO(), val, nil)

	if err != nil {
		slog.Error("Error creating assignment record: ", err)
		return err
	}

	slog.Debug("Assignment record created successfully")

	return nil
}

func (a *assignmentRepository) ValidateUser(userId string) (bool, error) {
	// out, err := exec.Command("bash", "-c", "az ad user show --id '"+userId+"' --query 'accountEnabled' --output json 2>/dev/null").Output()
	// if err != nil {
	// 	return false, err
	// }

	// outString := string(out)
	// outStringTrimmed := strings.TrimRight(outString, "\n")
	// //Convert to boolean.
	// boolVal, err := strconv.ParseBool(outStringTrimmed)
	// if err != nil {
	// 	return false, err
	// }

	// return boolVal, nil
	return true, nil
}
