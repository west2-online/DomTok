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

package rpc

import (
	"context"
	"fmt"

	"github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/kitex_gen/payment"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

func InitPaymentRPC() {
	c, err := client.InitPaymentRPC()
	if err != nil {
		logger.Fatalf("api.rpc.payment InitPayemntRPC failed, err is %v", err)
	}
	paymentClient = *c
}

func RequestPaymentTokenRPC(ctx context.Context, req *payment.PaymentTokenRequest) (response *model.PaymentTokenInfo, err error) {
	logger.Info("RequestPaymentTokenRPC called") // 这里打印一下，看看是否调用到了
	fmt.Println("RequestPaymentTokenRPC: called with OrderID:", req.OrderID, "UserID:", req.UserID)

	resp, err := paymentClient.RequestPaymentToken(ctx, req)
	// 这里的 err 是属于 RPC 间调用的错误，例如 network error
	// 而业务错误则是封装在 resp.base 当中的
	if err != nil {
		logger.Errorf("RequestPaymentTokenRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		// TODO
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.TokenInfo, nil
}

func RequestRefundRPC(ctx context.Context, req *payment.RefundTokenRequest) (response int64, err error) {
	logger.Infof("RequestRefundTokenRPC called") // 记录日志，确保调用成功
	// 调用 RPC 获取退款令牌
	resp, err := paymentClient.RequestRefundInfo(ctx, req)
	if err != nil {
		logger.Errorf("RequestRefundRPC: RPC call failed: %v", err.Error())
		return 0, errno.InternalServiceError.WithError(err)
	}

	// 解析业务错误（即 RPC 返回的错误信息）
	if !utils.IsSuccess(resp.Base) {
		logger.Errorf("RequestRefundRPC: Business error: %s", resp.Base.Msg)
		return 0, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}

	// 返回退款令牌信息
	return resp.RefundID, nil
}
