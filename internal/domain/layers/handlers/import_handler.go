package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/iki-rumondor/go-speech/internal/domain/layers/services"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/response"
	"github.com/iki-rumondor/go-speech/internal/utils"
)

type ImportHandler struct {
	Service *services.ImportService
}

func NewImportHandler(service *services.ImportService) *ImportHandler {
	return &ImportHandler{
		Service: service,
	}
}

func (h *ImportHandler) ImportTeachers(c *gin.Context) {

	teacherFile, err := c.FormFile("teachers")
	if err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if status := utils.CheckTypeFile(teacherFile, []string{"xlsx"}); !status {
		utils.HandleError(c, response.BADREQ_ERR("Tipe File Tidak Valid, Gunakan file xlsx"))
		return
	}

	tempFolder := "internal/files/temp"
	tempPath := filepath.Join(tempFolder, teacherFile.Filename)

	if err := utils.SaveUploadedFile(teacherFile, tempPath); err != nil {
		utils.HandleError(c, response.HANDLER_INTERR)
		return
	}

	defer func() {
		if err := os.Remove(tempPath); err != nil {
			log.Println(err.Error())
		}
	}()

	teachers, err := h.Service.SaveTeachers(tempPath)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.DATA_RES(teachers))
}

func (h *ImportHandler) ImportStudents(c *gin.Context) {

	studentsFile, err := c.FormFile("students")
	if err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if status := utils.CheckTypeFile(studentsFile, []string{"xlsx"}); !status {
		utils.HandleError(c, response.BADREQ_ERR("Tipe File Tidak Valid, Gunakan file xlsx"))
		return
	}

	tempFolder := "internal/files/temp"
	tempPath := filepath.Join(tempFolder, studentsFile.Filename)

	if err := utils.SaveUploadedFile(studentsFile, tempPath); err != nil {
		utils.HandleError(c, response.HANDLER_INTERR)
		return
	}

	defer func() {
		if err := os.Remove(tempPath); err != nil {
			log.Println(err.Error())
		}
	}()

	teachers, err := h.Service.SaveStudents(tempPath)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.DATA_RES(teachers))
}
