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

package usecase

import (
	"context"
	"fmt"

	"github.com/west2-online/DomTok/app/payment/domain/model"
	paymentStatus "github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

// CreatePayment 这里定义一些具体的方法和函数，比如校验密码，加密密码，创建用户之类的
func (uc *paymentUseCase) CreatePayment(ctx context.Context, orderID int64) (*model.PaymentOrder, error) {
	return nil, nil
}

func (uc *paymentUseCase) GetPaymentToken(ctx context.Context, orderID int64) (token string, expTime int64, err error) {
	// 1. 检查订单是否存在
	// TODO 记得删除注释
	/*var orderInfo bool
	orderInfo, err = uc.svc.CheckOrderExist(ctx, orderID)
	if err != nil {
		return "", 0, fmt.Errorf("check order existed failed:%w", err)
	}
	if orderInfo == paymentStatus.OrderNotExist {
		return "", 0, errno.NewErrNo(errno.ServicePaymentOrderNotExist, "order does not exist")
	}
	*/
	// 2. 获取用户id,无需检查用户是否存在
	// 获取用户id
	var uid int64
	uid, err = uc.svc.GetUserID(ctx)
	if err != nil {
		return "", 0, fmt.Errorf("get user id failed:%w", err)
	}

	// 3. 检查订单支付信息
	var paymentInfo bool
	paymentInfo, err = uc.db.CheckPaymentExist(ctx, orderID)
	if err != nil {
		return "", 0, fmt.Errorf("check payment existed failed:%w", err)
	}
	if paymentInfo == paymentStatus.PaymentNotExist { // 如果订单不存在
		// 创建支付订单
		// TODO 待完善
		_, err := uc.svc.CreatePaymentInfo(ctx, orderID)
		if err != nil {
			return "", 0, fmt.Errorf("create payment info failed:%w", err)
		}
	} else if paymentInfo == paymentStatus.PaymentExist { // 如果订单存在
		// 获取订单的支付状态
		payStatus, err := uc.db.GetPaymentInfo(ctx, orderID)
		if err != nil {
			return "", 0, fmt.Errorf("get payment info failed:%w", err)
		}
		// 如果订单正在支付或者已经支付完成，则拒绝进行接下来的生成令牌的活动
		if payStatus.Status == paymentStatus.PaymentStatusSuccessCode || payStatus.Status == paymentStatus.PaymentStatusProcessingCode {
			return "", 0, errno.Errorf(errno.ServicePaymentIsProcessing, "payment is processing or has already done")
		}
	}

	// 4. HMAC生成支付令牌
	token, expTime, err = uc.svc.GeneratePaymentToken(ctx, orderID)
	if err != nil {
		logger.Errorf("Error generating payment token: orderID:%d,err:%v", orderID, err)
		return "", 0, fmt.Errorf("generate payment token failed:%w", err)
	}
	var redisStatus bool
	// 5. 存储令牌到 Redis
	redisStatus, err = uc.svc.StorePaymentToken(ctx, token, expTime, uid, orderID)
	if err != nil && redisStatus != paymentStatus.RedisStoreSuccess {
		logger.Errorf("Error store payment token: orderID:%d,userID:%d,err:%v", orderID, uid, err)
		return "", 0, fmt.Errorf("store payment token failed:%w", err)
	}
	logger.Infof("Success generating payment token: orderID:%d,token:%s", orderID, token)
	return token, expTime, nil
}

// CreateRefund 发起退款请求
func (uc *paymentUseCase) CreateRefund(ctx context.Context, orderID int64) (refundStatus int64, refundID int64, err error) {
	// 1. 检查订单是否存在
	// TODO记得删除注释
	/*orderExists, err := uc.svc.CheckOrderExist(ctx, orderID)
	if err != nil {
		return 0, 0, fmt.Errorf("check order existence failed: %w", err)
	}
	if !orderExists {
		return 0, 0, errno.NewErrNo(errno.ServicePaymentOrderNotExist, "order does not exist")
	}*/
	// 2. 获取用户ID
	uid, err := uc.svc.GetUserID(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("get user id failed: %w", err)
	}
	// 3. Redis 限流检查
	var frequencyInfo bool
	var timeInfo bool
	frequencyInfo, timeInfo, err = uc.svc.CheckRedisRateLimiting(ctx, uid, orderID)
	if err != nil {
		return 0, 0, fmt.Errorf("check redis rate limiting failed: %w", err)
	}
	if frequencyInfo != paymentStatus.RedisValid {
		return 0, 0, fmt.Errorf("too many refund requests in a short time")
	}
	if timeInfo != paymentStatus.RedisValid {
		return 0, 0, fmt.Errorf("refund already requested for this order in the last 24 hours")
	}

	// 4. 创建退款信息
	refundID, err = uc.svc.CreateRefundInfo(ctx, orderID)
	if err != nil {
		return 0, 0, fmt.Errorf("create refund info failed: %w", err)
	}
	refundStatus = paymentStatus.RefundStatusProcessingCode
	return refundStatus, refundID, nil
}

// RefundReview 退款审核
func (uc *paymentUseCase) RefundReview(ctx context.Context, orderID int64, passed bool) error {
	// 1. 检查订单是否存在
	orderExist, orderExpired, err := uc.svc.GetOrderStatus(ctx, orderID)
	if err != nil {
		return err
	}
	// 订单不存在或者订单已经过期
	if !orderExist {
		return errno.Errorf(errno.ServiceOrderNotFound, "order does not exist")
	}
	if orderExpired {
		return errno.Errorf(errno.ServiceOrderExpired, "order has expired")
	}

	// 2. 用户是否存在
	uid, err := uc.svc.GetUserID(ctx)
	if err != nil {
		return err
	}

	// 3. 用户是否有权限发起退款
	hasPermission, err := uc.svc.CheckAdminPermission(ctx, uid)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errno.AuthNoOperatePermission
	}

	// 4. 检查退款信息是否存在
	refund, err := uc.db.GetRefundInfoByOrderID(ctx, orderID)
	if err != nil {
		return err
	}
	if refund == nil {
		return errno.Errorf(errno.ServicePaymentRefundNotExist, "refund does not exist")
	}

	// 5. 更新退款状态为处理中
	err = uc.db.UpdateRefundStatusByOrderIDAndStatus(ctx, orderID,
		paymentStatus.RefundStatusPendingCode, paymentStatus.RefundStatusProcessingCode)
	if err != nil {
		return err
	}

	if passed {
		refundAt, style, err := uc.svc.Refund(ctx)
		if err != nil {
			return err
		}
		err = uc.svc.CancelOrder(ctx, orderID, refundAt, style)
		if err != nil {
			return err
		}
		err = uc.db.UpdateRefundStatusToSuccessAndCreateLedgerAsTransaction(ctx, refund)
		if err != nil {
			return err
		}
	} else {
		err = uc.db.UpdateRefundStatusByOrderIDAndStatus(ctx, orderID,
			paymentStatus.RefundStatusProcessingCode, paymentStatus.RefundStatusFailedCode)
		if err != nil {
			return err
		}
	}

	return nil
}

