package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func CreateDirectory(pathName string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	savePath := filepath.Join(currentDir, pathName)

	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		if err := os.Mkdir(savePath, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func RandomFileName(file *multipart.FileHeader) string {

	fileExt := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, ".")+1:])
	randString := GenerateRandomString(12)

	return fmt.Sprintf("%s.%s", randString, fileExt)
}
