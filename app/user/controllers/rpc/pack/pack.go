package pack

import (
	"github.com/west2-online/DomTok/app/user/entities"
	"github.com/west2-online/DomTok/kitex_gen/model"
)

// BuildUser 将 entities 定义的 User 实体转换成 idl 定义的 RPC 交流实体，类似 dto
func BuildUser(u *entities.User) *model.UserInfo {
	return &model.UserInfo{
		Uid:   u.Uid,
		Name:  u.UserName,
		Email: u.Email,
		Phone: u.Phone,
	}
}
