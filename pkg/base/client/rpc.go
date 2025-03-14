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

package client

import (
	"errors"
	"fmt"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/streamclient"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/kitex_gen/cart/cartservice"
	"github.com/west2-online/DomTok/kitex_gen/commodity/commodityservice"
	"github.com/west2-online/DomTok/kitex_gen/order/orderservice"
	"github.com/west2-online/DomTok/kitex_gen/payment/paymentservice"
	"github.com/west2-online/DomTok/kitex_gen/user/userservice"
	"github.com/west2-online/DomTok/pkg/constants"
)

// 通用的RPC客户端初始化函数
func initRPCClient[T any](serviceName string, newClientFunc func(string, ...client.Option) (T, error)) (*T, error) {
	if config.Etcd == nil || config.Etcd.Addr == "" {
		return nil, errors.New("config.Etcd.Addr is nil")
	}
	// 初始化Etcd Resolver
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		return nil, fmt.Errorf("initRPCClient etcd.NewEtcdResolver failed: %w", err)
	}
	// 初始化具体的RPC客户端
	client, err := newClientFunc(serviceName,
		client.WithResolver(r),
		client.WithMuxConnection(constants.MuxConnection),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: fmt.Sprintf(constants.KitexClientEndpointInfoFormat, serviceName)}),
	)
	if err != nil {
		return nil, fmt.Errorf("initRPCClient NewClient failed: %w", err)
	}
	return &client, nil
}

func InitUserRPC() (*userservice.Client, error) {
	return initRPCClient("user", userservice.NewClient)
}

func InitOrderRPC() (*orderservice.Client, error) {
	return initRPCClient(constants.OrderServiceName, orderservice.NewClient)
}

func InitCommodityRPC() (*commodityservice.Client, error) {
	return initRPCClient(constants.CommodityServiceName, commodityservice.NewClient)
}

func InitCommodityStreamClientRPC() (*commodityservice.StreamClient, error) {
	if config.Etcd == nil || config.Etcd.Addr == "" {
		return nil, errors.New("config.Etcd.Addr is nil")
	}
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		return nil, fmt.Errorf("initRPCClient etcd.NewEtcdResolver failed: %w", err)
	}
	cli := commodityservice.MustNewStreamClient(constants.CommodityServiceName, streamclient.WithResolver(r))
	return &cli, nil
}

func InitCartRPC() (*cartservice.Client, error) {
	return initRPCClient(constants.CartServiceName, cartservice.NewClient)
}

func InitPaymentRPC() (*paymentservice.Client, error) {
	return initRPCClient(constants.PaymentServiceName, paymentservice.NewClient)
}
