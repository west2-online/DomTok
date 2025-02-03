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
	MuxConnection    = 1                     // (RPC) 最大连接数
	StreamBufferSize = 1024                  // (RPC) 流请求 Buffer 尺寸
	MaxQPS           = 100                   // (RPC) 最大 QPS
	RPCTimeout       = 3 * time.Second       // (RPC) RPC请求超时时间
	ConnectTimeout   = 50 * time.Millisecond // (RPC) 连接超时时间
)
