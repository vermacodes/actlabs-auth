package entity

type Assigment struct {
	Id      string `json:"id"`
	User    string `json:"user"`
	LabId   string `json:"labId"`
	LabName string `json:"labName"`
	Status  string `json:"status"`
}

type AssignmentService interface {
	GetAssignments() ([]Assigment, error)
	GetMyAssignments(userPrincipal string) ([]LabType, error)
	CreateAssignment(Assigment) error
	// TODO: UpdateAssignment(Assigment) error
	DeleteAssignment(Assigment) error
}

type AssignmentRepository interface {
	// List of all the available assignments.
	GetEnumerationResults() (EnumerationResults, error)
	GetAssignment(name string) (Assigment, error)
	DeleteAssignment(assignmentId string) error
	CreateAssignment(assignmentId string, assignment string) error

	ValidateUser(userId string) (bool, error)
}
