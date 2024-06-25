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

type FileHandler struct {
	Service *services.FileService
}

func NewFileHandler(service *services.FileService) *FileHandler {
	return &FileHandler{
		Service: service,
	}
}

func (h *FileHandler) CreateVideo(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if status := utils.CheckTypeFile(file, []string{"mp4", "webm"}); !status {
		utils.HandleError(c, response.BADREQ_ERR("Tipe File Tidak Valid, Gunakan mp4 atau webm"))
		return
	}

	title := c.PostForm("title")
	if title == "" {
		utils.HandleError(c, response.BADREQ_ERR("Judul Tidak Ditemukan"))
		return
	}

	classUuid := c.PostForm("class_uuid")
	if classUuid == "" {
		utils.HandleError(c, response.BADREQ_ERR("Uuid Kelas Tidak Ditemukan"))
		return
	}

	description := c.PostForm("description")
	if description == "" {
		utils.HandleError(c, response.BADREQ_ERR("Deskripsi Tidak Ditemukan"))
		return
	}

	tempFolder := "internal/files/videos"
	videoName := utils.RandomFileName(file)
	pathFile := filepath.Join(tempFolder, videoName)

	if err := utils.SaveUploadedFile(file, pathFile); err != nil {
		utils.HandleError(c, response.HANDLER_INTERR)
		return
	}

	if err := h.Service.CreateVideo(pathFile, videoName, title, description, classUuid); err != nil {
		if err := os.Remove(pathFile); err != nil {
			log.Println(err.Error())
		}
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SUCCESS_RES("Berhasil Menambahkan Video Pembelajaran"))
}

func (h *FileHandler) UpdateVideo(c *gin.Context) {
	var body request.UpdateVideo
	if err := c.BindJSON(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if _, err := govalidator.ValidateStruct(&body); err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	uuid := c.Param("uuid")
	if err := h.Service.UpdateVideo(uuid, &body); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Video Berhasil Diperbarui"))
}

func (h *FileHandler) DeleteVideo(c *gin.Context) {

	uuid := c.Param("uuid")
	if err := h.Service.DeleteVideo(uuid); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.SUCCESS_RES("Video Berhasil Diperbarui"))
}

func (h *FileHandler) CreateBook(c *gin.Context) {
	file, err := c.FormFile("book")
	if err != nil {
		utils.HandleError(c, response.BADREQ_ERR(err.Error()))
		return
	}

	if status := utils.CheckTypeFile(file, []string{"pdf"}); !status {
		utils.HandleError(c, response.BADREQ_ERR("Tipe File Tidak Valid, Gunakan pdf"))
		return
	}

	title := c.PostForm("title")
	if title == "" {
		utils.HandleError(c, response.BADREQ_ERR("Judul Tidak Ditemukan"))
		return
	}

	classUuid := c.PostForm("class_uuid")
	if classUuid == "" {
		utils.HandleError(c, response.BADREQ_ERR("Uuid Kelas Tidak Ditemukan"))
		return
	}

	description := c.PostForm("description")
	if description == "" {
		utils.HandleError(c, response.BADREQ_ERR("Deskripsi Tidak Ditemukan"))
		return
	}

	booksFolder := "internal/files/books"
	bookName := utils.RandomFileName(file)
	pathFile := filepath.Join(booksFolder, bookName)

	if err := utils.SaveUploadedFile(file, pathFile); err != nil {
		utils.HandleError(c, response.HANDLER_INTERR)
		return
	}

	if err := h.Service.CreateBook(bookName, title, description, classUuid); err != nil {
		if err := os.Remove(pathFile); err != nil {
			log.Println(err.Error())
		}
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SUCCESS_RES("Berhasil Menambahkan Buku Pembelajaran"))
}

func (h *FileHandler) DeleteBook(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.Service.DeleteBook(uuid); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.SUCCESS_RES("Berhasil Menambahkan Buku Pembelajaran"))
}

func (h *FileHandler) GetClassVideos(c *gin.Context) {
	classUuid := c.Param("uuid")
	resp, err := h.Service.GetClassVideos(classUuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *FileHandler) GetClassBooks(c *gin.Context) {
	classUuid := c.Param("uuid")
	resp, err := h.Service.GetClassBooks(classUuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *FileHandler) GetVideo(c *gin.Context) {
	uuid := c.Param("uuid")
	resp, err := h.Service.GetVideo(uuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *FileHandler) GetBook(c *gin.Context) {
	uuid := c.Param("uuid")
	resp, err := h.Service.GetBook(uuid)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DATA_RES(resp))
}

func (h *FileHandler) GetFileVideo(c *gin.Context) {
	videoName := c.Param("videoName")
	tempFolder := "internal/files/videos"
	pathFile := filepath.Join(tempFolder, videoName)
	c.File(pathFile)
}

func (h *FileHandler) GetFileSubtitle(c *gin.Context) {
	subtitle := c.Param("subtitle")
	tempFolder := "internal/files/subtitle"
	pathFile := filepath.Join(tempFolder, subtitle)
	c.File(pathFile)
}

func (h *FileHandler) GetFileBook(c *gin.Context) {
	bookName := c.Param("bookName")
	tempFolder := "internal/files/books"
	pathFile := filepath.Join(tempFolder, bookName)
	c.File(pathFile)
}
