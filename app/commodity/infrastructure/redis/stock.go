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
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

func (c *commodityCache) GetLockStockNum(ctx context.Context, key string) (int64, error) {
	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return 0, errno.Errorf(errno.InternalRedisErrorCode, "CommodityCache.GetLockStockNum failed :%v", err)
	}
	ret, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errno.Errorf(errno.InternalServiceErrorCode, "CommodityCache.GetLockStockNum failed :%v", err)
	}
	return ret, nil
}

func (c *commodityCache) SetLockStockNum(ctx context.Context, key string, num int64) {
	err := c.client.Set(ctx, key, num, constants.RedisLockStockExpireTime).Err()
	if err != nil {
		logger.Errorf("CommodityCache.SetLockStockNum failed :%v", err)
	}
}

func (c *commodityCache) IncrLockStockNum(ctx context.Context, infos []*model.SkuBuyInfo) error {
	_, err := c.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, info := range infos {
			key := c.GetLockStockKey(info.SkuID)
			err := pipe.IncrBy(ctx, key, info.Count).Err()
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return errno.Errorf(errno.InternalRedisErrorCode, "CommodityCache.IncrLockStockNum failed: %v", err)
	}
	return nil
}

func (c *commodityCache) DecrLockStockNum(ctx context.Context, infos []*model.SkuBuyInfo) error {
	// redis的事务是到最后一个操作的时候才执行全部操作的，Get这种获得值的操作就会为空值，使用Watch检测指定的key就可以避免这种情况，所以要先获取所有的key
	keys := make([]string, 0)
	for _, info := range infos {
		keys = append(keys, c.GetLockStockKey(info.SkuID))
	}
	err := c.client.Watch(ctx, func(tx *redis.Tx) error {
		for i := 0; i < len(infos); i++ {
			val, err := tx.Get(ctx, keys[i]).Int64()
			if err != nil {
				return errno.Errorf(errno.InternalRedisErrorCode, "CommodityCache.DecrLockStockNum failed :%v", err)
			}

			if val < 0 || val-infos[i].Count < 0 {
				return errno.NewErrNo(errno.InsufficientStockErrorCode, "CommodityCache.DecrLockStockNum failed: too many goods")
			}

			err = tx.DecrBy(ctx, keys[i], infos[i].Count).Err()
			if err != nil {
				return errno.Errorf(errno.InternalRedisErrorCode, "CommodityCache.DecrLockStockNum failed :%v", err)
			}
			err = tx.Expire(ctx, keys[i], constants.RedisLockStockExpireTime).Err()
			if err != nil {
				return errno.Errorf(errno.InternalRedisErrorCode, "CommodityCache.DecrLockStockNum failed :%v", err)
			}
		}
		return nil
	}, keys...)
	if err != nil {
		return err
	}
	return err
}

func (c *commodityCache) DecrStockNum(ctx context.Context, infos []*model.SkuBuyInfo) error {
	stockKeys := make([]string, 0)
	lockStockKeys := make([]string, 0)
	for _, info := range infos {
		stockKeys = append(stockKeys, c.GetStockKey(info.SkuID))
		lockStockKeys = append(lockStockKeys, c.GetLockStockKey(info.SkuID))
	}
	combinedKeys := append(stockKeys, lockStockKeys...)
	err := c.client.Watch(ctx, func(tx *redis.Tx) error {
		for i := 0; i < len(infos); i++ {
			val, err := tx.Get(ctx, stockKeys[i]).Int64()
			if err != nil {
				return errno.Errorf(errno.InternalRedisErrorCode, "CommodityCache.DecrStockNum failed :%v", err)
			}

			if val <= 0 || val-infos[i].Count < 0 {
				return errno.NewErrNo(errno.InsufficientStockErrorCode, "CommodityCache.DecrStockNum failed: too many goods for stock")
			}

			lockVal, err := tx.Get(ctx, lockStockKeys[i]).Int64()
			if err != nil {
				return errno.Errorf(errno.InternalRedisErrorCode, "CommodityCache.DecrStockNum failed :%v", err)
			}

			if lockVal <= 0 || lockVal-infos[i].Count < 0 || val <= lockVal {
				return errno.NewErrNo(errno.InsufficientStockErrorCode, "CommodityCache.DecrStockNum failed: too many goods for lock stock")
			}

			err = tx.DecrBy(ctx, stockKeys[i], infos[i].Count).Err()
			if err != nil {
				return errno.Errorf(errno.InternalRedisErrorCode, "CommodityCache.DecrStockNum failed :%v", err)
			}

			err = tx.Expire(ctx, stockKeys[i], constants.RedisStockExpireTime).Err()
			if err != nil {
				return errno.Errorf(errno.InternalRedisErrorCode, "CommodityCache.DecrStockNum failed :%v", err)
			}

			err = tx.DecrBy(ctx, lockStockKeys[i], infos[i].Count).Err()
			if err != nil {
				return errno.Errorf(errno.InternalRedisErrorCode, "CommodityCache.DecrStockNum failed :%v", err)
			}

			err = tx.Expire(ctx, lockStockKeys[i], constants.RedisLockStockExpireTime).Err()
			if err != nil {
				return errno.Errorf(errno.InternalRedisErrorCode, "CommodityCache.DecrStockNum failed :%v", err)
			}
		}
		return nil
	}, combinedKeys...)
	if err != nil {
		return err
	}
	return err
}
