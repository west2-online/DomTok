package rpc

import (
	"context"

	"github.com/west2-online/DomTok/app/user/controllers/rpc/pack"
	"github.com/west2-online/DomTok/app/user/entities"
	"github.com/west2-online/DomTok/kitex_gen/user"
	"github.com/west2-online/DomTok/pkg/base"
	"github.com/west2-online/DomTok/pkg/logger"
)

// UseCasePort 在这里定义出 use case 的接口，原本接口应该是放在 domain/entities 中的
// golang 的接口是[隐式]实现的，所以不会出现 use case 实现了接口而导致的[循环依赖问题]
// 命名后缀带上 port 表示接口，命名来源于[六边形架构]
type UseCasePort interface {
	RegisterUser(ctx context.Context, user *entities.User) (uid int64, err error)
	Login(ctx context.Context, user *entities.User) (*entities.User, error)
}

// UserHandler 实现 idl 中定义的 RPC 接口
type UserHandler struct {
	useCase UseCasePort
}

func NewUserHandler(useCase UseCasePort) *UserHandler {
	return &UserHandler{useCase: useCase}
}

func (h *UserHandler) Register(ctx context.Context, req *user.RegisterRequest) (r *user.RegisterResponse, err error) {
	resp := new(user.RegisterResponse)
	// 将 req 转换为 entities.User
	user := &entities.User{
		UserName: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
	}
	// 调用 use case
	uid, err := h.useCase.RegisterUser(ctx, user)
	if err != nil {
		logger.Infof("UserHandler.Register: RegisterUser failed, err: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.UserID = uid
	resp.Base = base.BuildSuccessResp()
	return resp, nil
}

func (h *UserHandler) Login(ctx context.Context, req *user.LoginRequest) (r *user.LoginResponse, err error) {
	resp := new(user.LoginResponse)

	user := &entities.User{
		UserName: req.Username,
		Password: req.Password,
	}

	ans, err := h.useCase.Login(ctx, user)
	if err != nil {
		logger.Infof("UserHandler.Login: Login failed, err: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.User = pack.BuildUser(ans)
	return resp, nil
}
