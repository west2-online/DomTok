/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package entities

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// User 用于在 handler ---use case --- repository 之间传递数据的实体类
// 目的是方便 use case 操作对应的业务
type User struct {
	Uid      int64
	UserName string
	Password string
	Email    string
	Phone    string
	// AvatarURL string
}

// IsValidEmail 根据正则匹配来判断是否是合法的邮箱
func (u *User) IsValidEmail() bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(u.Email)
}

// EncryptPassword 加密密码, 直接修改 User 中的 Password 字段
func (u *User) EncryptPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) IsValidPassword() bool {
	// 密码长度至少为8个字符
	if len(u.Password) < ValidPasswordLength {
		return false
	}

	// 至少包含一个大写字母、一个小写字母、一个数字和一个特殊字符
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range u.Password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",<>\.\?/\\|]`).MatchString(string(char)):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}
