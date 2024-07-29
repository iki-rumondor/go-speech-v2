package services

import (
	"fmt"
	"log"
	"strconv"

	"github.com/iki-rumondor/go-speech/internal/domain/layers/interfaces"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/models"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/request"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/response"
	"github.com/iki-rumondor/go-speech/internal/utils"
)

type UserService struct {
	Repo interfaces.UserInterface
}

func NewUserService(repo interfaces.UserInterface) *UserService {
	return &UserService{
		Repo: repo,
	}
}

func (s *UserService) CreateUser(req *request.User) error {

	var model models.User

	roleID, _ := strconv.Atoi(req.RoleID)

	if roleID == 2 {
		if req.DepartmentUuid == "" {
			return response.BADREQ_ERR("Field Department Tidak Ditemukan")
		}

		if req.Nip == "" {
			return response.BADREQ_ERR("Field Nip Tidak Ditemukan")
		}

		var department models.Department
		condition := fmt.Sprintf("uuid = '%s'", req.DepartmentUuid)
		if err := s.Repo.First(&department, condition); err != nil {
			log.Println(err)
			return response.SERVICE_INTERR
		}

		model = models.User{
			Name:     req.Name,
			Username: req.Email,
			Password: req.Password,
			RoleID:   uint(roleID),
			Teacher: &models.Teacher{
				Nidn:         req.Nip,
				DepartmentID: department.ID,
			},
		}
	} else {
		if req.Nim == "" {
			return response.BADREQ_ERR("Field Nim Tidak Ditemukan")
		}

		model = models.User{
			Active:   true,
			Name:     req.Name,
			Username: req.Email,
			Password: req.Password,
			RoleID:   uint(roleID),
			Student: &models.Student{
				Nim: req.Nim,
			},
		}
	}

	if err := s.Repo.Create(&model); err != nil {
		log.Println(err)
		if utils.IsErrorType(err) {
			return err
		}
		return response.SERVICE_INTERR
	}

	return nil
}

func (s *UserService) GetUser(uuid string) (*response.User, error) {

	var model models.User
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&model, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	resp := response.User{
		Uuid:  model.Uuid,
		Email: model.Username,
		Name:  model.Name,
		Role:  model.Role.Name,
	}

	return &resp, nil

}

func (s *UserService) GetRoles() (*[]response.Role, error) {

	var model []models.Role

	if err := s.Repo.Find(&model, "", "id"); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp []response.Role
	for _, item := range model {
		resp = append(resp, response.Role{
			ID:   fmt.Sprintf("%d", item.ID),
			Name: item.Name,
		})
	}

	return &resp, nil

}

func (s *UserService) GetTeachers() (*[]response.Teacher, error) {

	var model []models.Teacher

	if err := s.Repo.Find(&model, "", "id"); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp []response.Teacher
	for _, item := range model {
		resp = append(resp, response.Teacher{
			Uuid:           item.Uuid,
			Nidn:           item.Nidn,
			Department:     item.Department.Name,
			DepartmentUuid: item.Department.Uuid,
			Email:          item.User.Username,
			Name:           item.User.Name,
			Active:         item.User.Active,
		})
	}

	return &resp, nil
}

func (s *UserService) VerifyUser(req *request.SignIn) (map[string]string, error) {
	var user models.User
	condition := fmt.Sprintf("username = '%s'", req.Email)
	if err := s.Repo.First(&user, condition); err != nil {
		log.Println(err)
		return nil, response.NOTFOUND_ERR("Username atau Password Salah")
	}

	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return nil, response.NOTFOUND_ERR("Username atau Password Salah")
	}

	if !user.Active {
		return nil, response.NOTFOUND_ERR("Akun Anda Belum Diaktifkan")
	}

	jwt, err := utils.GenerateToken(user.Uuid, user.Role.Name)
	if err != nil {
		return nil, err
	}

	resp := map[string]string{
		"token": jwt,
	}

	return resp, nil

}

