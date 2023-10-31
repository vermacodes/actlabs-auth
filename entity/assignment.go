package entity

type Assignment struct {
	Id      string `json:"id"`
	User    string `json:"user"`
	LabId   string `json:"labId"`
	LabName string `json:"labName"`
	Status  string `json:"status"`
}

type AssignmentService interface {
	GetAssignments() ([]Assignment, error)
	GetMyAssignments(userPrincipal string) ([]LabType, error)
	CreateAssignment(Assignment) error
	// TODO: UpdateAssignment(Assignment) error
	DeleteAssignment(Assignment) error
}

type AssignmentRepository interface {
	// List of all the available assignments.
	GetEnumerationResults() (EnumerationResults, error)
	GetAssignment(name string) (Assignment, error)
	DeleteAssignment(assignmentId string) error
	CreateAssignment(assignmentId string, assignment string) error

	ValidateUser(userId string) (bool, error)
}
