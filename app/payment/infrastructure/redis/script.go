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

package redis

import (
	"context"
	"fmt"
)

type ScriptKey string

const (
	CheckAndDelScript  ScriptKey = "CheckAndDel"
	GetTTLAndDelScript ScriptKey = "GetTTLAndDel"
)

type _Script struct {
	Hash string
	Cmd  string
}

var scripts = map[ScriptKey]*_Script{
	CheckAndDelScript: {
		Cmd: `
local exists = redis.call("EXISTS", KEYS[1])
if exists == 1 and redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("DEL", KEYS[1])
	return 1
end
return 0
`,
	},
	GetTTLAndDelScript: {
		Cmd: `
local exists = redis.call("EXISTS", KEYS[1])
if exists == 1 and redis.call("GET", KEYS[1]) == ARGV[1] then
	local ttl = redis.call("TTL", KEYS[1])
	redis.call("DEL", KEYS[1])
	return {ttl, 1}
end
return {-1, 0}
`,
	},
}

// loadScript 加载脚本，并将哈希值存储在内存中
func (p *paymentRedis) loadScript() (err error) {
	ctx := context.Background()
	for key, value := range scripts {
		hash, err := p.client.ScriptLoad(ctx, value.Cmd).Result()
		if err != nil {
			return fmt.Errorf("load script %s failed: %w", key, err)
		}
		value.Hash = hash
	}
	return nil
}

// execScript 执行脚本
func (p *paymentRedis) execScript(ctx context.Context, k ScriptKey, keys []string, args ...interface{}) (interface{}, error) {
	script := scripts[k]
	return p.client.EvalSha(ctx, script.Hash, keys, args...).Result()
}
