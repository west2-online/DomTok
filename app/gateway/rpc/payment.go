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

func RequestPaymentTokenRPC(ctx context.Context, req *payment.PaymentTokenRequest) (token string, err error) {
	resp, err := paymentClient.RequestPaymentToken(ctx, req)
	// 这里的 err 是属于 RPC 间调用的错误，例如 network error
	// 而业务错误则是封装在 resp.base 当中的
	if err != nil {
		logger.Errorf("RequestPaymentTokenRPC: RPC called failed: %v", err.Error())
		return "", errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		// TODO
		return "", errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.PaymentToken, nil
}
