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

package errno

// 错误码的设计原则是, 以尽可能少的错误码来传递必要的信息,
// 让前端能够根据尽量少的 error code 和具体的场景来告知用户错误信息
// 总的来说前端不依赖于后端传递的 msg 来告知用户, 而是通过 code 来额外处理
// 当然如果有一些强指向性错误信息, 你当然可以再写进来一个 code, 比如密码错误或者用户已存在
// 我们将这种与业务强相关的 code 也放在 errno 包中, 主要是为了方便统一管理与避免 code 冲突

const (
	// SuccessCode For microservices
	SuccessCode = 10000
	SuccessMsg  = "ok"
)

// 200xx: 参数错误，Param 打头
const (
	ParamVerifyErrorCode   = 20000 + iota // 参数校验失败
	ParamMissingErrorCode                 // 参数缺失
	ParamMissingHeaderCode                // 请求头缺失
	ParamInvalidHeaderCode                // 请求头无效
)

// 300xx: 鉴权错误，Auth 打头
const (
	AuthInvalidCode             = 30000 + iota // 鉴权失败
	AuthAccessExpiredCode                      // 访问令牌过期
	AuthRefreshExpiredCode                     // 刷新令牌过期
	AuthNoTokenCode                            // 没有 token
	AuthNoOperatePermissionCode                // 没有操作权限
	AuthMissingTokenCode                       // 缺少 token
	IllegalOperatorCode                        // 不合格的操作(比如传入 payment status时传入了一个不存在的 status)
)

// 500xx: 内部错误，Internal 打头
// 服务级别的错误, 发生的时候说明我们程序自身出了问题
// 比如数据库断联, 编码错误等. 需要我们人为的去维护
const (
	InternalServiceErrorCode  = 50000 + iota // 内部服务错误
	InternalDatabaseErrorCode                // 数据库错误
	InternalRedisErrorCode                   // Redis错误
	InternalNetworkErrorCode                 // 网络错误
	InternalESErrorCode                      // ES错误
	InternalKafkaErrorCode                   // kafka 错误
	OSOperateErrorCode
	IOOperateErrorCode
	InsufficientStockErrorCode
	InternalRPCErrorCode
	InternalRocketmqErrorCode
)

const (
	UpYunFileErrorCode = 60000 + iota
	RedisKeyNotExist
	RepeatedOperation
)
