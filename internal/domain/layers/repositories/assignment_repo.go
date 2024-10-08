package repositories

import (
	"fmt"
	"log"

	"github.com/iki-rumondor/go-speech/internal/domain/layers/interfaces"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/models"
	"github.com/iki-rumondor/go-speech/internal/utils"
	"gorm.io/gorm"
)

type AssignmentRepo struct {
	db *gorm.DB
}

func NewAssignmentInterface(db *gorm.DB) interfaces.AssignmentInterface {
	return &AssignmentRepo{
		db: db,
	}
}

func (r *AssignmentRepo) CreateAssignment(classUuid string, model *models.Assignment) error {
	var class models.Class
	if err := r.db.Preload("Teacher.User").First(&class, "uuid = ?", classUuid).Error; err != nil {
		return err
	}
	model.ClassID = class.ID

	if err := r.db.Create(model).Error; err != nil {
		return err
	}

	notification := models.ClassNotification{
		ClassID: class.ID,
		Title:   "Tugas Baru",
		Body:    fmt.Sprintf("Tugas telah ditambahkan oleh %s dengan deadline pada tanggal %s", class.Teacher.User.Name, utils.UnixToDate(model.Deadline)),
	}

	if err := r.db.Create(&notification).Error; err != nil {
		log.Println(err)
	}

	return nil
}

func (r *AssignmentRepo) DeleteAssignment(uuid string) error {
	var assignment models.Assignment
	if err := r.db.First(&assignment, "uuid = ?", uuid).Error; err != nil {
		return err
	}
	return r.db.Select("Answers").Delete(&assignment).Error
}

func (r *AssignmentRepo) FindAssignmentsByClass(classUuid string) (*[]models.Assignment, error) {
	var assignments []models.Assignment
	if err := r.db.Joins("Class").Preload("Class.Teacher.User").Preload("Answers.User.Student").Find(&assignments, "Class.uuid = ?", classUuid).Error; err != nil {
		return nil, err
	}
	return &assignments, nil
}

func (r *AssignmentRepo) FindAssignmentsByUser(userUuid string) (*[]models.Assignment, error) {
	var user models.User
	if err := r.db.Preload("Class.Teacher.User").First(&user, "uuid = ?", userUuid).Error; err != nil {
		return nil, err
	}

	var studentClasses []uint
	if err := r.db.Model(&models.StudentClasses{}).Where("student_id = ?", user.Student.ID).Pluck("class_id", &studentClasses).Error; err != nil {
		return nil, err
	}

	var assignments []models.Assignment
	if err := r.db.Find(&assignments, "class_id IN (?)", studentClasses).Error; err != nil {
		return nil, err
	}
	return &assignments, nil
}

func (r *AssignmentRepo) FindAssignmentsByStudent(studentUuid string) (*[]models.Assignment, error) {
	var student models.Student
	if err := r.db.First(&student, "uuid = ?", studentUuid).Error; err != nil {
		return nil, err
	}

	var studentClasses []uint
	if err := r.db.Model(&models.StudentClasses{}).Where("student_id = ?", student.ID).Pluck("class_id", &studentClasses).Error; err != nil {
		return nil, err
	}

	var assignments []models.Assignment
	if err := r.db.Preload("Class").Find(&assignments, "class_id IN (?)", studentClasses).Error; err != nil {
		return nil, err
	}
	return &assignments, nil
}

func (r *AssignmentRepo) FirstAssignmentByUuid(uuid string) (*models.Assignment, error) {
	var assignment models.Assignment
	if err := r.db.Preload("Answers").First(&assignment, "uuid = ?", uuid).Error; err != nil {
		return nil, err
	}
	return &assignment, nil

}
func (r *AssignmentRepo) FirstStudentByUuid(uuid string) (*models.Student, error) {
	var student models.Student
	if err := r.db.Preload("User").First(&student, "uuid = ?", uuid).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *AssignmentRepo) UpdateAssignment(model *models.Assignment) error {
	return r.db.Updates(model).Error
}

func (r *AssignmentRepo) FirstAnswerByUser(userUuid, assignmentUuid string) (*models.Answer, error) {
	var user models.User
	if err := r.db.First(&user, "uuid = ?", userUuid).Error; err != nil {
		return nil, err
	}

	var assignment models.Assignment
	if err := r.db.First(&assignment, "uuid = ?", assignmentUuid).Error; err != nil {
		return nil, err
	}

	var answer models.Answer
	if err := r.db.First(&answer, "user_id = ? AND assignment_id = ?", user.ID, assignment.ID).Error; err != nil {
		return nil, err
	}

	return &answer, nil
}

func (r *AssignmentRepo) CreateAnswer(assignmentUuid, userUuid string, model *models.Answer) error {
	var assignment models.Assignment
	if err := r.db.First(&assignment, "uuid = ?", assignmentUuid).Error; err != nil {
		return err
	}

	var user models.User
	if err := r.db.First(&user, "uuid = ?", userUuid).Error; err != nil {
		return err
	}

	model.AssignmentID = assignment.ID
	model.UserID = user.ID

	return r.db.Create(model).Error
}

func (r *AssignmentRepo) UpdateAnswer(answerUuid string, model *models.Answer) error {
	var answer models.Answer
	if err := r.db.First(&answer, "uuid = ?", answerUuid).Error; err != nil {
		return err
	}
	model.ID = answer.ID

	return r.db.Updates(model).Error
}
