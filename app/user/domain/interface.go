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

package domain

import "context"

// UserDB 表示持久化存储接口 (或者也可以叫做 DBPort)
type UserDB interface {
	IsUserExist(ctx context.Context, username string) (bool, error)
	CreateUser(ctx context.Context, user *User) error
}

type UserUseCase interface {
	RegisterUser(ctx context.Context, user *User) (uid int64, err error)
	Login(ctx context.Context, user *User) (*User, error)
}
