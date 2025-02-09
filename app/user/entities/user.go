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

// IsValidEmail TODO: 根据正则匹配 ？来判断是否是合法的邮箱
func (u *User) IsValidEmail() bool {
	return true
}

// EncryptPassword TODO: 加密密码, 直接修改 User 中的 Password 字段
func (u *User) EncryptPassword() error {
	return nil
}

// CheckPassword TODO: 检查密码是否正确
func (u *User) CheckPassword() bool {
	return true
}
