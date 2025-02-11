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

package service

import (
	"regexp"

	"github.com/west2-online/DomTok/app/user/domain/repository"
	"github.com/west2-online/DomTok/pkg/utils"
)

type UserService struct {
	db      repository.UserDB
	sf      *utils.Snowflake
	emailRe *regexp.Regexp
}

// NewUserService 返回一个 NewUserService 实例
func NewUserService(db repository.UserDB, sf *utils.Snowflake) *UserService {
	if db == nil {
		panic("userService`s db should not be nil")
	}
	if sf == nil {
		panic("userService`s sf should not be nil")
	}

	svc := &UserService{db: db}
	svc.init()

	return svc
}

func (svc *UserService) init() {
	svc.initEmailRe()
}

// 理论上来说这也是需要外部注入以及常量管理规则的, 但我懒得搞了, 就放在这吧, 有心思的可以优化一下
func (svc *UserService) initEmailRe() {
	svc.emailRe = regexp.MustCompile(`^[A-Za-z0-9\x{4e00}-\x{9fa5}]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
}
