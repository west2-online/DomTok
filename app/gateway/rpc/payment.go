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
