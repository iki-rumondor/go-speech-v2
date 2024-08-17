package repositories

import (
	"fmt"

	"github.com/iki-rumondor/go-speech/internal/domain/layers/interfaces"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserInterface(db *gorm.DB) interfaces.UserInterface {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Create(model interface{}) error {
	return r.db.Create(model).Error
}

func (r *UserRepo) Find(dest interface{}, condition, order string) error {
	return r.db.Preload(clause.Associations).Order(order).Find(dest, condition).Error
}

func (r *UserRepo) First(dest interface{}, condition string) error {
	return r.db.Preload(clause.Associations).First(dest, condition).Error
}

func (r *UserRepo) Update(model interface{}, condition string) error {
	return r.db.Where(condition).Updates(model).Error
}

func (r *UserRepo) Delete(data interface{}, withAssociation []string) error {
	return r.db.Select(withAssociation).Delete(data).Error
}

func (r *UserRepo) Distinct(model interface{}, column, condition string, dest *[]string) error {
	return r.db.Model(model).Distinct().Where(condition).Pluck(column, dest).Error
}

func (r *UserRepo) Truncate(tableName string) error {
	return r.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error
}

func (r *UserRepo) Pluck(model, dest interface{}, columnName, condition string) error {
	return r.db.Model(model).Where(condition).Pluck(columnName, dest).Error
}

func (r *UserRepo) FindClasses(model *[]models.Class, condition string) error {
	return r.db.Preload(clause.Associations).Preload("Teacher.User").Preload("Teacher.Department").Find(model, condition).Error
}

func (r *UserRepo) FindTeacherStudents(model *[]models.ClassRequest, teacherID uint) error {
	classIDs := r.db.Model(&models.Class{}).Where("teacher_id = ?", teacherID).Select("id")
	return r.db.Find(model, "class_id IN (?) AND status = ?", classIDs, 2).Error

}

func (r *UserRepo) SelectColumn(model interface{}, condition, columnName string) *gorm.DB {
	return r.db.Model(model).Where(condition).Select(columnName)
}

func (r *UserRepo) Include(model interface{}, condition, colName string, selectCols *gorm.DB) error {
	conds := fmt.Sprintf("%s IN (?)", colName)
	return r.db.Where(condition).Find(model, conds, selectCols).Error
}

func (r *UserRepo) FindStudentClasses(userUuid string) (*[]int, error) {
	var user models.User
	if err := r.db.Preload("Student").First(&user, "uuid = ?", userUuid).Error; err != nil {
		return nil, err
	}

	var studentClassesIDs []int
	if err := r.db.Model(&models.StudentClasses{}).Where("student_id = ?", user.Student.ID).Pluck("class_id", &studentClassesIDs).Error; err != nil {
		return nil, err
	}

	return &studentClassesIDs, nil
}

//	func (r *UserRepo) FindTeacherClassReq(dest *[]models.ClassRequest, teacherID uint) error {
//		subQuery := r.db.Model(&models.Class{}).Where("academic_year_id = ?", yearID).Select("facility_id")
//		return r.db.Preload(clause.Associations).Preload("Teacher.User").Preload("Teacher.Department").Find(model, condition).Error
//	}

func (r *UserRepo) FindClassNotifications(userUuid string) (*[]models.ClassNotification, error) {
	var user models.User
	if err := r.db.Preload("Student").First(&user, "uuid = ?", userUuid).Error; err != nil {
		return nil, err
	}

	var classIDs []uint
	if err := r.db.Model(&models.StudentClasses{}).Where("student_id = ?", user.Student.ID).Pluck("class_id", &classIDs).Error; err != nil {
		return nil, err
	}

	var notifications []models.ClassNotification
	if err := r.db.Order("id DESC").Find(&notifications, "class_id IN (?)", classIDs).Error; err != nil {
		return nil, err
	}

	return &notifications, nil
}

func (r *UserRepo) ReadNotification(userUuid, notificationUuid string) error {
	var user models.User
	if err := r.db.Preload("Student").First(&user, "uuid = ?", userUuid).Error; err != nil {
		return err
	}

	var notification models.ClassNotification
	if err := r.db.First(&notification, "uuid = ?", notificationUuid).Error; err != nil {
		return err
	}

	if row := r.db.First(&models.ReadNotification{}, "student_id = ? AND notification_id = ?", user.Student.ID, notification.ID).RowsAffected; row != 0 {
		return fmt.Errorf("notifikasi sudah dibaca")
	}

	model := models.ReadNotification{
		NotificationID: notification.ID,
		StudentID:      user.Student.ID,
	}

	return r.db.Create(&model).Error
}

func (r *UserRepo) GetReadNotifications(userUuid string) (*[]models.ReadNotification, error) {
	var user models.User
	if err := r.db.Preload("Student").First(&user).Error; err != nil {
		return nil, err
	}

	var resp []models.ReadNotification
	if err := r.db.Find(&resp, "student_id = ?", user.Student.ID).Error; err != nil {
		return nil, err
	}

	return &resp, nil
}

func (r *UserRepo) GetUnreadNotification(userUuid string) (*[]models.ClassNotification, error) {
	var user models.User
	if err := r.db.Preload("Student").First(&user).Error; err != nil {
		return nil, err
	}

	var notificationIDs []uint
	if err := r.db.Model(&models.ReadNotification{}).Pluck("notification_id", &notificationIDs).Error; err != nil {
		return nil, err
	}

	var notifications []models.ClassNotification
	if err := r.db.Find(&notifications, "id NOT IN (?)", notificationIDs).Error; err != nil {
		return nil, err
	}

	return &notifications, nil
}
