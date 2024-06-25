package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/iki-rumondor/go-speech/internal/domain/layers/interfaces"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/models"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/request"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/response"
	"github.com/iki-rumondor/go-speech/internal/utils"
)

type FileService struct {
	Repo interfaces.MasterInterface
}

func NewFileService(repo interfaces.MasterInterface) *FileService {
	return &FileService{
		Repo: repo,
	}
}

func (s *FileService) CreateVideo(pathFile, videoName, title, description, classUuid string) error {
	var class models.Class
	condition := fmt.Sprintf("uuid = '%s'", classUuid)
	if err := s.Repo.First(&class, condition); err != nil {
		return response.NOTFOUND_ERR("Kelas Tidak Ditemukan")
	}

	audioName := utils.GenerateRandomString(12)
	audioBasePath := "internal/files/audio"
	audioPath := filepath.Join(audioBasePath, fmt.Sprintf("%s.mp3", audioName))

	if err := os.MkdirAll(filepath.Dir(audioPath), 0750); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	cmd := exec.Command("ffmpeg", "-i", pathFile, "-vn", "-acodec", "libmp3lame", "-q:a", "4", audioPath)

	if err := cmd.Run(); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	defer func() {
		if err := os.Remove(audioPath); err != nil {
			log.Println(err.Error())
		}
	}()

	result, err := s.Repo.UploadAudio(audioPath)
	if err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	subtitle, err := s.Repo.AudioToSubtitleTranscript(result["upload_url"].(string))
	if err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	subtitleName := fmt.Sprintf("%s.vtt", utils.GenerateRandomString(12))
	subtitleBasePath := "internal/files/subtitle"
	subtitlePath := filepath.Join(subtitleBasePath, subtitleName)

	if err := os.MkdirAll(filepath.Dir(subtitlePath), 0750); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	if err := os.WriteFile(subtitlePath, subtitle, 0644); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	model := models.Video{
		ClassID:      class.ID,
		Title:        title,
		Description:  description,
		VideoName:    videoName,
		SubtitleName: subtitleName,
	}

	if err := s.Repo.Create(&model); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	return nil
}

func (s *FileService) UpdateVideo(uuid string, req *request.UpdateVideo) error {
	var video models.Video
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&video, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	model := models.Video{
		ID:          video.ID,
		Title:       req.Title,
		Description: req.Description,
	}

	if err := s.Repo.Update(&model, ""); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}
	return nil
}

func (s *FileService) CreateBook(fileName, title, description, classUuid string) error {

	// resp, err := s.Repo.UploadFlipbookHeyzine("https://www.dropbox.com/scl/fi/c2gpatpiibt2u7iqb1s2w/4neJp3syRoWT.pdf?raw=1")
	// if err != nil {
	// 	log.Println(err)
	// 	return response.SERVICE_INTERR
	// }

	// log.Println(resp)
	// if err := s.Repo.UploadToDropbox(fileName); err != nil {
	// 	log.Println(err)
	// 	return response.SERVICE_INTERR
	// }

	// url, err := s.Repo.GetDropboxURL(fileName)
	// if err != nil {
	// 	log.Println(err)
	// 	return response.SERVICE_INTERR
	// }

	var class models.Class
	condition := fmt.Sprintf("uuid = '%s'", classUuid)
	if err := s.Repo.First(&class, condition); err != nil {
		return response.NOTFOUND_ERR("Kelas Tidak Ditemukan")
	}

	model := models.Book{
		ClassID:     class.ID,
		Title:       title,
		Description: description,
		FileName:    fileName,
	}

	if err := s.Repo.Create(&model); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	return nil
}

func (s *FileService) DeleteVideo(uuid string) error {
	var video models.Video
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&video, condition); err != nil {
		return response.SERVICE_INTERR
	}

	if err := s.Repo.Delete(&video, nil); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	subtitleFolder := "internal/files/subtitle"
	subtitlePath := filepath.Join(subtitleFolder, video.SubtitleName)

	videoFolder := "internal/files/videos"
	videoPath := filepath.Join(videoFolder, video.VideoName)

	if err := os.Remove(subtitlePath); err != nil {
		log.Println(err.Error())
	}

	if err := os.Remove(videoPath); err != nil {
		log.Println(err.Error())
	}

	return nil
}

func (s *FileService) DeleteBook(uuid string) error {
	var book models.Book
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&book, condition); err != nil {
		return response.SERVICE_INTERR
	}

	if err := s.Repo.Delete(&book, nil); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	folder := "internal/files/books"
	pathFile := filepath.Join(folder, book.FileName)

	if err := os.Remove(pathFile); err != nil {
		log.Println(err.Error())
	}

	return nil
}

func (s *FileService) GetClassVideos(classUuid string) (*[]response.Video, error) {

	var class models.Class
	condition := fmt.Sprintf("uuid = '%s'", classUuid)
	if err := s.Repo.First(&class, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var teacher models.Teacher
	condition = fmt.Sprintf("id = '%d'", class.TeacherID)
	if err := s.Repo.First(&teacher, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var model []models.Video
	condition = fmt.Sprintf("class_id = '%d'", class.ID)
	if err := s.Repo.Find(&model, condition, "id"); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp []response.Video
	for _, item := range model {
		resp = append(resp, response.Video{
			Uuid:         item.Uuid,
			Title:        item.Title,
			Description:  item.Description,
			VideoName:    item.VideoName,
			SubtitleName: item.SubtitleName,
			Teacher:      teacher.User.Name,
			CreatedAt:    item.CreatedAt,
		})
	}

	return &resp, nil
}

func (s *FileService) GetClassBooks(classUuid string) (*[]response.Book, error) {

	var class models.Class
	condition := fmt.Sprintf("uuid = '%s'", classUuid)
	if err := s.Repo.First(&class, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var teacher models.Teacher
	condition = fmt.Sprintf("id = '%d'", class.TeacherID)
	if err := s.Repo.First(&teacher, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var model []models.Book
	condition = fmt.Sprintf("class_id = '%d'", class.ID)
	if err := s.Repo.Find(&model, condition, "id"); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp []response.Book
	for _, item := range model {
		resp = append(resp, response.Book{
			Uuid:        item.Uuid,
			Title:       item.Title,
			Description: item.Description,
			FileName:    item.FileName,
			Teacher:     teacher.User.Name,
			CreatedAt:   item.CreatedAt,
		})
	}

	return &resp, nil
}

func (s *FileService) GetVideo(uuid string) (*response.Video, error) {

	var model models.Video
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&model, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp = response.Video{
		Uuid:         model.Uuid,
		Title:        model.Title,
		Description:  model.Description,
		VideoName:    model.VideoName,
		SubtitleName: model.SubtitleName,
		CreatedAt:    model.CreatedAt,
	}

	return &resp, nil
}

func (s *FileService) GetBook(uuid string) (*response.Book, error) {

	var model models.Book
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&model, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp = response.Book{
		Uuid:        model.Uuid,
		Title:       model.Title,
		Description: model.Description,
		FileName:    model.FileName,
	}

	return &resp, nil
}
