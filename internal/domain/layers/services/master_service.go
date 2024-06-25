package services

import (
	"fmt"
	"log"

	"github.com/iki-rumondor/go-speech/internal/domain/layers/interfaces"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/models"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/request"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/response"
	"github.com/iki-rumondor/go-speech/internal/utils"
)

type MasterService struct {
	Repo interfaces.MasterInterface
}

func NewMasterService(repo interfaces.MasterInterface) *MasterService {
	return &MasterService{
		Repo: repo,
	}
}

func (s *MasterService) CreateClass(userUuid string, req *request.Class) error {
	var user models.User
	condition := fmt.Sprintf("uuid = '%s'", userUuid)
	if err := s.Repo.First(&user, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	model := models.Class{
		TeacherID: user.Teacher.ID,
		Name:      req.Name,
		Code:      req.Code,
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

func (s *MasterService) DeleteClass(userUuid, classUuid string) error {

	var user models.User
	condition := fmt.Sprintf("uuid = '%s'", userUuid)
	if err := s.Repo.First(&user, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	var model models.Class
	condition = fmt.Sprintf("teacher_id = '%d' AND uuid = '%s'", user.Teacher.ID, classUuid)
	if err := s.Repo.First(&model, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	if err := s.Repo.Delete(&model, nil); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}
	return nil
}

func (s *MasterService) CreateDepartment(req *request.Department) error {

	model := models.Department{
		Name: req.Name,
	}

	if err := s.Repo.Create(&model); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}
	return nil
}

func (s *MasterService) UpdateDepartment(uuid string, req *request.Department) error {

	model := models.Department{
		Name: req.Name,
	}

	condition := fmt.Sprintf("uuid = '%s'", uuid)

	if err := s.Repo.Update(&model, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}
	return nil
}

func (s *MasterService) DeleteDepartment(uuid string) error {

	var model models.Department
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&model, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	if err := s.Repo.Delete(&model, nil); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}
	return nil
}

func (s *MasterService) GetAllDepartment() (*[]response.Department, error) {

	var model []models.Department

	if err := s.Repo.Find(&model, "", "id"); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp []response.Department
	for _, item := range model {
		resp = append(resp, response.Department{
			Uuid: item.Uuid,
			Name: item.Name,
		})
	}

	return &resp, nil
}

func (s *MasterService) GetClasses(userUuid string) (*[]response.Class, error) {
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

	var resp []response.Class
	for _, item := range model {
		resp = append(resp, response.Class{
			Uuid: item.Uuid,
			Name: item.Name,
			Code: item.Code,
		})
	}

	return &resp, nil
}

func (s *MasterService) GetAllClasses() (*[]response.Class, error) {

	var model []models.Class
	if err := s.Repo.Find(&model, "", "id"); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp []response.Class
	for _, item := range model {
		var teacher models.Teacher
		condition := fmt.Sprintf("id = '%d'", item.TeacherID)
		if err := s.Repo.First(&teacher, condition); err != nil {
			return nil, response.SERVICE_INTERR
		}

		resp = append(resp, response.Class{
			Uuid:              item.Uuid,
			Name:              item.Name,
			Code:              item.Code,
			Teacher:           teacher.User.Name,
			TeacherDepartment: teacher.Department.Name,
		})
	}

	return &resp, nil
}

func (s *MasterService) GetStudentClasses(userUuid string) (*[]response.Class, error) {
	var user models.User
	condition := fmt.Sprintf("uuid = '%s'", userUuid)
	if err := s.Repo.First(&user, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var model []models.Class
	if err := s.Repo.FindStudentClasses(user.Student.ID, &model); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp []response.Class
	for _, item := range model {
		resp = append(resp, response.Class{
			Uuid: item.Uuid,
			Name: item.Name,
			Code: item.Code,
		})
	}

	return &resp, nil
}

func (s *MasterService) GetDepartment(uuid string) (*response.Department, error) {

	var model models.Department
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&model, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp = response.Department{
		Uuid: model.Uuid,
		Name: model.Name,
	}

	return &resp, nil
}

func (s *MasterService) GetClass(uuid string) (*response.Class, error) {

	var model models.Class
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&model, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var classReq []models.ClassRequest
	condition = fmt.Sprintf("class_id = '%d' AND status = '%d'", model.ID, 2)
	if err := s.Repo.Find(&classReq, condition, "id"); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var students []response.Student
	for _, item := range classReq {
		var student models.Student
		condition = fmt.Sprintf("id = '%d'", item.StudentID)
		if err := s.Repo.First(&student, condition); err != nil {
			log.Println(err)
			return nil, response.SERVICE_INTERR
		}

		students = append(students, response.Student{
			Uuid:              student.Uuid,
			Nim:               student.Nim,
			Name:              student.User.Name,
			Email:             student.User.Email,
			RegisterClassTime: item.CreatedAt,
		})
	}

	var resp = response.Class{
		Uuid:     model.Uuid,
		Name:     model.Name,
		Code:     model.Code,
		Students: &students,
	}

	return &resp, nil
}

func (s *MasterService) UpdateClass(userUuid, uuid string, req *request.Class) error {
	var user models.User
	condition := fmt.Sprintf("uuid = '%s'", userUuid)
	if err := s.Repo.First(&user, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	var class models.Class
	condition = fmt.Sprintf("uuid = '%s' AND teacher_id = '%d'", uuid, user.Teacher.ID)
	if err := s.Repo.First(&class, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	model := models.Class{
		ID:   class.ID,
		Name: req.Name,
		Code: req.Code,
	}

	if err := s.Repo.Update(&model, ""); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}
	return nil
}

func (s *MasterService) CreateNote(req *request.Note) error {
	var class models.Class
	condition := fmt.Sprintf("uuid = '%s'", req.ClassUuid)
	if err := s.Repo.First(&class, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	model := models.Note{
		Title:   req.Title,
		Body:    req.Body,
		ClassID: class.ID,
	}

	if err := s.Repo.Create(&model); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}
	return nil
}

func (s *MasterService) GetNotes(classUuid string) (*[]response.Note, error) {

	var class models.Class
	condition := fmt.Sprintf("uuid = '%s'", classUuid)
	if err := s.Repo.First(&class, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var notes []models.Note
	condition = fmt.Sprintf("class_id = '%d'", class.ID)
	if err := s.Repo.Find(&notes, condition, "created_at desc"); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp []response.Note
	for _, item := range notes {
		resp = append(resp, response.Note{
			Uuid:      item.Uuid,
			Title:     item.Title,
			Body:      item.Body,
			CreatedAt: item.CreatedAt,
		})
	}

	return &resp, nil
}

func (s *MasterService) GetNote(uuid string) (*response.Note, error) {

	var note models.Note
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&note, condition); err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp = response.Note{
		Uuid:      note.Uuid,
		Title:     note.Title,
		Body:      note.Body,
		CreatedAt: note.CreatedAt,
	}

	return &resp, nil
}

func (s *MasterService) UpdateNote(uuid string, req *request.Note) error {
	var note models.Note
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&note, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	model := models.Note{
		ID:    note.ID,
		Title: req.Title,
		Body:  req.Body,
	}

	if err := s.Repo.Update(&model, ""); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}
	return nil
}

func (s *MasterService) DeleteNote(uuid string) error {
	var note models.Note
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&note, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	if err := s.Repo.Delete(&note, nil); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}
	return nil
}
