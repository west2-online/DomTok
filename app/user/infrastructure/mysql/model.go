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

package mysql

import "github.com/west2-online/DomTok/pkg/constants"

// User 是 mysql 【独有】的，和 db 中的表数据一一对应，和 entities 层的 User 的作用域不一样
type User struct {
	// model    gorm.Model
	ID       int64
	Username string
	Password string
	Email    string
	Phone    string
}

func (User) TableName() string {
	return constants.UserTableName
}
