package rpc

import (
	"context"

	"github.com/west2-online/DomTok/kitex_gen/user"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

func InitUserRPC() {
	c, err := client.InitUserRPC()
	if err != nil {
		logger.Fatalf("api.rpc.user InitUserRPC failed, err is %v", err)
	}
	userClient = *c
}

func RegisterRPC(ctx context.Context, req *user.RegisterRequest) (uid int64, err error) {
	resp, err := userClient.Register(ctx, req)
	// 这里的 err 是属于 RPC 间调用的错误，例如 network error
	// 而业务错误则是封装在 resp.base 当中的
	if err != nil {
		logger.Errorf("RegisterRPC: RPC called failed: %v", err.Error())
		return 0, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return 0, errno.BizError.WithMessage(resp.Base.Msg)
	}
	return resp.UserID, nil

}
