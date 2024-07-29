package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/iki-rumondor/go-speech/internal/consts"
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

		t := time.Unix(item.CreatedAt/1000, 0)
		formattedDate := t.Format("02-01-2006")

		students = append(students, response.Student{
			Uuid:              student.Uuid,
			Nim:               student.Nim,
			Name:              student.User.Name,
			Email:             student.User.Username,
			RegisterClassTime: item.CreatedAt,
			RegTimeString:     formattedDate,
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

func (s *MasterService) GetClassesReport(uuid string) error {

	class, err := s.GetClass(uuid)
	if err != nil {
		return err
	}

	if err := s.Repo.LaravelClassReport(class); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}
	return nil
}

func (s *MasterService) CreateMaterial(req *request.Material, videoName, bookName string) error {
	var class models.Class
	condition := fmt.Sprintf("uuid = '%s'", req.ClassUuid)
	if err := s.Repo.First(&class, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	subtitleName, err := s.Assembly_VideoToSubtitle(videoName)
	if err != nil {
		log.Println(err)
		return &response.Error{
			Code:    500,
			Message: "Gagal Konversi Video To Audio AssemblyAI",
		}
	}

	model := models.Material{
		ClassID:     class.ID,
		Title:       req.Title,
		Description: req.Description,
		Video: &models.VideoPart{
			VideoName:    videoName,
			SubtitleName: *subtitleName,
		},
		Book: &models.BookPart{
			FileName: bookName,
		},
	}

	if err := s.Repo.Create(&model); err != nil {
		DeleteSubtitle(*subtitleName)
		return response.SERVICE_INTERR
	}

	return nil
}

func DeleteSubtitle(subtitleName string) {
	if err := os.Remove(filepath.Join(consts.SUBTITLE_FOLDER, subtitleName)); err != nil {
		log.Println(err)
	}
}

func (s *MasterService) GetAllMaterials(classUuid string) (*[]response.Material, error) {
	var class models.Class
	condition := fmt.Sprintf("uuid = '%s'", classUuid)
	if err := s.Repo.First(&class, condition); err != nil {
		return nil, response.SERVICE_INTERR
	}

	var teacher models.Teacher
	condition = fmt.Sprintf("id = '%d'", class.TeacherID)
	if err := s.Repo.First(&teacher, condition); err != nil {
		return nil, response.SERVICE_INTERR
	}

	var materials []models.Material
	condition = fmt.Sprintf("class_id = '%d'", class.ID)
	if err := s.Repo.Find(&materials, condition, "created_at ASC"); err != nil {
		return nil, response.SERVICE_INTERR
	}

	var resp []response.Material
	for _, item := range materials {
		resp = append(resp, response.Material{
			Uuid:         item.Uuid,
			Title:        item.Title,
			Description:  item.Description,
			BookName:     item.Book.FileName,
			VideoName:    item.Video.VideoName,
			VideoUuid:    item.Video.Uuid,
			SubtitleName: item.Video.SubtitleName,
			CreatedAt:    item.CreatedAt,
			UpdateAt:     item.UpdatedAt,
			Class: &response.Class{
				Name:    class.Name,
				Code:    class.Code,
				Teacher: teacher.User.Name,
			},
		})
	}

	return &resp, nil
}

func (s *MasterService) GetMaterial(uuid string) (*response.Material, error) {
	var material models.Material
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&material, condition); err != nil {
		return nil, response.SERVICE_INTERR
	}

	var teacher models.Teacher
	condition = fmt.Sprintf("id = '%d'", material.Class.TeacherID)
	if err := s.Repo.First(&teacher, condition); err != nil {
		return nil, response.SERVICE_INTERR
	}

	var resp = response.Material{
		Uuid:         material.Uuid,
		Title:        material.Title,
		Description:  material.Description,
		BookName:     material.Book.FileName,
		VideoName:    material.Video.VideoName,
		VideoUuid:    material.Video.Uuid,
		SubtitleName: material.Video.SubtitleName,
		CreatedAt:    material.CreatedAt,
		UpdateAt:     material.UpdatedAt,
		Class: &response.Class{
			Name:    material.Class.Name,
			Code:    material.Class.Code,
			Teacher: teacher.User.Name,
		},
	}

	return &resp, nil
}

func (s *MasterService) UpdateMaterial(uuid string, req *request.UpdateMaterial) error {
	var material models.Material
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&material, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	model := models.Material{
		ID:          material.ID,
		Title:       req.Title,
		Description: req.Description,
	}

	if err := s.Repo.Update(&model, ""); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}
	return nil
}

