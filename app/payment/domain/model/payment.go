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

package model

import (
	"github.com/shopspring/decimal"

	"github.com/west2-online/DomTok/pkg/constants"
)

// PaymentOrder 支付订单表（删除了三个时间）
type PaymentOrder struct {
	ID                        int64           `gorm:"primaryKey;autoIncrement;comment:支付订单的唯一标识"`
	OrderID                   int64           `gorm:"not null;comment:商户订单号"`
	UserID                    int64           `gorm:"not null;comment:用户的唯一标识"`
	Amount                    decimal.Decimal `gorm:"type:decimal(14,4);not null;comment:订单总金额"`
	Status                    int64           `gorm:"not null;default:0;comment:支付状态：0-待支付，1-处理中，2-成功支付，3-支付失败"`
	MaskedCreditCardNumber    string          `gorm:"size:19;comment:信用卡号（仅存储掩码，如 **** **** **** 1234）"`
	CreditCardExpirationYear  int64           `gorm:"comment:信用卡到期年"`
	CreditCardExpirationMonth int64           `gorm:"comment:信用卡到期月"`
	Description               string          `gorm:"size:255;comment:订单描述信息"`
}

// PaymentRefund 退款表
type PaymentRefund struct {
	ID                        int64           `gorm:"primaryKey;autoIncrement;comment:支付退款的唯一标识"`
	OrderID                   int64           `gorm:"size:64;not null;comment:关联的商户订单号"`
	UserID                    int64           `gorm:"not null;comment:用户的唯一标识"`
	RefundAmount              decimal.Decimal `gorm:"type:decimal(15,4);not null;comment:退款金额，单位为元"`
	RefundReason              string          `gorm:"size:255;comment:退款原因"`
	Status                    int64           `gorm:"not null;default:0;comment:退款状态：0-待处理，1-处理中，2-成功退款，3-退款失败"`
	MaskedCreditCardNumber    string          `gorm:"size:19;comment:信用卡号（仅存储掩码，如 **** **** **** 1234）"`
	CreditCardExpirationYear  int64           `gorm:"comment:信用卡到期年"`
	CreditCardExpirationMonth int64           `gorm:"comment:信用卡到期月"`
}

// PaymentLedger 流水信息表
type PaymentLedger struct {
	ID              int64           `gorm:"primaryKey;autoIncrement;comment:流水ID"`
	ReferenceID     int64           `gorm:"not null;comment:关联的支付订单或退款订单ID"`
	UserID          int64           `gorm:"not null;comment:用户ID"`
	Amount          decimal.Decimal `gorm:"type:decimal(15,4);not null;comment:交易金额（正数表示收入，负数表示支出）"`
	TransactionType int64           `gorm:"not null;comment:交易类型：1-支付，2-退款，3-手续费，4-调整"`
	Status          int64           `gorm:"not null;default:0;comment:交易状态：0-待处理，1-成功，2-失败"`
}

func (p *PaymentOrder) ToLedger() *PaymentLedger {
	return &PaymentLedger{
		ReferenceID:     p.ID,
		UserID:          p.UserID,
		Amount:          p.Amount,
		TransactionType: constants.LedgerTransactionTypePaymentCode,
	}
}

func (p *PaymentRefund) ToLedger() *PaymentLedger {
	return &PaymentLedger{
		ReferenceID:     p.ID,
		UserID:          p.UserID,
		Amount:          p.RefundAmount,
		TransactionType: constants.LedgerTransactionTypeRefundCode,
	}
}
