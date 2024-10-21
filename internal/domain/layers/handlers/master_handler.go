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

type MasterHandler struct {
	Service *services.MasterService
}

func NewMasterHandler(service *services.MasterService) *MasterHandler {
	return &MasterHandler{
		Service: service,
	}
}

func (h *MasterHandler) CreateClass(c *gin.Context) {
	var body request.Class
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if err := h.Service.CreateClass(&body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SUCCESS_RES("Kelas Berhasil Ditambahkan"))
}

func (h *MasterHandler) CreateSubject(c *gin.Context) {
	var body request.Class
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if err := h.Service.CreateSubject(&body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SUCCESS_RES("Mata Kuliah Berhasil Ditambahkan"))
}

func (h *MasterHandler) DeleteClass(c *gin.Context) {

	uuid := c.Param("uuid")

	if err := h.Service.DeleteClass(uuid); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Kelas Berhasil Dihapus"))
}

func (h *MasterHandler) CreateDepartment(c *gin.Context) {
	var body request.Department
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if err := h.Service.CreateDepartment(&body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SUCCESS_RES("Program Studi Berhasil Ditambahkan"))
}

func (h *MasterHandler) UpdateDepartment(c *gin.Context) {
	var body request.Department
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	uuid := c.Param("uuid")
	if err := h.Service.UpdateDepartment(uuid, &body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Program Studi Berhasil Diperbarui"))
}

func (h *MasterHandler) DeleteDepartment(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.Service.DeleteDepartment(uuid); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Program Studi Berhasil Dihapus"))
}

func (h *MasterHandler) GetAllDepartment(c *gin.Context) {

	resp, err := h.Service.GetAllDepartment()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *MasterHandler) GetTeacherClasses(c *gin.Context) {
	userUuid := c.GetString("uuid")
	if userUuid == "" {
		utils.HandleError(c, response.HANDLER_INTERR)
		return
	}

	resp, err := h.Service.GetTeacherClasses(userUuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *MasterHandler) GetAllClasses(c *gin.Context) {
	resp, err := h.Service.GetAllClasses()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *MasterHandler) GetAllSubjects(c *gin.Context) {
	resp, err := h.Service.GetAllSubjects()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *MasterHandler) GetStudentClasses(c *gin.Context) {
	userUuid := c.Param("userUuid")
	resp, err := h.Service.GetClass(userUuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *MasterHandler) GetDepartment(c *gin.Context) {

	uuid := c.Param("uuid")
	resp, err := h.Service.GetDepartment(uuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *MasterHandler) GetClass(c *gin.Context) {
	uuid := c.Param("uuid")
	resp, err := h.Service.GetClass(uuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *MasterHandler) UpdateClass(c *gin.Context) {
	var body request.Class
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	uuid := c.Param("uuid")
	if err := h.Service.UpdateClass(uuid, &body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Program Studi Berhasil Diperbarui"))
}

func (h *MasterHandler) CreateNote(c *gin.Context) {
	var body request.Note
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if err := h.Service.CreateNote(&body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SUCCESS_RES("Catatan Berhasil Ditambahkan"))
}

func (h *MasterHandler) GetNotes(c *gin.Context) {
	classUuid := c.Param("uuid")
	resp, err := h.Service.GetNotes(classUuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *MasterHandler) GetNote(c *gin.Context) {
	uuid := c.Param("uuid")
	resp, err := h.Service.GetNote(uuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *MasterHandler) UpdateNote(c *gin.Context) {
	var body request.Note
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	uuid := c.Param("uuid")
	if err := h.Service.UpdateNote(uuid, &body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SUCCESS_RES("Catatan Berhasil Diperbarui"))
}

func (h *MasterHandler) DeleteNote(c *gin.Context) {

	uuid := c.Param("uuid")
	if err := h.Service.DeleteNote(uuid); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SUCCESS_RES("Catatan Berhasil Dihapus"))
}

func (h *MasterHandler) GetClassesReport(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.Service.GetClassesReport(uuid); err != nil {
		utils.HandleError(c, err)
		return
	}

	reportFolder := "internal/files/reports"
	pathFile := filepath.Join(reportFolder, "class_students.pdf")
	c.File(pathFile)
}

func (h *MasterHandler) GetStudentAssignmentsReport(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.Service.GetClassesReport(uuid); err != nil {
		utils.HandleError(c, err)
		return
	}

	reportFolder := "internal/files/reports"
	pathFile := filepath.Join(reportFolder, "class_students.pdf")
	c.File(pathFile)
}

func (h *MasterHandler) CreateMaterial(c *gin.Context) {
	var body request.Material
	if err := c.Bind(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	video, err := c.FormFile("video")
	if err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if status := utils.CheckTypeFile(video, []string{"mp4", "webm"}); !status {
		utils.HandleError(c, response.BADREQ_ERR("Tipe File Tidak Valid, Gunakan mp4 atau webm"))
		return
	}

	book, err := c.FormFile("book")
	if err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if status := utils.CheckTypeFile(book, []string{"pdf"}); !status {
		utils.HandleError(c, response.BADREQ_ERR("Tipe File Tidak Valid, Gunakan pdf"))
		return
	}

	videosFolder := "internal/files/videos"
	videoName := utils.RandomFileName(video)
	pathVideo := filepath.Join(videosFolder, videoName)

	if err := utils.SaveUploadedFile(video, pathVideo); err != nil {
		utils.HandleError(c, response.HANDLER_INTERR)
		return
	}

	booksFolder := "internal/files/books"
	bookName := utils.RandomFileName(book)
	pathBook := filepath.Join(booksFolder, bookName)

	if err := utils.SaveUploadedFile(book, pathBook); err != nil {
		utils.HandleError(c, response.HANDLER_INTERR)
		return
	}

	if err := h.Service.CreateMaterial(&body, videoName, bookName); err != nil {
		if err := os.Remove(pathVideo); err != nil {
			log.Println(err.Error())
		}

		if err := os.Remove(pathBook); err != nil {
			log.Println(err.Error())
		}

		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SUCCESS_RES("Materi Berhasil Ditambahkan"))
}

func (h *MasterHandler) GetAllMaterials(c *gin.Context) {
	classUuid := c.Param("class_uuid")
	resp, err := h.Service.GetAllMaterials(classUuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *MasterHandler) GetMaterial(c *gin.Context) {
	uuid := c.Param("uuid")
	resp, err := h.Service.GetMaterial(uuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *MasterHandler) UpdateMaterial(c *gin.Context) {
	var body request.UpdateMaterial
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	uuid := c.Param("uuid")
	if err := h.Service.UpdateMaterial(uuid, &body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SUCCESS_RES("Materi Berhasil Diperbarui"))
}

func (h *MasterHandler) DeleteMaterial(c *gin.Context) {

	uuid := c.Param("uuid")
	if err := h.Service.DeleteMaterial(uuid); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SUCCESS_RES("Materi Berhasil Dihapus"))
}

func (h *MasterHandler) GetTeacher(c *gin.Context) {
	uuid := c.Param("uuid")
	resp, err := h.Service.GetTeacher(uuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *MasterHandler) CreateTeacher(c *gin.Context) {
	var body request.CreateTeacher
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if err := h.Service.CreateTeacher(&body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SUCCESS_RES("Dosen Berhasil Ditambahkan"))
}

func (h *MasterHandler) UpdateTeacher(c *gin.Context) {
	var body request.UpdateTeacher
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	uuid := c.Param("uuid")
	if err := h.Service.UpdateTeacher(uuid, &body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Dosen Berhasil Diperbarui"))
}

func (h *MasterHandler) DeleteTeacher(c *gin.Context) {

	uuid := c.Param("uuid")
	if err := h.Service.DeleteTeacher(uuid); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Dosen Berhasil Dihapus"))
}

func (h *MasterHandler) GetStudents(c *gin.Context) {

	resp, err := h.Service.GetStudents()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *MasterHandler) GetStudent(c *gin.Context) {
	uuid := c.Param("uuid")
	resp, err := h.Service.GetStudent(uuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *MasterHandler) CreateStudent(c *gin.Context) {
	var body request.CreateStudent
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if err := h.Service.CreateStudent(&body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Data Mahasiswa Berhasil Ditambahkan"))
}

func (h *MasterHandler) UpdateStudent(c *gin.Context) {
	var body request.UpdateStudent
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	uuid := c.Param("uuid")
	if err := h.Service.UpdateStudent(uuid, &body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Data Mahasiswa Berhasil Diperbarui"))
}

func (h *MasterHandler) DeleteStudent(c *gin.Context) {

	uuid := c.Param("uuid")
	if err := h.Service.DeleteStudent(uuid); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Data Mahasiswa Berhasil Dihapus"))
}
