package mysql

import (
	"gorm.io/gorm"
)

type User struct {
	model    gorm.Model
	UserName string
	Password string
	Email    string
}

func (User) TableName() string {
	return "users"
}
