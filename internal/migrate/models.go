package migrate

import "github.com/iki-rumondor/go-speech/internal/domain/structs/models"

type Model struct {
	Model interface{}
}

func GetAllModels() []Model {
	return []Model{
		{Model: models.Role{}},
		{Model: models.User{}},
		{Model: models.Department{}},
		{Model: models.Teacher{}},
		{Model: models.Student{}},
		{Model: models.Class{}},
		{Model: models.ClassRequest{}},
		{Model: models.Video{}},
		{Model: models.Book{}},
		{Model: models.Note{}},
		{Model: models.FlipBook{}},
		{Model: models.VideoPart{}},
		{Model: models.BookPart{}},
		{Model: models.Material{}},
		{Model: models.StudentClasses{}},
		{Model: models.Assignment{}},
		{Model: models.Answer{}},
		{Model: models.Notification{}},
		{Model: models.ClassNotification{}},
		{Model: models.ReadNotification{}},
	}
}
