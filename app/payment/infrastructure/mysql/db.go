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
	"errors"

	"gorm.io/gorm"

	"github.com/west2-online/DomTok/app/payment/domain/model"
	"github.com/west2-online/DomTok/app/payment/domain/repository"
	paymentStatus "github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

// paymentDB impl domain.PaymentDB defined domain
type paymentDB struct {
	client *gorm.DB
}

func NewPaymentDB(client *gorm.DB) repository.PaymentDB {
	return &paymentDB{client: client}
}

// TODO 这里是直接传token去查询，还是要把token解析出来？
func (db *paymentDB) GetOrderByToken(ctx context.Context, paramToken string) (int64, error) {
	var paymentOrder PaymentOrder
	err := db.client.WithContext(ctx).Where("token = ?", paramToken).First(&paymentOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return paymentStatus.PaymentOrderNotExist, nil
		}
		// 这里报错了就不是业务错误了, 而是服务级别的错误
		return paymentStatus.PaymentOrderNotExist, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query payment token: %v", err)
	}
	return paymentOrder.OrderID, nil // 查询成功，返回 order_id
}

func (db *paymentDB) GetUserByToken(ctx context.Context, paramToken string) (int64, error) {
	var paymentOrder PaymentOrder
	err := db.client.WithContext(ctx).Where("token = ?", paramToken).First(&paymentOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return paymentStatus.UserNotExist, nil
		}
		// 这里报错了就不是业务错误了, 而是服务级别的错误
		return paymentStatus.UserNotExist, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query payment token: %v", err)
	}
	return paymentOrder.UserID, nil // 查询成功，返回 user_id
}

func (db *paymentDB) GetPaymentInfo(ctx context.Context, paramToken string) (int, error) {
	var paymentOrder PaymentOrder
	err := db.client.WithContext(ctx).Where("token = ?", paramToken).First(&paymentOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return paymentStatus.UserNotExist, nil
		}
		// 这里报错了就不是业务错误了, 而是服务级别的错误
		return paymentStatus.UserNotExist, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query payment token: %v", err)
	}
	return int(paymentOrder.Status), nil // 查询成功，返回支付状态
}

// ConvertPayment TODO 后面把转换函数单独抽出来
func (db *paymentDB) ConvertPayment(ctx context.Context, p *model.PaymentOrder) (*model.PaymentOrder, error) {
	return nil, nil
}

func (db *paymentDB) CreatePayment(ctx context.Context, p *model.PaymentOrder) error {
	// 将 entity 转换成 mysql 这边的 paymentOrder
	// TODO 可以考虑整一个函数统一转化, 放在这里占了太多行, 而且这不是这个方法该做的. 这个方法应该做的是创建用户
	paymentOrder := PaymentOrder{
		OrderID: p.OrderID,
		UserID:  p.UserID,
	}
	if err := db.client.WithContext(ctx).Create(paymentOrder).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create payment: %v", err)
	}
	return nil
}
