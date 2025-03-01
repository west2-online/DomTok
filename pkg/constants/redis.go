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
	RedisSlowQuery = 10 // ms redis默认的慢查询时间，适用于 logger
)

// Redis Key and Expire Time
const (
	RedisCartExpireTime      = 5 * 60 * time.Second
	RedisCartStoreNum        = 30
	RedisSpuImageExpireTime  = 5 * 60 * time.Second

	RedisSkuExpireTime = 5 * 60 * time.Second
	RedisLockStockExpireTime = 24 * 60 * 60 * time.Second
	RedisStockExpireTime     = 24 * 60 * 60 * time.Second
	RedisNXExpireTime        = 3 * time.Second
	RedisMaxLockRetryTime    = 400 * time.Millisecond
	RedisRetryStopTime       = 100 * time.Millisecond
)

// Redis DB Name
const (
	RedisDBOrder     = 0
	RedisDBCommodity = 1
	RedisDBCart      = 2

	RedSyncDBId = 0
)

// Redis Connection Pool Configuration
const (
	RedisPoolSize           = 50              // 最大连接数
	RedisMinIdleConnections = 10              // 最小空闲连接数
	RedisDialTimeout        = 5 * time.Second // 连接超时时间
)

// Order
const (
	OrderCacheOrderExpireFormat   = "order-expire-%d"
	OrderCachePaymentStatusFormat = "payment-status-%d"
	OrderPaymentStatusExpireTime  = 12 * time.Minute // 大于 SkuStockRollbackTopicDelayTimeLevel 即可
	OrderCacheLuaKeyExistFlag     = 1
	OrderRedSyncDefaultTTL        = 8 * time.Second
	OrderRedSyncDefaultInterval   = 1 * time.Second

	OrderLockFormat                   = "lock-order-%d"
	OrderUpdatePaymentStatusLuaScript = `
        local exK = KEYS[1]
        local sK = KEYS[2]
        local orderExpire = ARGV[1]
        local status = ARGV[2]
        local expire = tonumber(ARGV[3])

        local exist = redis.call('EXISTS', exK)
        if exist == 0 then
            return 0
        end

        redis.call('SET', exK, orderExpire, 'EX', expire)
        redis.call('SET', sK, status, 'EX', expire)

        return 1
    `
)
