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

import "time"

const (
	// 请求体最大体积
	ServerMaxRequestBodySize = 1 << 31

	CorsMaxAge = 12 * time.Hour

	SentinelThreshold        = 100
	SentinelStatIntervalInMs = 1000
	LoginDataKey             = "loginData"

	BloomFilterSize   = 10000
	FalsePositiveRate = 0.05
	BFMaxSize         = 20
)

const (
	UserLogout = iota
	UserBanned
)
