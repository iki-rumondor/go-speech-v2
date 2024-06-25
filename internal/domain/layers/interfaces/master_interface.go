package interfaces

import (
	"github.com/iki-rumondor/go-speech/internal/domain/structs/models"
)

type MasterInterface interface {
	Find(dest interface{}, condition, order string) error
	First(dest interface{}, condition string) error
	Truncate(tableName string) error
	Distinct(model interface{}, column, condition string, dest *[]string) error
	Create(data interface{}) error
	Update(data interface{}, condition string) error
	Delete(data interface{}, withAssociation []string) error

	FindStudentClasses(studentID uint, dest *[]models.Class) error

	UploadAudio(audioPath string) (map[string]interface{}, error)
	AudioToTextAPI(audioPath string) error
	AudioToSubtitleTranscript(audioUrl string) ([]byte, error)
	UploadToDropbox(bookName string) error
	GetDropboxURL(bookName string) (string, error)
	UploadFlipbookHeyzine(pdfURL string) (map[string]interface{}, error)
}