func (s *UserService) ActivateUser(teacherUuid string) error {

	var teacher models.Teacher
	condition := fmt.Sprintf("uuid = '%s'", teacherUuid)
	if err := s.Repo.First(&teacher, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	var model = models.User{
		Active: true,
	}

	condition = fmt.Sprintf("id = '%d'", teacher.UserID)

	if err := s.Repo.Update(&model, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	return nil
}

func (s *UserService) GetAllClasses() (*[]response.Class, error) {

	var model []models.Class
	if err := s.Repo.FindClasses(&model, ""); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp []response.Class
	for _, item := range model {
		resp = append(resp, response.Class{
			Uuid:              item.Uuid,
			Name:              item.Name,
			Code:              item.Code,
			Teacher:           item.Teacher.User.Name,
			TeacherDepartment: item.Teacher.Department.Name,
		})
	}

	return &resp, nil
}

func (s *UserService) CreateClassRequest(userUuid string, req *request.ClassRequest) error {

	var user models.User
	condition := fmt.Sprintf("uuid = '%s'", userUuid)
	if err := s.Repo.First(&user, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	var class models.Class
	condition = fmt.Sprintf("uuid = '%s'", req.ClassUuid)
	if err := s.Repo.First(&class, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	var classRequest *models.ClassRequest
	condition = fmt.Sprintf("student_id = '%d' AND class_id = '%d'", user.Student.ID, class.ID)
	if err := s.Repo.First(&classRequest, condition); err == nil {
		message := "Anda Sudah Mendaftar Untuk Kelas Ini, Harap Menunggu Konfirmasi Dari Dosen"
		if classRequest.Status == 2 {
			message = "Anda Sudah Terdaftar Di Kelas Ini, Silahkan Akses Dashboard Anda"
		}
		return response.BADREQ_ERR(message)
	}

	model := models.ClassRequest{
		StudentID: user.Student.ID,
		ClassID:   class.ID,
		Status:    1,
	}

	if err := s.Repo.Create(&model); err != nil {
		log.Println(err)
		if utils.IsErrorType(err) {
			return err
		}
		return response.SERVICE_INTERR
	}
	return nil
}

func (s *UserService) JoinClass(userUuid string, req *request.ClassRequest) error {

	var user models.User
	condition := fmt.Sprintf("uuid = '%s'", userUuid)
	if err := s.Repo.First(&user, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	var class models.Class
	condition = fmt.Sprintf("uuid = '%s'", req.ClassUuid)
	if err := s.Repo.First(&class, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	var studentClasses models.StudentClasses
	condition = fmt.Sprintf("student_id = '%d' AND class_id = '%d'", user.Student.ID, class.ID)
	if err := s.Repo.First(&studentClasses, condition); err == nil {
		log.Println(err)
		return response.BADREQ_ERR("Anda Sudah Tergabung Dalam Kelas Tersebut")
	}

	model := models.StudentClasses{
		StudentID: user.Student.ID,
		ClassID:   class.ID,
	}

	if err := s.Repo.Create(&model); err != nil {
		log.Println(err)
		if utils.IsErrorType(err) {
			return err
		}
		return response.SERVICE_INTERR
	}
	return nil
}

func (s *UserService) GetStudentRequestClasses(userUuid string) (*[]response.RequestClass, error) {

	var user models.User
	condition := fmt.Sprintf("uuid = '%s'", userUuid)
	if err := s.Repo.First(&user, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var model []models.ClassRequest
	condition = fmt.Sprintf("student_id = '%d'", user.Student.ID)
	if err := s.Repo.Find(&model, condition, "id"); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp []response.RequestClass
	for _, item := range model {

		var class models.Class
		condition = fmt.Sprintf("id = '%d'", item.Class.ID)
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

		resp = append(resp, response.RequestClass{
			Uuid:      item.Uuid,
			ClassName: class.Name,
			ClassCode: class.Code,
			Teacher:   teacher.User.Name,
			Status:    item.Status,
			CreatedAt: item.CreatedAt,
		})
	}

	return &resp, nil
}

func (s *UserService) GetRequestClasses(userUuid string) (*[]response.RequestClass, error) {

	var user models.User
	condition := fmt.Sprintf("uuid = '%s'", userUuid)
	if err := s.Repo.First(&user, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var model []models.Class
	condition = fmt.Sprintf("teacher_id = '%d'", user.Teacher.ID)
	if err := s.Repo.Find(&model, condition, "id"); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp []response.RequestClass
	for _, item := range model {
		for _, item := range *item.ClassRequests {
			var class models.Class
			condition = fmt.Sprintf("id = '%d'", item.ClassID)
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

			var student models.Student
			condition = fmt.Sprintf("id = '%d'", item.StudentID)
			if err := s.Repo.First(&student, condition); err != nil {
				log.Println(err)
				return nil, response.SERVICE_INTERR
			}

			resp = append(resp, response.RequestClass{
				Uuid:      item.Uuid,
				ClassName: class.Name,
				ClassCode: class.Code,
				Teacher:   teacher.User.Name,
				Student:   student.User.Name,
				Status:    item.Status,
				CreatedAt: item.CreatedAt,
			})
		}

	}

	return &resp, nil
}

func (s *UserService) DashboardAdmin() (*response.DashboardAdmin, error) {

	var departments []models.Department
	if err := s.Repo.Find(&departments, "", ""); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var teacher []models.Teacher
	if err := s.Repo.Find(&teacher, "", ""); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	resp := response.DashboardAdmin{
		JumlahProdi: len(departments),
		JumlahDosen: len(teacher),
	}

	return &resp, nil
}

func (s *UserService) DashboardTeacher(userUuid string) (*response.DashboardTeacher, error) {

	var user models.User
	condition := fmt.Sprintf("uuid = '%s'", userUuid)
	if err := s.Repo.First(&user, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var class []models.Class
	condition = fmt.Sprintf("teacher_id = '%d'", user.Teacher.ID)
	classIDs := s.Repo.SelectColumn(&models.Class{}, condition, "id")
	if err := s.Repo.Find(&class, condition, ""); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var students []models.ClassRequest
	if err := s.Repo.FindTeacherStudents(&students, user.Teacher.ID); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var books []models.Book
	if err := s.Repo.Include(&books, "", "class_id", classIDs); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var materials []models.Material
	if err := s.Repo.Include(&materials, "", "class_id", classIDs); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var videos []models.Video
	if err := s.Repo.Include(&videos, "", "class_id", classIDs); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	resp := response.DashboardTeacher{
		JumlahMahasiswa: len(students),
		JumlahKelas:     len(class),
		JumlahBuku:      len(books),
		JumlahVideo:     len(videos),
		JumlahMateri:    len(materials),
	}

	return &resp, nil
}

func (s *UserService) DashboardStudent(userUuid string) (*response.DashboardStudent, error) {

	var user models.User
	condition := fmt.Sprintf("uuid = '%s'", userUuid)
	if err := s.Repo.First(&user, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var classVerified []models.ClassRequest
	condition = fmt.Sprintf("student_id = '%d' AND status = '%d'", user.Student.ID, 2)
	if err := s.Repo.Find(&classVerified, condition, ""); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var classNot []models.ClassRequest
	condition = fmt.Sprintf("student_id = '%d' AND status = '%d'", user.Student.ID, 0)
	if err := s.Repo.Find(&classNot, condition, ""); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	resp := response.DashboardStudent{
		JumlahKelasVerified: len(classVerified),
		JumlahKelasNot:      len(classNot),
	}

	return &resp, nil
}

func (s *UserService) UpdateStatusClassReq(uuid string, req *request.StatusClassReq) error {

	var class models.ClassRequest
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&class, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	model := models.ClassRequest{
		ID:     class.ID,
		Status: req.Status,
	}

	if err := s.Repo.Update(&model, ""); err != nil {
		log.Println(err)
		if utils.IsErrorType(err) {
			return err
		}
		return response.SERVICE_INTERR
	}
	return nil
}

func (s *UserService) GetAllClassesByStudent(userUuid string) (*[]response.StudentClass, error) {
	var user models.User
	condition := fmt.Sprintf("uuid = '%s'", userUuid)
	if err := s.Repo.First(&user, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var classes []models.Class
	if err := s.Repo.Find(&classes, "", "id"); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var classIDs []int
	condition = fmt.Sprintf("student_id = '%d'", user.Student.ID)

	if err := s.Repo.Pluck(&models.StudentClasses{}, &classIDs, "class_id", condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp []response.StudentClass
	for _, item := range classes {
		var join = false
		if contain := utils.CheckContainsInt(classIDs, int(item.ID)); contain {
			join = true
		}

		var teacher models.Teacher
		condition := fmt.Sprintf("id = '%d'", item.TeacherID)
		if err := s.Repo.First(&teacher, condition); err != nil {
			log.Println(err)
			return nil, response.SERVICE_INTERR
		}

		resp = append(resp, response.StudentClass{
			Join: join,
			Class: &response.Class{
				Uuid:              item.Uuid,
				Name:              item.Name,
				Code:              item.Code,
				Teacher:           teacher.User.Name,
				TeacherDepartment: teacher.Department.Name,
			},
		})
	}

	return &resp, nil
}
