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

	"github.com/west2-online/DomTok/app/payment/domain/model"
	"github.com/west2-online/DomTok/app/payment/domain/repository"
	"github.com/west2-online/DomTok/pkg/constants"
	paymentStatus "github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"gorm.io/gorm"
)

// paymentDB impl domain.PaymentDB defined domain
type paymentDB struct {
	client *gorm.DB
}

func NewPaymentDB(client *gorm.DB) repository.PaymentDB {
	return &paymentDB{client: client}
}

// CheckPaymentExist 检查是否已经发起过支付，利用orderID在订单支付表里查询
func (db *paymentDB) CheckPaymentExist(ctx context.Context, orderID int64) (paymentInfo bool, err error) {
	var paymentOrder PaymentOrder
	// 利用orderID在订单支付表里查询是否已经发起过支付申请了（注意是订单支付表不是订单表）
	err = db.client.WithContext(ctx).Table(constants.PaymentTableName).Where("order_id = ?", orderID).First(&paymentOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return paymentStatus.PaymentNotExist, nil
		}
		// 这里报错了就不是业务错误了, 而是服务级别的错误
		return paymentStatus.PaymentNotExist, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query payment order: %v", err)
	}
	return paymentStatus.PaymentExist, nil // 查询成功，返回支付状态
}

// GetPaymentInfo 通过orderID查询payment的信息
func (db *paymentDB) GetPaymentInfo(ctx context.Context, orderID int64) (interface{}, error) {
	var paymentOrder PaymentOrder
	err := db.client.WithContext(ctx).Table(constants.PaymentTableName).Where("order_id = ?", orderID).First(&paymentOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return paymentStatus.PaymentNotExist, nil
		}
		// 这里报错了就不是业务错误了, 而是服务级别的错误
		return paymentStatus.PaymentNotExist, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query payment order_id: %v", err)
	}
	return paymentOrder.Status, nil // 查询成功，返回支付状态
}

// ConvertToDBModel 转换函数
func ConvertToDBModel(p *model.PaymentOrder) (*PaymentOrder, error) {
	if p == nil {
		// TODO 这里应该是用errno.ParamVerifyErrorCode这个错误码？
		return nil, errno.Errorf(errno.ParamVerifyErrorCode, "ConvertToDBModel: input payment order is nil")
	}
	return &PaymentOrder{
		OrderID:                   p.OrderID,
		UserID:                    p.UserID,
		Amount:                    p.Amount,
		Status:                    p.Status,
		MaskedCreditCardNumber:    p.MaskedCreditCardNumber,
		CreditCardExpirationYear:  p.CreditCardExpirationYear,
		CreditCardExpirationMonth: p.CreditCardExpirationMonth,
		Description:               p.Description,
	}, nil
}

func (db *paymentDB) CreatePayment(ctx context.Context, p *model.PaymentOrder) error {
	// 将 entity 转换成 mysql 这边的 paymentOrder
	paymentOrder, err := ConvertToDBModel(p)
	if err != nil {
		return errno.Errorf(errno.InternalServiceErrorCode, "CreatePayment: failed to convert payment order: %v", err)
	}
	if err := db.client.WithContext(ctx).Create(paymentOrder).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create payment: %v", err)
	}
	return nil
}
