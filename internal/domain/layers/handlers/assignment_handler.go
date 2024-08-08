package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/iki-rumondor/go-speech/internal/domain/layers/services"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/request"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/response"
	"github.com/iki-rumondor/go-speech/internal/utils"
)

type AssignmentHandler struct {
	Service *services.AssignmentService
}

func NewAssignmentHandler(service *services.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{
		Service: service,
	}
}

func (h *AssignmentHandler) CreateAssignment(c *gin.Context) {

	var body request.Assignment
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if err := h.Service.CreateAssignment(&body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Tugas Berhasil Ditambahkan"))
}

func (h *AssignmentHandler) FirstAssignmentByUuid(c *gin.Context) {

	uuid := c.Param("uuid")
	resp, err := h.Service.FirstAssignmentByUuid(uuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *AssignmentHandler) FindAssignmentByClass(c *gin.Context) {

	classUuid := c.Param("classUuid")
	resp, err := h.Service.FindAssignmentByClass(classUuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *AssignmentHandler) FindAssignmentByStudent(c *gin.Context) {

	studentUuid := c.Param("studentUuid")
	resp, err := h.Service.FindAssignmentByStudent(studentUuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *AssignmentHandler) FindAssignmentByUser(c *gin.Context) {

	classUuid := c.Param("classUuid")
	userUuid := c.GetString("uuid")
	if userUuid == "" {
		utils.HandleError(c, response.BADREQ_ERR("Uuid Tidak Ditemukan"))
		return
	}
	resp, err := h.Service.FindAssignmentsByUser(userUuid, classUuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *AssignmentHandler) UpdateAssignment(c *gin.Context) {

	var body request.UpdateAssignment
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}
	uuid := c.Param("uuid")
	if err := h.Service.UpdateAssignment(uuid, &body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Tugas Berhasil Diperbarui"))
}

func (h *AssignmentHandler) DeleteAssignment(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.Service.DeleteAssignment(uuid); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Tugas Berhasil Dihapus"))
}

func (h *AssignmentHandler) UploadAnswer(c *gin.Context) {
	var assignmentUuid = c.Param("assignmentUuid")
	var userUuid = c.GetString("uuid")
	if userUuid == "" {
		utils.HandleError(c, response.BADREQ_ERR("Uuid user tidak ditemukan"))
		return
	}

	answerFile, err := c.FormFile("answer")
	if err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if status := utils.CheckTypeFile(answerFile, []string{"pdf"}); !status {
		utils.HandleError(c, response.BADREQ_ERR("Tipe File Tidak Valid, Gunakan file xlsx"))
		return
	}

	answerFolder := "internal/files/answers"
	fileName := utils.RandomFileName(answerFile)
	answerPath := filepath.Join(answerFolder, fileName)

	if err := utils.SaveUploadedFile(answerFile, answerPath); err != nil {
		utils.HandleError(c, response.HANDLER_INTERR)
		return
	}

	if err := h.Service.CreateAnswer(assignmentUuid, userUuid, fileName); err != nil {
		if err := os.Remove(answerPath); err != nil {
			log.Println(err.Error())
		}
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SUCCESS_RES("Berhasil Mengupload Tugas"))
}

func (h *AssignmentHandler) GradeAnswer(c *gin.Context) {

	var body request.Grading
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	answerUuid := c.Param("answerUuid")
	if err := h.Service.GradeAnswer(answerUuid, &body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Jawaban Berhasil Dinilai"))
}
