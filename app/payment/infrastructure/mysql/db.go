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
	"github.com/west2-online/DomTok/pkg/constants"
	paymentStatus "github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
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
			return paymentStatus.PaymentNotExist, errno.Errorf(errno.ServicePaymentOrderNotExist, "payment order not found")
		}
		// 这里报错了就不是业务错误了, 而是服务级别的错误
		return paymentStatus.PaymentNotExist, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query payment order: %v", err)
	}
	return paymentStatus.PaymentExist, nil // 查询成功，返回支付状态
}

// GetPaymentInfo 通过orderID查询payment的信息
func (db *paymentDB) GetPaymentInfo(ctx context.Context, orderID int64) (*model.PaymentOrder, error) {
	var paymentOrder PaymentOrder
	err := db.client.WithContext(ctx).Table(constants.PaymentTableName).Where("order_id = ?", orderID).First(&paymentOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.Errorf(errno.ServicePaymentOrderNotExist, "payment order not found")
		}
		// 这里报错了就不是业务错误了, 而是服务级别的错误
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query payment order_id: %v", err)
	}
	paymentOrderInfo := &model.PaymentOrder{
		ID:                        paymentOrder.ID,
		OrderID:                   paymentOrder.OrderID,
		UserID:                    paymentOrder.UserID,
		Amount:                    paymentOrder.Amount,
		Status:                    paymentOrder.Status,
		MaskedCreditCardNumber:    paymentOrder.MaskedCreditCardNumber,
		CreditCardExpirationYear:  paymentOrder.CreditCardExpirationYear,
		CreditCardExpirationMonth: paymentOrder.CreditCardExpirationMonth,
		Description:               paymentOrder.Description,
	}
	return paymentOrderInfo, nil // 查询成功，返回支付状态
}

// ConvertToDBModel 转换函数
func ConvertToDBModel(p *model.PaymentOrder) (*PaymentOrder, error) {
	if p == nil {
		// TODO 这里应该是用errno.ParamVerifyErrorCode这个错误码？
		return nil, errno.Errorf(errno.ParamVerifyErrorCode, "ConvertToDBModel: input payment order is nil")
	}
	return &PaymentOrder{
		ID:                        p.ID,
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

func (db *paymentDB) UpdatePaymentStatus(ctx context.Context, orderID int64, status int) error {
	// 更新支付状态
	result := db.client.WithContext(ctx).Table(constants.PaymentTableName).
		Where("order_id = ?", orderID).
		Update("status", status)
	if result.Error != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update payment status: %v", result.Error)
	}
	if result.RowsAffected <= 0 {
		return errno.Errorf(errno.PaymentOrderNotExist, "mysql: payment order not found but tried to update")
	}
	return nil
}

// UpdatePaymentStatusToSuccessAndCreateLedgerAsTransaction 更新支付状态为成功并创建流水表项(作为事务)
func (db *paymentDB) UpdatePaymentStatusToSuccessAndCreateLedgerAsTransaction(ctx context.Context,
	order *model.PaymentOrder,
) error {
	err := db.client.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)
		// 更新支付状态
		result := tx.Table(constants.PaymentTableName).
			Where("order_id = ?", order.OrderID).
			Update("status", paymentStatus.PaymentStatusSuccessCode)
		if result.Error != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update payment status: %v", result.Error)
		}
		if result.RowsAffected <= 0 {
			return errno.Errorf(errno.PaymentOrderNotExist, "mysql: payment order not found but tried to update")
		}

		ledger := order.ToLedger()
		ledger.Status = paymentStatus.LedgerStatusSuccessCode
		dbLedger, err := ConvertLedgerToDBModel(ledger)
		if err != nil {
			return errno.Errorf(errno.InternalServiceErrorCode, "failed to convert ledger: %v", err)
		}
		// 创建流水表项
		if err = tx.Table(constants.PaymentLedgerTableName).WithContext(ctx).Create(dbLedger).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create ledger: %v", err)
		}
		return nil
	})
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update payment status and ledger status: %v", err)
	}
	return nil
}

// ConvertRefundToDBModel 转换函数
func ConvertRefundToDBModel(p *model.PaymentRefund) (*PaymentRefund, error) {
	if p == nil {
		// TODO 这里应该是用errno.ParamVerifyErrorCode这个错误码？
		return nil, errno.Errorf(errno.ParamVerifyErrorCode, "ConvertToDBModel: input payment order is nil")
	}
	return &PaymentRefund{
		ID:                        p.ID,
		OrderID:                   p.OrderID,
		UserID:                    p.UserID,
		RefundAmount:              p.RefundAmount,
		RefundReason:              p.RefundReason,
		Status:                    p.Status,
		MaskedCreditCardNumber:    p.MaskedCreditCardNumber,
		CreditCardExpirationYear:  p.CreditCardExpirationYear,
		CreditCardExpirationMonth: p.CreditCardExpirationMonth,
	}, nil
}