func (s *MasterService) DeleteMaterial(uuid string) error {
	var material models.Material
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&material, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	if err := s.Repo.Delete(&material, []string{"Video", "Book"}); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	subtitleFolder := "internal/files/subtitle"
	subtitlePath := filepath.Join(subtitleFolder, material.Video.SubtitleName)

	videoFolder := "internal/files/videos"
	videoPath := filepath.Join(videoFolder, material.Video.VideoName)

	bookFolder := "internal/files/books"
	bookPath := filepath.Join(bookFolder, material.Book.FileName)

	if err := os.Remove(subtitlePath); err != nil {
		log.Println(err.Error())
	}

	if err := os.Remove(videoPath); err != nil {
		log.Println(err.Error())
	}

	if err := os.Remove(bookPath); err != nil {
		log.Println(err.Error())
	}

	return nil
}

func (s *MasterService) GetTeacher(uuid string) (*response.Teacher, error) {
	var teacher models.Teacher
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&teacher, condition); err != nil {
		return nil, response.SERVICE_INTERR
	}

	var resp = response.Teacher{
		Uuid:           teacher.Uuid,
		Nidn:           teacher.Nidn,
		Department:     teacher.Department.Name,
		DepartmentUuid: teacher.Department.Uuid,
		Email:          teacher.User.Username,
		Name:           teacher.User.Name,
		Active:         teacher.User.Active,
	}

	return &resp, nil
}

func (s *MasterService) UpdateTeacher(uuid string, req *request.UpdateTeacher) error {
	var department models.Department
	condition := fmt.Sprintf("uuid = '%s'", req.DepartmentUuid)
	if err := s.Repo.First(&department, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	var teacher models.Teacher
	condition = fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&teacher, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	teacherModel := models.Teacher{
		ID:           teacher.ID,
		Nidn:         req.Nidn,
		DepartmentID: department.ID,
	}

	userModel := models.User{
		ID:       teacher.UserID,
		Name:     req.Name,
		Password: req.Nidn,
		Username: req.Nidn,
	}

	if err := s.Repo.UpdateTeacher(&teacherModel, &userModel); err != nil {
		log.Println(err)
		if utils.IsErrorType(err) {
			return err
		}
		return response.SERVICE_INTERR
	}
	return nil
}

func (s *MasterService) DeleteTeacher(uuid string) error {
	var teacher models.Teacher
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&teacher, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	if err := s.Repo.Delete(teacher.User, []string{"Teacher"}); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}
	return nil
}

func (s *MasterService) GetStudents() (*[]response.Student, error) {

	var students []models.Student
	if err := s.Repo.Find(&students, "", "created_at ASC"); err != nil {
		return nil, response.SERVICE_INTERR
	}

	var resp []response.Student
	for _, item := range students {
		resp = append(resp, response.Student{
			Uuid: item.Uuid,
			Name: item.User.Name,
			Nim:  item.Nim,
		})
	}

	return &resp, nil
}

func (s *MasterService) GetStudent(uuid string) (*response.Student, error) {
	var student models.Student
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&student, condition); err != nil {
		return nil, response.SERVICE_INTERR
	}

	var resp = response.Student{
		Uuid:  student.Uuid,
		Nim:   student.Nim,
		Email: student.User.Username,
		Name:  student.User.Name,
	}

	return &resp, nil
}

func (s *MasterService) UpdateStudent(uuid string, req *request.UpdateStudent) error {

	var student models.Student
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&student, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	studentModel := models.Student{
		ID:  student.ID,
		Nim: req.Nim,
	}

	userModel := models.User{
		ID:       student.UserID,
		Name:     req.Name,
		Password: req.Nim,
		Username: req.Nim,
	}

	if err := s.Repo.UpdateStudent(&studentModel, &userModel); err != nil {
		log.Println(err)
		if utils.IsErrorType(err) {
			return err
		}
		return response.SERVICE_INTERR
	}
	return nil
}

func (s *MasterService) DeleteStudent(uuid string) error {
	var student models.Student
	condition := fmt.Sprintf("uuid = '%s'", uuid)
	if err := s.Repo.First(&student, condition); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	if err := s.Repo.Delete(student.User, []string{"Student"}); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}
	return nil
}
