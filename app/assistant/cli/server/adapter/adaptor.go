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

package adapter

import (
	"context"

	"github.com/west2-online/DomTok/app/gateway/model/api/cart"
	"github.com/west2-online/DomTok/app/gateway/model/api/order"
)

// ServerCaller is the interface for calling the server
// It is used by the AI client to call the server
// List required methods here
type ServerCaller interface {
	// Ping An example method
	Ping(ctx context.Context) ([]byte, error)

	CartShow(ctx context.Context, params *cart.ShowCartGoodsListRequest) ([]byte, error)
	CartPurchase(ctx context.Context, params *cart.PurChaseCartGoodsRequest) ([]byte, error)
	OrderList(ctx context.Context, params *order.ViewOrderListReq) ([]byte, error)
	OrderView(ctx context.Context, params *order.ViewOrderReq) ([]byte, error)
}
