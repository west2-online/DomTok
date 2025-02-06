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
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

// 通用的RPC客户端初始化函数
func initRPCClient[T any](serviceName string, newClientFunc func(string, ...client.Option) (T, error)) (*T, error) { //nolint
	if config.Etcd == nil || config.Etcd.Addr == "" {
		return nil, errors.New("config.Etcd.Addr is nil")
	}
	// 初始化Etcd Resolver
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalETCDErrorCode, fmt.Sprintf("initRPCClient etcd.NewEtcdResolver failed: %v", err))
	}
	// 初始化具体的RPC客户端
	client, err := newClientFunc(serviceName, client.WithResolver(r), client.WithMuxConnection(constants.MuxConnection))
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalInitRpcErrorCode, fmt.Sprintf("initRPCClient NewClient failed: %v", err))
	}
	return &client, nil
}
