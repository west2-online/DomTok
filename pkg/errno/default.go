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

// 一些常用的的错误, 如果你懒得单独定义也可以直接使用
var (
	Success = NewErrNo(SuccessCode, "ok")

	ParamVerifyError  = NewErrNo(ParamVerifyErrorCode, "parameter validation failed")
	ParamMissingError = NewErrNo(ParamMissingErrorCode, "missing parameter")

	AuthInvalid             = NewErrNo(AuthInvalidCode, "authentication failure")
	AuthAccessExpired       = NewErrNo(AuthAccessExpiredCode, "token expiration")
	AuthNoToken             = NewErrNo(AuthNoTokenCode, "lack of token")
	AuthNoOperatePermission = NewErrNo(AuthNoOperatePermissionCode, "No permission to operate")

	InternalServiceError = NewErrNo(InternalServiceErrorCode, "internal server error")
	OSOperationError     = NewErrNo(OSOperateErrorCode, "os operation failed")
	IOOperationError     = NewErrNo(IOOperateErrorCode, "io operation failed")

	UpYunFileError = NewErrNo(UpYunFileErrorCode, "upyun operation failed")
)
