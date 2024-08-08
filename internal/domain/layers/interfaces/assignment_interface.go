package interfaces

import "github.com/iki-rumondor/go-speech/internal/domain/structs/models"

type AssignmentInterface interface {
	CreateAssignment(classUuid string, model *models.Assignment) error
	FindAssignmentsByClass(classUuid string) (*[]models.Assignment, error)
	FindAssignmentsByUser(userUuid string) (*[]models.Assignment, error)
	FirstAssignmentByUuid(uuid string) (*models.Assignment, error)
	UpdateAssignment(model *models.Assignment) error
	DeleteAssignment(uuid string) error

	// FindAnswers(assignmentUuid string) (*[]models.Answer, error)
	FirstAnswerByUser(userUuid, assignmentUuid string) (*models.Answer, error)
	CreateAnswer(assignmentUuid, userUuid string, model *models.Answer) error
	UpdateAnswer(answerUuid string, model *models.Answer) error
}
