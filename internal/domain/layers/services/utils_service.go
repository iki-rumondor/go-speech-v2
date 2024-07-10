package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/iki-rumondor/go-speech/internal/consts"
	"github.com/iki-rumondor/go-speech/internal/utils"
)

func (s *MasterService) Assembly_VideoToSubtitle(videoName string) (*string, error) {

	audioName := utils.GenerateRandomString(12)
	audioBasePath := "internal/files/audio"
	audioPath := filepath.Join(audioBasePath, fmt.Sprintf("%s.mp3", audioName))

	if err := os.MkdirAll(filepath.Dir(audioPath), 0750); err != nil {
		log.Println(err)
		return nil, err
	}

	videoPath := filepath.Join(consts.VIDEOS_FOLDER, videoName)

	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vn", "-acodec", "libmp3lame", "-q:a", "4", audioPath)

	if err := cmd.Run(); err != nil {
		log.Println(err)
		return nil, err
	}

	defer func() {
		if err := os.Remove(audioPath); err != nil {
			log.Println(err.Error())
		}
	}()

	result, err := s.Repo.UploadAudioToAssembly(audioPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	subtitle, err := s.Repo.AudioToSubtitleTranscript(result["upload_url"].(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	subtitleName := fmt.Sprintf("%s.vtt", utils.GenerateRandomString(12))
	subtitlePath := filepath.Join(consts.GetSubtitleFolder(), subtitleName)

	if err := os.MkdirAll(filepath.Dir(subtitlePath), 0750); err != nil {
		log.Println(err)
		return nil, err
	}

	if err := os.WriteFile(subtitlePath, subtitle, 0644); err != nil {
		log.Println(err)
		return nil, err
	}

	return &subtitleName, nil
}
