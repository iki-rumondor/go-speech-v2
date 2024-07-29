package interfaces

import (
	"github.com/iki-rumondor/go-speech/internal/domain/structs/models"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/response"
)

type MasterInterface interface {
	Find(dest interface{}, condition, order string) error
	First(dest interface{}, condition string) error
	Truncate(tableName string) error
	Distinct(model interface{}, column, condition string, dest *[]string) error
	Create(data interface{}) error
	Update(data interface{}, condition string) error
	Delete(data interface{}, withAssociation []string) error
	FirstOrCreate(model interface{}, condition interface{}) error

	FindStudentClasses(studentID uint, dest *[]models.Class) error
	UpdateTeacher(teacher *models.Teacher, user *models.User) error
	UpdateStudent(student *models.Student, user *models.User) error

	UploadAudioToAssembly(audioPath string) (map[string]interface{}, error)
	AudioToTextAPI(audioPath string) error
	AudioToSubtitleTranscript(audioUrl string) ([]byte, error)
	UploadToDropbox(bookName string) error
	GetDropboxURL(bookName string) (string, error)
	UploadFlipbookHeyzine(pdfURL string) (map[string]interface{}, error)

	LaravelClassReport(data *response.Class) error
}
