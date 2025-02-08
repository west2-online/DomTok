package pack

import (
	"github.com/west2-online/DomTok/gateway/model/model"
	rpcModel "github.com/west2-online/DomTok/kitex_gen/model"
)

// BuildUserInfo 将 RPC 交流实体转换成 http 返回的实体
func BuildUserInfo(u *rpcModel.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		UID:   u.Uid,
		Name:  u.Name,
		Email: u.Email,
	}
}
