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

	"github.com/west2-online/DomTok/app/payment/domain/repository"
	orderrpc "github.com/west2-online/DomTok/kitex_gen/order"
	"github.com/west2-online/DomTok/kitex_gen/order/orderservice"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/utils"
)

type paymentRPC struct {
	order orderservice.Client
}

func NewPaymentRPC(order orderservice.Client) repository.PaymentRPC {
	return &paymentRPC{order: order}
}

func (rpc *paymentRPC) PaymentIsOrderExist(ctx context.Context, orderID int64) (orderExistInfo bool, err error) {
	orderRpcReq := &orderrpc.IsOrderExistReq{
		OrderID: orderID,
	}
	resp, err := rpc.order.IsOrderExist(ctx, orderRpcReq)
	/*if err != nil {
		return false, fmt.Errorf("rpc.order.IsOrderExist: %w", err)
	}
	if !utils.IsSuccess(resp.Base) {
		return false, fmt.Errorf("rpc.order.IsOrderExist: %v", resp.Base.Msg)
	}*/
	if err = utils.ProcessRpcError("payment.IsOrderExist", resp, err); err != nil {
		return false, err
	}
	orderExistInfo = resp.Exist
	return orderExistInfo, nil
}

func (rpc *paymentRPC) GetOrderStatus(ctx context.Context, orderID int64) (exist bool, expire int64, err error) {
	resp, err := rpc.order.IsOrderExist(ctx, &orderrpc.IsOrderExistReq{OrderID: orderID})
	if err = utils.ProcessRpcError("rpc.order.IsOrderExist", resp, err); err != nil {
		return false, 0, err
	}
	return resp.Exist, resp.OrderExpire, nil
}

func (rpc *paymentRPC) OrderPaymentCancel(ctx context.Context, orderID int64, paymentAt int64, paymentStyle string) error {
	req := &orderrpc.UpdateOrderStatusReq{
		OrderID:       orderID,
		PaymentStatus: constants.PaymentStatusFailedCode,
		PaymentAt:     paymentAt,
		PaymentStyle:  paymentStyle,
	}
	resp, err := rpc.order.OrderPaymentCancel(ctx, req)
	if err = utils.ProcessRpcError("rpc.order.OrderPaymentCancel", resp, err); err != nil {
		return err
	}
	return nil
}

func (rpc *paymentRPC) OrderPaymentSuccess(ctx context.Context, orderID int64, paymentAt int64, paymentStyle string) error {
	req := &orderrpc.UpdateOrderStatusReq{
		OrderID:       orderID,
		PaymentStatus: constants.PaymentStatusSuccessCode,
		PaymentAt:     paymentAt,
		PaymentStyle:  paymentStyle,
	}
	resp, err := rpc.order.OrderPaymentSuccess(ctx, req)
	if err = utils.ProcessRpcError("rpc.order.OrderPaymentSuccess", resp, err); err != nil {
		return err
	}
	return nil
}
