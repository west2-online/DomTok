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

package constants

import "time"

const (
	OrderNotExist            = false
	PaymentExist             = true
	PaymentNotExist          = false
	UserNotExist             = false
	PaymentSecretKey         = "west2online"
	RefundSecretKey          = "west2online"
	ExpirationDuration       = 15 * time.Minute
	RefundExpirationDuration = 15 * time.Minute
	// TODO 这一个常量要改
	PingTime = 2
)

const (
	RedisStoreSuccess       = true  // 成功
	RedisStoreFailed        = false // Redis 存储失败
	RedisValid              = true
	RedisMinute             = 60
	RedisHour               = 3600
	RedisDay                = 86400
	RedisDayPlaceholder     = "1"
	RedisCheckTimesInMinute = 3
)

const (
	PaymentStatusPendingCode    = iota // 待支付
	PaymentStatusProcessingCode        // 处理中
	PaymentStatusSuccessCode           // 成功支付
	PaymentStatusFailedCode            // 支付失败
)

const (
	PaymentStatusPending    = "待支付"  // 待支付
	PaymentStatusProcessing = "处理中"  // 处理中
	PaymentStatusSuccess    = "成功支付" // 成功支付
	PaymentStatusFailed     = "支付失败" // 支付失败
)

const (
	RefundStatusPendingCode = iota
	RefundStatusProcessingCode
	RefundStatusSuccessCode
	RefundStatusFailedCode
)

const (
	RefundStatusPending    = "待退款"
	RefundStatusProcessing = "退款中"
	RefundStatusSuccess    = "成功退款"
	RefundStatusFailed     = "退款失败"
)

const (
	LedgerStatusPendingCode = iota
	LedgerStatusSuccessCode
	LedgerStatusFailedCode
)

const (
	LedgerStatusPending = "待处理"
	LedgerStatusSuccess = "成功"
	LedgerStatusFailed  = "失败"
)

const (
	LedgerTransactionTypePaymentCode = iota + 1
	LedgerTransactionTypeRefundCode
	LedgerTransactionTypeFeeCode
	LedgerTransactionTypeAdjustmentCode
)

const (
	LedgerTransactionTypePayment    = "支付"
	LedgerTransactionTypeRefund     = "退款"
	LedgerTransactionTypeFee        = "手续费"
	LedgerTransactionTypeAdjustment = "调整"
)

func GetPaymentStatus(code int8) string {
	switch code {
	case PaymentStatusPendingCode:
		return PaymentStatusPending
	case PaymentStatusProcessingCode:
		return PaymentStatusProcessing
	case PaymentStatusSuccessCode:
		return PaymentStatusSuccess
	case PaymentStatusFailedCode:
		return PaymentStatusFailed
	default:
		return PaymentStatusFailed
	}
}

func GetRefundStatus(code int8) string {
	switch code {
	case RefundStatusPendingCode:
		return RefundStatusPending
	case RefundStatusProcessingCode:
		return RefundStatusProcessing
	case RefundStatusSuccessCode:
		return RefundStatusSuccess
	case RefundStatusFailedCode:
		return RefundStatusFailed
	default:
		return RefundStatusFailed
	}
}

func GetLedgerStatus(code int8) string {
	switch code {
	case LedgerStatusPendingCode:
		return LedgerStatusPending
	case LedgerStatusSuccessCode:
		return LedgerStatusSuccess
	case LedgerStatusFailedCode:
		return LedgerStatusFailed
	default:
		return LedgerStatusFailed
	}
}

func GetLedgerTransactionType(code int8) string {
	switch code {
	case LedgerTransactionTypePaymentCode:
		return LedgerTransactionTypePayment
	case LedgerTransactionTypeRefundCode:
		return LedgerTransactionTypeRefund
	case LedgerTransactionTypeFeeCode:
		return LedgerTransactionTypeFee
	case LedgerTransactionTypeAdjustmentCode:
		return LedgerTransactionTypeAdjustment
	default:
		return LedgerTransactionTypeAdjustment
	}
}
