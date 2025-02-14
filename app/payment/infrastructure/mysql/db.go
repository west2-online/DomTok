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

import (
	"context"
	"github.com/west2-online/DomTok/app/payment/domain/model"
	"github.com/west2-online/DomTok/app/payment/domain/repository"
	"github.com/west2-online/DomTok/pkg/errno"
	"gorm.io/gorm"
)

type PaymentDB struct {
	client *gorm.DB
}

// 为什么这里会报错 return &paymentDB{client: client}
func NewPaymentDB(client *gorm.DB) repository.PaymentDB {
	return &paymentDB{client: client}
}

// ConvertPayment 后面把转换函数单独抽出来
func (db *paymentDB) ConvertPayment(ctx context.Context, p *model.Payment) (*model.Payment, error) {
	return nil, nil
}

func (db *paymentDB) CreatePayment(ctx context.Context, p *model.Payment) error {
	// 将 entity 转换成 mysql 这边的 model
	// TODO 可以考虑整一个函数统一转化, 放在这里占了太多行, 而且这不是这个方法该做的. 这个方法应该做的是创建用户
	model := Payment{
		UserName: p.UserName,
		Password: p.Password,
		Email:    p.Email,
	}

	if err := db.client.WithContext(ctx).Create(model).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create user: %v", err)
	}
	return nil
}
