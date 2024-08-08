package services

import (
	"errors"
	"log"

	"github.com/iki-rumondor/go-speech/internal/domain/layers/interfaces"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/models"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/request"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/response"
	"github.com/iki-rumondor/go-speech/internal/utils"
	"gorm.io/gorm"
)

type AssignmentService struct {
	Repo interfaces.AssignmentInterface
}

func NewAssignmentService(repo interfaces.AssignmentInterface) *AssignmentService {
	return &AssignmentService{
		Repo: repo,
	}
}

func (s *AssignmentService) CreateAssignment(body *request.Assignment) error {

	model := models.Assignment{
		Title:       body.Title,
		Description: body.Description,
		Deadline:    body.Deadline,
	}

	if err := s.Repo.CreateAssignment(body.ClassUuid, &model); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	return nil
}

func (s *AssignmentService) FindAssignmentByClass(classUuid string) (*[]response.Assignment, error) {

	assignments, err := s.Repo.FindAssignmentsByClass(classUuid)
	if err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp []response.Assignment
	for _, item := range *assignments {
		var answers []response.Answer
		if item.Answers != nil {
			for _, item := range *item.Answers {
				answers = append(answers, response.Answer{
					Uuid:     item.Uuid,
					Ontime:   item.Ontime,
					Grade:    item.Grade,
					Filename: item.Filename,
					Student: &response.Student{
						Name: item.User.Name,
						Nim:  item.User.Student.Nim,
					},
				})
			}
		}

		resp = append(resp, response.Assignment{
			Uuid:        item.Uuid,
			Title:       item.Title,
			Description: item.Description,
			Deadline:    item.Deadline,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
			Class: &response.Class{
				Teacher: item.Class.Teacher.User.Name,
			},
			Answers: &answers,
		})
	}

	return &resp, nil
}

func (s *AssignmentService) FindAssignmentByStudent(studentUuid string) (*[]response.Assignment, error) {

	student, err := s.Repo.FirstStudentByUuid(studentUuid)
	if err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	assignments, err := s.Repo.FindAssignmentsByStudent(studentUuid)
	if err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp []response.Assignment
	for _, item := range *assignments {
		answer, err := s.Repo.FirstAnswerByUser(student.User.Uuid, item.Uuid)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		var studentAnswer response.Answer
		if answer != nil {
			studentAnswer = response.Answer{
				Uuid:      answer.Uuid,
				Ontime:    answer.Ontime,
				Filename:  answer.Filename,
				Grade:     answer.Grade,
				Submitted: true,
				Student: &response.Student{
					Name: student.User.Name,
					Nim:  student.Nim,
				},
			}
		}

		resp = append(resp, response.Assignment{
			Uuid:          item.Uuid,
			Title:         item.Title,
			Description:   item.Description,
			Deadline:      item.Deadline,
			CreatedAt:     item.CreatedAt,
			UpdatedAt:     item.UpdatedAt,
			StudentAnswer: &studentAnswer,
			Class: &response.Class{
				Name: item.Class.Name,
			},
		})
	}

	return &resp, nil
}

func (s *AssignmentService) FindAssignmentsByUser(userUuid, classUuid string) (*[]response.Assignment, error) {

	assignments, err := s.Repo.FindAssignmentsByClass(classUuid)
	if err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp []response.Assignment
	for _, item := range *assignments {
		answer, err := s.Repo.FirstAnswerByUser(userUuid, item.Uuid)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		var studentAnswer response.Answer
		if answer != nil {
			studentAnswer = response.Answer{
				Uuid:      answer.Uuid,
				Ontime:    answer.Ontime,
				Filename:  answer.Filename,
				Grade:     answer.Grade,
				Submitted: true,
			}
		}

		resp = append(resp, response.Assignment{
			Uuid:        item.Uuid,
			Title:       item.Title,
			Description: item.Description,
			Deadline:    item.Deadline,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
			Class: &response.Class{
				Teacher: item.Class.Teacher.User.Name,
			},
			StudentAnswer: &studentAnswer,
		})
	}

	return &resp, nil
}

func (s *AssignmentService) FirstAssignmentByUuid(uuid string) (*response.Assignment, error) {

	assignment, err := s.Repo.FirstAssignmentByUuid(uuid)
	if err != nil {
		log.Println(err)
		return nil, response.SERVICE_INTERR
	}

	var resp = response.Assignment{
		Uuid:        assignment.Uuid,
		Title:       assignment.Title,
		Description: assignment.Description,
		CreatedAt:   assignment.CreatedAt,
		UpdatedAt:   assignment.UpdatedAt,
	}

	return &resp, nil
}

func (s *AssignmentService) UpdateAssignment(uuid string, body *request.UpdateAssignment) error {

	assignment, err := s.Repo.FirstAssignmentByUuid(uuid)
	if err != nil {
		return response.SERVICE_INTERR
	}

	model := models.Assignment{
		ID:          assignment.ID,
		Title:       body.Title,
		Description: body.Description,
	}

	if err := s.Repo.UpdateAssignment(&model); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	return nil
}

func (s *AssignmentService) DeleteAssignment(uuid string) error {

	if err := s.Repo.DeleteAssignment(uuid); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	return nil
}

func (s *AssignmentService) CreateAnswer(assignmentUuid, userUuid, fileName string) error {

	answer, err := s.Repo.FirstAnswerByUser(userUuid, assignmentUuid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if answer != nil {
		return response.BADREQ_ERR("Anda Sudah Melakukan Pengumpulan Tugas")
	}

	assignment, err := s.Repo.FirstAssignmentByUuid(assignmentUuid)
	if err != nil {
		return response.SERVICE_INTERR
	}

	var ontime = true
	if utils.IsAfterUnix(assignment.Deadline) {
		ontime = false
	}

	model := models.Answer{
		Filename: fileName,
		Ontime:   ontime,
	}

	if err := s.Repo.CreateAnswer(assignmentUuid, userUuid, &model); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	return nil
}

func (s *AssignmentService) GradeAnswer(answerUuid string, body *request.Grading) error {

	model := models.Answer{
		Grade: body.Grade,
	}

	if err := s.Repo.UpdateAnswer(answerUuid, &model); err != nil {
		log.Println(err)
		return response.SERVICE_INTERR
	}

	return nil
}
