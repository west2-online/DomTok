package pack

import (
	"github.com/west2-online/DomTok/app/gateway/model/model"
	rpcModel "github.com/west2-online/DomTok/kitex_gen/model"
)

// BuildTokenInfo 将 RPC 交流实体转换成 http 返回的实体
func BuildTokenInfo(p *rpcModel.PaymentTokenInfo) *model.{
	return &model.UserInfo{
		UserId: u.UserId,
		Name:   u.Name,
	}
}
