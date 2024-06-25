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

// func (r *UserRepo) FindTeacherClassReq(dest *[]models.ClassRequest, teacherID uint) error {
// 	subQuery := r.db.Model(&models.Class{}).Where("academic_year_id = ?", yearID).Select("facility_id")
// 	return r.db.Preload(clause.Associations).Preload("Teacher.User").Preload("Teacher.Department").Find(model, condition).Error
// }
