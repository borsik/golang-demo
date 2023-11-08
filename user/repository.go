package user

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	Insert(user User) (User, error)
	Select(name string, country string, offset int, limit int) ([]User, int64)
	SelectByID(id string) (User, error)
	Update(id string, input InputUser) error
	Delete(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Insert(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *repository) Select(name string, country string, offset int, limit int) ([]User, int64) {
	var users []User
	var totalCount int64
	db := r.db
	if name != "" {
		db = db.Where("first_name ILIKE ?", fmt.Sprintf("%%%s%%", name)).Or("last_name ILIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	if country != "" {
		db = db.Where("country = ?", strings.ToUpper(country))
	}
	db.Model(&users).Count(&totalCount)
	db.Limit(limit).Offset(offset).Find(&users)
	return users, totalCount
}

func (r *repository) SelectByID(id string) (User, error) {
	var user User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *repository) Update(id string, input InputUser) error {
	values := map[string]interface{}{
		"first_name": input.FirstName,
		"last_name":  input.LastName,
		"nickname":   input.Nickname,
		"password":   input.Password,
		"email":      input.Email,
		"country":    input.Country,
	}
	err := r.db.Model(&User{ID: id}).Updates(values).Error
	return err
}

func (r *repository) Delete(id string) error {
	err := r.db.Delete(&User{ID: id}).Error
	return err
}