// PaymentCheckout 支付结算
func (uc *paymentUseCase) PaymentCheckout(ctx context.Context, orderID int64, token string) error {
	// 1. 检查订单是否存在
	orderExist, orderExpired, err := uc.svc.GetOrderStatus(ctx, orderID)
	if err != nil {
		return err
	}
	// 订单不存在或者订单已经过期
	if !orderExist {
		return errno.Errorf(errno.ServiceOrderNotFound, "order does not exist")
	}
	if orderExpired {
		return errno.Errorf(errno.ServiceOrderExpired, "order has expired")
	}

	// 2. 用户是否存在
	uid, err := uc.svc.GetUserID(ctx)
	if err != nil {
		return err
	}

	// 3. 支付令牌是否在 Redis 中
	exist, exp, err := uc.svc.GetExpiredAtAndDelPaymentToken(ctx, token, uid, orderID)
	if !exist && err == nil {
		return errno.Errorf(errno.IllegalOperatorCode, "duplicate payment request, mismatched order or token has expired")
	}

	var order *model.PaymentOrder
	// 4.1 Redis 错误，进入第二层校验
	if err != nil {
		// 4.1.1 数据库中查询 status 状态
		order, err = uc.db.GetPaymentInfo(ctx, orderID)
		if err != nil {
			return err
		}

		// 4.1.2 校验状态是否为待支付
		if order.Status != paymentStatus.PaymentStatusPendingCode {
			// 表示重复操作，返回重复支付错误
			return errno.Errorf(errno.IllegalOperatorCode, "duplicate payment request, mismatched order or token has expired")
		}

		// 4.1.3 更新支付状态为处理中
		err = uc.db.UpdatePaymentStatus(ctx, orderID, paymentStatus.PaymentStatusProcessingCode)
		if err != nil {
			return err
		}
	}

	var rollbackBeforePay func() error
	if order == nil {
		// Redis 未出错，则需要回滚 Redis
		rollbackBeforePay = func() error {
			errRollback := uc.svc.PutBackPaymentToken(ctx, token, uid, orderID, exp)
			if errRollback != nil {
				return errRollback
			}
			return nil
		}
	} else {
		// Redis出错，则需要回滚数据库
		rollbackBeforePay = func() error {
			errRollback := uc.db.UpdatePaymentStatus(ctx, orderID, paymentStatus.PaymentStatusPendingCode)
			if errRollback != nil {
				return errRollback
			}
			return nil
		}
	}

	// 临时挂起支付信息，用于预确认订单
	payAt, style, err := uc.svc.GetPayInfo(ctx)
	if err != nil {
		return err
	}

	err = uc.svc.ConfirmOrder(ctx, orderID, payAt, style)
	// 确认订单失败，回滚 Redis 或数据库
	if err != nil {
		rollbackErr := rollbackBeforePay()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	transactionSuccess := true
	// 5. 调用支付接口
	// 此时的支付信息为确认订单的应填支付信息
	// 真实场景下应当更新订单的相关字段，由于支付过程为模拟，故此处不做处理
	_, _, err = uc.svc.Pay(ctx)
	if err != nil {
		transactionSuccess = false
	}

	// 如果支付成功且订单为空，则从数据库中获取订单信息-用于后续创建流水表项
	if transactionSuccess && order == nil {
		order, err = uc.db.GetPaymentInfo(ctx, orderID)
		if err != nil {
			transactionSuccess = false
		}
	}

	// 更新支付状态为成功并创建流水表项
	if transactionSuccess {
		err = uc.db.UpdatePaymentStatusToSuccessAndCreateLedgerAsTransaction(ctx, order)
		if err != nil {
			transactionSuccess = false
		}
	}

	if !transactionSuccess {
		// 支付失败，更新支付状态为失败
		errX := uc.db.UpdatePaymentStatus(ctx, orderID, paymentStatus.PaymentStatusFailedCode)
		if errX != nil {
			return errX
		}
		return err
	}

	// 6. 将支付结果写入 Redis(未开放查看支付结果接口，故此处不做处理)
	return nil
}
