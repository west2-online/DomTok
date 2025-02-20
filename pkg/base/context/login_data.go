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

package context

import (
	"context"
	"strconv"

	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

// WithLoginData 将LoginData加入到context中，通过metainfo传递到RPC server
func WithLoginData(ctx context.Context, uid int64) context.Context {
	return newContext(ctx, constants.LoginDataKey, strconv.FormatInt(uid, 10))
}

// GetLoginData 从context中取出LoginData
func GetLoginData(ctx context.Context) (int64, error) {
	user, ok := fromContext(ctx, constants.LoginDataKey)
	if !ok {
		return -1, errno.NewErrNo(errno.ParamMissingErrorCode, "Failed to get header in context")
	}

	value, err := strconv.ParseInt(user, 10, 64)
	if err != nil {
		return -1, errno.NewErrNo(errno.InternalServiceErrorCode, "Failed to get header in context when parse loginData")
	}
	return value, nil
}

// GetStreamLoginData 流式传输传递ctx, 获取loginData
func GetStreamLoginData(ctx context.Context) (int64, error) {
	uid, success := streamFromContext(ctx, constants.LoginDataKey)
	if !success {
		return -1, errno.NewErrNo(errno.ParamMissingErrorCode, "Failed to get info in context")
	}

	value, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return -1, errno.NewErrNo(errno.InternalServiceErrorCode, "Failed to get info in context when parse loginData")
	}
	return value, nil
}

// SetStreamLoginData 流式传输传递ctx, 设置ctx值
func SetStreamLoginData(ctx context.Context, uid int64) context.Context {
	value := strconv.FormatInt(uid, 10)
	return streamAppendContext(ctx, constants.LoginDataKey, value)
}
