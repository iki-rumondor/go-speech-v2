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
	}
}
