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

package constants

import (
	"fmt"
	"time"
)

const (
	TypeAccessToken  = 0
	TypeRefreshToken = 1

	AccessTokenTTL  = time.Hour * 24 * 7  // Access Token 有效期7天
	RefreshTokenTTL = time.Hour * 24 * 30 // Refresh Token 有效期30天
	Issuer          = "west2-online"      // token 颁发者

	AuthHeader         = "Authorization" // 获取 Token 时的请求头
	AccessTokenHeader  = "Access-Token"  // 响应时的访问令牌头
	RefreshTokenHeader = "Refresh-Token" // 响应时的刷新令牌头
)

var PublicKey = fmt.Sprintf("%v\n%v\n%v", "-----BEGIN PUBLIC KEY-----",
	"MCowBQYDK2VwAyEAT+ypuz7wIltf8HoFUEI/rDBrQNhZShqLv88j4aAWnT0=",
	"-----END PUBLIC KEY-----")
