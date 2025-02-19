package service

import (
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

type OrderVerifyOps func() error

// Verify 验证多个条件
func (svc *OrderService) Verify(ops ...OrderVerifyOps) error {
	for _, opt := range ops {
		if err := opt(); err != nil {
			return err
		}
	}
	return nil
}

// VerifyOrderStatus 验证订单状态
func (svc *OrderService) VerifyOrderStatus(status int32) func() error {
	return func() error {
		if svc.GetOrderStatusMsg(int8(status)) == constants.OrderStatusUnknown {
			return errno.NewErrNo(errno.ServiceError, "invalid order status")
		}
		return nil
	}
}
