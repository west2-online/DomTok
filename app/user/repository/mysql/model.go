package mysql

import (
	"gorm.io/gorm"
)

// User 是 mysql 【独有】的，和 db 中的表数据一一对应，和 entities 层的 User 的作用域不一样
type User struct {
	model    gorm.Model
	UserName string
	Password string
	Email    string
}

func (User) TableName() string {
	return "users"
}
