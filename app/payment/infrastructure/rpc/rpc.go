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
	"fmt"

	"github.com/west2-online/DomTok/app/payment/domain/repository"
	orderrpc "github.com/west2-online/DomTok/kitex_gen/order"
	"github.com/west2-online/DomTok/kitex_gen/order/orderservice"
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
	if err != nil {
		return false, fmt.Errorf("rpc.order.IsOrderExist: %w", err)
	}
	orderExistInfo = resp.Exist
	return orderExistInfo, nil
}
