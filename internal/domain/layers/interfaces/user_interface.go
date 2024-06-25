package interfaces

import (
	"github.com/iki-rumondor/go-speech/internal/domain/structs/models"
	"gorm.io/gorm"
)

type UserInterface interface {
	Find(dest interface{}, condition, order string) error
	First(dest interface{}, condition string) error
	Truncate(tableName string) error
	Distinct(model interface{}, column, condition string, dest *[]string) error
	Create(data interface{}) error
	Update(data interface{}, condition string) error
	Delete(data interface{}, withAssociation []string) error

	FindClasses(model *[]models.Class, condition string) error
	FindTeacherStudents(model *[]models.ClassRequest, teacherID uint) error
	SelectColumn(model interface{}, condition, columnName string) *gorm.DB
	Include(model interface{}, condition, colName string, selectCols *gorm.DB) error
}