func (db *paymentDB) CreateRefund(ctx context.Context, p *model.PaymentRefund) error {
	// 将 entity 转换成 MySQL 需要的 refundOrder 结构
	paymentOrder, err := db.GetPaymentInfo(ctx, p.OrderID)
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "CreateRefund: failed to get payment info: %v", err)
	}
	// 通过去表格里查询获取这四个数据，发起退款的时候前端不需要传这些东西，如果查不到就说明有错误直接报错就好
	p.MaskedCreditCardNumber = paymentOrder.MaskedCreditCardNumber
	p.CreditCardExpirationYear = paymentOrder.CreditCardExpirationYear
	p.CreditCardExpirationYear = paymentOrder.CreditCardExpirationYear
	p.UserID = paymentOrder.UserID
	refundOrder, err := ConvertRefundToDBModel(p)
	if err != nil {
		logger.Errorf("CreateRefund: failed to convert refund order: %v", err)
		return errno.Errorf(errno.InternalServiceErrorCode, "CreateRefund: failed to convert refund order: %v", err)
	}

	// 插入数据库
	if err = db.client.WithContext(ctx).Create(refundOrder).Error; err != nil {
		logger.Errorf("CreateRefund: failed to create refund order: %v", err)
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create refund: %v", err)
	}
	logger.Infof("CreateRefund: refund order created successfully")
	return nil
}

func (db *paymentDB) GetRefundInfoByOrderID(ctx context.Context, orderID int64) (*model.PaymentRefund, error) {
	var refundOrder PaymentRefund
	err := db.client.WithContext(ctx).Table(constants.PaymentRefundTableName).Where("order_id = ?", orderID).First(&refundOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NewErrNo(errno.PaymentRefundNotExist, "refund order not found")
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query refund order_id: %v", err)
	}
	return &model.PaymentRefund{
		ID:                        refundOrder.ID,
		OrderID:                   refundOrder.OrderID,
		UserID:                    refundOrder.UserID,
		RefundAmount:              refundOrder.RefundAmount,
		RefundReason:              refundOrder.RefundReason,
		Status:                    refundOrder.Status,
		MaskedCreditCardNumber:    refundOrder.MaskedCreditCardNumber,
		CreditCardExpirationYear:  refundOrder.CreditCardExpirationYear,
		CreditCardExpirationMonth: refundOrder.CreditCardExpirationMonth,
	}, nil
}

// UpdateRefundStatusByOrderIDAndStatus 更新退款状态
func (db *paymentDB) UpdateRefundStatusByOrderIDAndStatus(ctx context.Context, orderID int64, prevStatus, nextStatus int) error {
	result := db.client.WithContext(ctx).Table(constants.PaymentRefundTableName).
		Where("order_id = ?", orderID).Where("status = ?", prevStatus).
		Update("status", nextStatus)
	if result.Error != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update refund status: %v", result.Error)
	}
	if result.RowsAffected <= 0 {
		return errno.Errorf(errno.PaymentRefundNotExist, "mysql: refund order not found but tried to update")
	}
	return nil
}

// UpdateRefundStatusToSuccessAndCreateLedgerAsTransaction 更新退款状态为成功并创建流水表项(作为事务)
func (db *paymentDB) UpdateRefundStatusToSuccessAndCreateLedgerAsTransaction(ctx context.Context,
	refund *model.PaymentRefund,
) error {
	err := db.client.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)
		// 更新退款状态
		result := tx.Table(constants.PaymentRefundTableName).
			Where("order_id = ?", refund.OrderID).
			Update("status", paymentStatus.RefundStatusSuccessCode)
		if result.Error != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update refund status: %v", result.Error)
		}
		if result.RowsAffected <= 0 {
			return errno.Errorf(errno.PaymentRefundNotExist, "mysql: refund order not found but tried to update")
		}

		ledger := refund.ToLedger()
		ledger.Status = paymentStatus.LedgerStatusSuccessCode
		dbLedger, err := ConvertLedgerToDBModel(ledger)
		if err != nil {
			return errno.Errorf(errno.InternalServiceErrorCode, "failed to convert ledger: %v", err)
		}
		// 创建流水表项
		if err = tx.Table(constants.PaymentLedgerTableName).WithContext(ctx).Create(dbLedger).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create ledger: %v", err)
		}
		return nil
	})
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update refund status and ledger status: %v", err)
	}
	return nil
}

// ConvertLedgerToDBModel 转换函数
func ConvertLedgerToDBModel(p *model.PaymentLedger) (*PaymentLedger, error) {
	if p == nil {
		return nil, errno.Errorf(errno.ParamVerifyErrorCode, "ConvertToDBModel: input payment order is nil")
	}
	return &PaymentLedger{
		ID:              p.ID,
		ReferenceID:     p.ReferenceID,
		UserID:          p.UserID,
		Amount:          p.Amount,
		TransactionType: p.TransactionType,
		Status:          p.Status,
	}, nil
}
