package handlers

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/iki-rumondor/go-speech/internal/domain/layers/services"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/request"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/response"
	"github.com/iki-rumondor/go-speech/internal/utils"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var body request.User
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if err := h.Service.CreateUser(&body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SUCCESS_RES("Pengguna Berhasil Ditambahkan"))
}

func (h *UserHandler) GetUser(c *gin.Context) {
	uuid := c.Param("uuid")
	resp, err := h.Service.GetUser(uuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *UserHandler) VerifyUser(c *gin.Context) {
	var body request.SignIn
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	resp, err := h.Service.VerifyUser(&body)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *UserHandler) GetRoles(c *gin.Context) {

	resp, err := h.Service.GetRoles()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *UserHandler) GetTeachers(c *gin.Context) {

	resp, err := h.Service.GetTeachers()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *UserHandler) ActivateUser(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.Service.ActivateUser(uuid); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Aktivasi User Berhasil"))
}

func (h *UserHandler) GetAllClasses(c *gin.Context) {

	resp, err := h.Service.GetAllClasses()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *UserHandler) CreateClassRequest(c *gin.Context) {
	var body request.ClassRequest
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	userUuid := c.GetString("uuid")
	if userUuid == "" {
		utils.HandleError(c, response.HANDLER_INTERR)
		return
	}

	if err := h.Service.CreateClassRequest(userUuid, &body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Berhasil Mendaftar Di Kelas"))
}

func (h *UserHandler) GetStudentRequestClasses(c *gin.Context) {
	userUuid := c.GetString("uuid")
	if userUuid == "" {
		utils.HandleError(c, response.HANDLER_INTERR)
		return
	}
	resp, err := h.Service.GetStudentRequestClasses(userUuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *UserHandler) DashboardAdmin(c *gin.Context) {
	resp, err := h.Service.DashboardAdmin()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *UserHandler) DashboardTeacher(c *gin.Context) {
	userUuid := c.GetString("uuid")
	resp, err := h.Service.DashboardTeacher(userUuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *UserHandler) DashboardStudent(c *gin.Context) {
	userUuid := c.GetString("uuid")
	resp, err := h.Service.DashboardStudent(userUuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *UserHandler) GetRequestClasses(c *gin.Context) {
	userUuid := c.GetString("uuid")
	if userUuid == "" {
		utils.HandleError(c, response.HANDLER_INTERR)
		return
	}
	resp, err := h.Service.GetRequestClasses(userUuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *UserHandler) UpdateStatusClassReq(c *gin.Context) {
	var body request.StatusClassReq
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	uuid := c.Param("uuid")

	if err := h.Service.UpdateStatusClassReq(uuid, &body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Berhasil Mendaftar Di Kelas"))
}
