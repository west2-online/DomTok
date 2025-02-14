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
	"github.com/west2-online/DomTok/pkg/constants"
	"time"
)

type Payment struct {
	// model    gorm.Model
	ID                        int64      `gorm:"primaryKey;autoIncrement;comment:支付订单的唯一标识"`
	OrderID                   int64      `gorm:"notNull;comment:商户订单号"`
	UserID                    int64      `gorm:"notNull;comment:用户的唯一标识"`
	Amount                    float64    `gorm:"type:decimal(15,4);notNull;comment:订单总金额"`
	Status                    int8       `gorm:"type:tinyint;notNull;default:0;comment:支付状态：0-待支付，1-处理中，2-成功支付 3-支付失败"`
	MaskedCreditCardNumber    string     `gorm:"type:varchar(19);comment:信用卡号 国际信用卡号的最大长度为19 (仅存储掩码，如 **** **** **** 1234)"`
	CreditCardExpirationYear  int        `gorm:"type:int;comment:信用卡到期年"`
	CreditCardExpirationMonth int        `gorm:"type:int;comment:信用卡到期月"`
	Description               string     `gorm:"type:varchar(255);comment:订单描述信息"`
	CreatedAt                 time.Time  `gorm:"notNull;default:CURRENT_TIMESTAMP;comment:订单创建时间"`
	UpdatedAt                 time.Time  `gorm:"notNull;default:CURRENT_TIMESTAMP;comment:订单最后更新时间"`
	DeletedAt                 *time.Time `gorm:"type:timestamp;comment:订单删除时间"`
}

func (Payment) TableName() string {
	return constants.PaymentTableName
}
