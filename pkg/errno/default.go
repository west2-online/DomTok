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

// This file is be designed to define any common error so that we can use it in any service simply.

package errno

import "github.com/cloudwego/hertz/pkg/protocol/consts"

var (
	Success                   = NewErrNo(SuccessCode, "ok")
	CustomLaunchScreenSuccess = NewErrNo(consts.StatusOK, "ok") // 兼容处理

	AuthError          = NewErrNo(AuthErrorCode, "鉴权失败")            // 鉴权失败，通常是内部错误，如解析失败
	AuthInvalid        = NewErrNo(AuthInvalidCode, "鉴权无效")          // 鉴权无效，如令牌颁发者不是 west2-online
	AuthAccessExpired  = NewErrNo(AuthAccessExpiredCode, "访问令牌过期")  // 访问令牌过期
	AuthRefreshExpired = NewErrNo(AuthRefreshExpiredCode, "刷新令牌过期") // 刷新令牌过期
	AuthMissing        = NewErrNo(AuthInvalidCode, "缺失合法鉴权数据")      // 鉴权缺失，如访问令牌缺失

	ParamError         = NewErrNo(ParamErrorCode, "参数错误") // 参数校验失败，可能是参数为空、参数类型错误等
	ParamMissingHeader = NewErrNo(ParamMissingHeaderCode, "缺失合法学生请求头数据")

	BizError             = NewErrNo(BizErrorCode, "请求业务出现问题")
	InternalServiceError = NewErrNo(InternalServiceErrorCode, "内部服务错误")

	SuffixError           = NewErrNo(ParamErrorCode, "文件不可用")
	NoRunningPictureError = NewErrNo(BizErrorCode, "没有可用图片")
	NoMatchingPlanError   = NewErrNo(BizErrorCode, "没有匹配的计划")

	// internal error
	UpcloudError = NewErrNo(BizFileUploadErrorCode, "云服务商交互错误")

	// redis
	RedisError = NewErrNo(InternalRedisErrorCode, "缓存服务出现问题")
)
