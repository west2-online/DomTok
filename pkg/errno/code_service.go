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

// 业务强相关, 范围是 1000-9999
// User
const (
	ServiceWrongPassword = 1000 + iota
	ServiceUserExist
	ServiceUserNotExist
	AddressNotExist

	ErrRecordNotFound
	UserLogOut
	UserBaned
	UserAlreadyLogin
)

// order
const (
	ServiceOrderNotFound = 2000 + iota
	UnknownOrderStatus
	OrderShouldNotBeChange
	ServiceOrderExpired
	ServiceOrderStatusInvalid
)

// commodity
const (
	ServiceSpuNotExist = 3000 + iota
	ServiceImgNotExist
	ServiceSkuExist

	ServiceSkuNotExist
	ServiceSkuImageNotExist
	ServiceSkuAttrNotExist

	ServiceCategoryExist
	ServiceCategorynotExist
	ServiceListCategoryFailed
	ServiceCategoryCreateFail
)

// payment
const (
	ServicePaymentOrderNotExist = 4000 + iota
	ServicePaymentRefundNotExist
	ServiceStorePaymentRedisTokenFailed
	ServiceCreatePaymentIDFailed
	ServiceCreatePaymentOrderFailed
	ServiceGetLoginDataFailed
	ServiceCheckOrderExistFailed
	ServiceGenerateHMACFailed
	ServiceRedisTimeLimited
	ServiceCreateRefundIDFailed
	ServiceCreateRefundFailed
	ServicePaymentIsProcessing
)

// cart
const (
	InvalidDeleteCartCode = 5000 + iota
)

// assistant
const (
	ServiceUserCloseWebsocketConn = 8000 + iota
)
