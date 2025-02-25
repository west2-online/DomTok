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

package utils

import (
	"fmt"
	"github.com/west2-online/DomTok/kitex_gen/model"
	"github.com/west2-online/DomTok/pkg/errno"
)

// IsSuccess 通用的rpc结果处理
func IsSuccess(baseResp *model.BaseResp) bool {
	return baseResp.Code == errno.SuccessCode
}

type Baser interface {
	GetBase() *model.BaseResp
	IsSetBase() bool
}

func ProcessRpcError(calledMethod string, resp any, err error) error {
	// err 不为 nil 就直接返回即可, 这里用 errno 是因为在我们的 rpc 体系中, 能返回出来的 error 都不是我们自己定义的
	// 所有底层为 errno 的 err 都在中间件中被捕获了, 最后从 rpc 返回的都是 nil, 此处不为 nil 说明不是服务传出来的, 可能是框架或者网络错误
	if err != nil {
		return errno.NewErrNo(errno.InternalRPCErrorCode, fmt.Sprintf("failed to call %s,err: %v", calledMethod, err))
	}
	// 这里用了 any 来让 resp 传进来并且判断是否为 nil 是为了避免 nil 地狱
	if resp == nil {
		return errno.NewErrNo(errno.InternalRPCErrorCode, fmt.Sprintf("success call %s but resp is nil", calledMethod))
	}
	// 如果不能被断言为 Baser, 那这个 resp 有大问题, 甚至这里可以 panic. 因为所有我们自己的 resp 都是非 nil 且含有 Base 的
	baser, ok := resp.(Baser)
	if !ok {
		return errno.NewErrNo(errno.InternalServiceErrorCode, fmt.Sprintf("rpc`s resp that passed by %s don`t have model.Base", calledMethod))
	}

	if !baser.IsSetBase() {
		return errno.NewErrNo(errno.InternalRPCErrorCode, fmt.Sprintf("success call %s, but its resp.Base is nil", calledMethod))
	}
	base := baser.GetBase()
	// 这里也算是我们这个调用方的一个能得知的最根源的错误了, 所以当然使用 errno
	if !IsSuccess(base) {
		return errno.NewErrNo(base.Code, fmt.Sprintf(" call %s failed because: %v", calledMethod, base.Msg))
	}

	return nil
}
