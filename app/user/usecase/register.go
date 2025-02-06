package usecase

import (
	"context"
	"fmt"

	"github.com/west2-online/DomTok/app/user/entities"
	"github.com/west2-online/DomTok/pkg/errno"
)

func (u *UseCase) RegisterUser(ctx context.Context, entity *entities.User) (uid int64, err error) {
	// 判断是否已经注册过
	exist, err := u.DB.IsUserExist(ctx, entity.UserName)
	if err != nil {
		return 0, fmt.Errorf("check user exist failed: %w", err)
	}
	if exist {
		// 返回错误码定义？

		// 原始错误
		return 0, errno.NewErrNo(errno.BizErrorCode, "user already exist")
	}
	// 校验 email
	if valid := entity.IsValidEmail(); !valid {
		return 0, errno.NewErrNo(errno.BizErrorCode, "invalid email")
	}
	// 加密密码，准备存入数据库
	if err = entity.EncryptPassword(); err != nil {
		return 0, errno.NewErrNo(errno.BizErrorCode, "encrypt password failed")
	}
	// 存入数据库
	if err = u.DB.CreateUser(ctx, entity); err != nil {
		return 0, fmt.Errorf("create user failed: %w", err)
	}
	return entity.Uid, nil
}
