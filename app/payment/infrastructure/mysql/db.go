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

// CheckPaymentExist 检查是否已经发起过支付，利用orderID在订单支付表里查询
func (db *paymentDB) CheckPaymentExist(ctx context.Context, orderID int64) (paymentInfo interface{}, err error) {
	var paymentOrder PaymentOrder
	// 利用orderID在订单支付表里查询是否已经发起过支付申请了（注意是订单支付表不是订单表）
	err = db.client.WithContext(ctx).Where("order_id = ?", orderID).First(&paymentOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return paymentStatus.PaymentNotExist, nil
		}
		// 这里报错了就不是业务错误了, 而是服务级别的错误
		return paymentStatus.PaymentNotExist, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query payment order: %v", err)
	}
	return paymentStatus.PaymentExist, nil // 查询成功，返回 user_id
}

// GetPaymentInfo 通过orderID查询payment的信息
func (db *paymentDB) GetPaymentInfo(ctx context.Context, orderID int64) (interface{}, error) {
	var paymentOrder PaymentOrder
	err := db.client.WithContext(ctx).Where("order_id = ?", orderID).First(&paymentOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return paymentStatus.PaymentNotExist, nil
		}
		// 这里报错了就不是业务错误了, 而是服务级别的错误
		return paymentStatus.PaymentNotExist, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query payment order_id: %v", err)
	}
	return paymentOrder.Status, nil // 查询成功，返回支付状态
}

// ConvertPayment TODO 后面把转换函数单独抽出来
func (db *paymentDB) ConvertPayment(ctx context.Context, p *model.PaymentOrder) (*model.PaymentOrder, error) {
	return nil, nil
}

func (db *paymentDB) CreatePayment(ctx context.Context, p *model.PaymentOrder) error {
	// 将 entity 转换成 mysql 这边的 paymentOrder
	// TODO 可以考虑整一个函数统一转化, 放在这里占了太多行, 而且这不是这个方法该做的. 这个方法应该做的是创建支付订单
	paymentOrder := PaymentOrder{
		OrderID: p.OrderID,
		UserID:  p.UserID,
	}
	if err := db.client.WithContext(ctx).Create(paymentOrder).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create payment: %v", err)
	}
	return nil
}
