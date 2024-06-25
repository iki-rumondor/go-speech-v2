package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	aai "github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/iki-rumondor/go-speech/internal/domain/layers/interfaces"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MasterRepo struct {
	db *gorm.DB
}

const (
	DropboxToken = ""
)

func NewMasterInterface(db *gorm.DB) interfaces.MasterInterface {
	return &MasterRepo{
		db: db,
	}
}

func (r *MasterRepo) Create(model interface{}) error {
	return r.db.Create(model).Error
}

func (r *MasterRepo) Find(dest interface{}, condition, order string) error {
	return r.db.Preload(clause.Associations).Order(order).Find(dest, condition).Error
}

func (r *MasterRepo) First(dest interface{}, condition string) error {
	return r.db.Preload(clause.Associations).First(dest, condition).Error
}

func (r *MasterRepo) Update(model interface{}, condition string) error {
	return r.db.Where(condition).Updates(model).Error
}

func (r *MasterRepo) Delete(data interface{}, withAssociation []string) error {
	return r.db.Select(withAssociation).Delete(data).Error
}

func (r *MasterRepo) Distinct(model interface{}, column, condition string, dest *[]string) error {
	return r.db.Model(model).Distinct().Where(condition).Pluck(column, dest).Error
}

func (r *MasterRepo) Truncate(tableName string) error {
	return r.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error
}

func (r *MasterRepo) FindStudentClasses(studentID uint, dest *[]models.Class) error {
	subQuery := r.db.Model(&models.ClassRequest{}).Where("student_id = ? AND status = ?", studentID, 2).Select("class_id")
	return r.db.Find(dest, "id IN (?)", subQuery).Error
}

func (r *MasterRepo) AudioToTextAPI(audioUrl string) error {
	assemblyKey := os.Getenv("ASSEMBLY_KEY")
	if assemblyKey == "" {
		return errors.New("assembly env not found")
	}

	apiUrl := "https://api.assemblyai.com/v2/transcript"

	values := map[string]string{"audio_url": audioUrl}
	jsonData, err := json.Marshal(values)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", assemblyKey)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	transcript, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Transcript:", string(transcript))

	return nil
}

func (r *MasterRepo) AudioToSubtitleTranscript(audioUrl string) ([]byte, error) {
	assemblyKey := os.Getenv("ASSEMBLY_KEY")
	if assemblyKey == "" {
		return nil, errors.New("assembly env not found")
	}

	ctx := context.Background()

	client := aai.NewClient(assemblyKey)

	transcript, err := client.Transcripts.TranscribeFromURL(ctx, audioUrl, &aai.TranscriptOptionalParams{
		// LanguageDetection: aai.Bool(true),
	})

	// log.Println(aai.ToString(transcript.Text))

	if err != nil {
		return nil, err
	}

	params := &aai.TranscriptGetSubtitlesOptions{
		CharsPerCaption: 32,
	}

	vtt, err := client.Transcripts.GetSubtitles(ctx, aai.ToString(transcript.ID), "vtt", params)
	if err != nil {
		return nil, err
	}

	return vtt, nil
}

func (r *MasterRepo) UploadAudio(audioPath string) (map[string]interface{}, error) {
	assemblyKey := os.Getenv("ASSEMBLY_KEY")
	if assemblyKey == "" {
		return nil, errors.New("assembly env not found")
	}

	apiUrl := "https://api.assemblyai.com/v2/upload"

	audioFile, err := os.ReadFile(audioPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(audioFile))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", assemblyKey)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result, nil
}

func (r *MasterRepo) UploadToDropbox(bookName string) error {
	booksFolder := "internal/files/books"
	filePath := filepath.Join(booksFolder, bookName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	dbxArg := map[string]interface{}{
		"path":            fmt.Sprintf("/%s", bookName),
		"mode":            "overwrite",
		"autorename":      true,
		"mute":            false,
		"strict_conflict": false,
	}

	dbxArgJSON, err := json.Marshal(dbxArg)
	if err != nil {
		return err
	}

	url := "https://content.dropboxapi.com/2/files/upload"
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+DropboxToken)
	req.Header.Set("Dropbox-API-Arg", string(dbxArgJSON))
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	if resp.StatusCode != 200 {
		fmt.Println("headers:", resp.Header)
		return errors.New("something went error")
	}

	return nil
}

func (r *MasterRepo) GetDropboxURL(bookName string) (string, error) {

	data := map[string]interface{}{
		"path": fmt.Sprintf("/%s", bookName),
		"settings": map[string]interface{}{
			"access":               "viewer",
			"allow_download":       true,
			"audience":             "public",
			"requested_visibility": "public",
		},
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	url := "https://api.dropboxapi.com/2/sharing/create_shared_link_with_settings"
	req, err := http.NewRequest("POST", url, bytes.NewReader(dataJSON))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+DropboxToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(body))

	if resp.StatusCode != 200 {
		fmt.Println("headers:", resp.Header)
		return "", errors.New("something went error")
	}

	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return "", err
	}

	url, ok := respData["url"].(string)
	if !ok {
		return "", errors.New("failed to get shared link URL")
	}

	return url, nil
}

func (r *MasterRepo) UploadFlipbookHeyzine(pdfURL string) (map[string]interface{}, error) {
	heyzineAPIKey := os.Getenv("HEYZINE_KEY")
	if heyzineAPIKey == "" {
		return nil, errors.New("heyzine key not found")
	}

	url := fmt.Sprintf("https://heyzine.com/api1/rest?pdf=%s&k=%s", pdfURL, heyzineAPIKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respData map[string]interface{}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return nil, err
	}

	if !respData["success"].(bool) {
		log.Println(respData)
		return nil, errors.New("something went error")
	}

	return respData, nil
}
