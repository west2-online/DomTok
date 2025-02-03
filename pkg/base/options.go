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

package base

import (
	cli "github.com/cloudwego/hertz/pkg/app/client"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// WithRedisClient will create redis object
func WithRedisClient(dbName int) Option {
	return func(clientSet *ClientSet) {
		redisClient, err := client.NewRedisClient(dbName)
		if err != nil {
			logger.Fatalf("init cache failed, err: %v", err)
		}
		clientSet.CacheClient = cache.NewCache(redisClient)
		clientSet.cleanups = append(clientSet.cleanups, func() {
			err = redisClient.Close()
			if err != nil {
				logger.Errorf("close cache failed, err: %v", err)
			}
		})
	}
}

// WithDBClient will create database object
func WithDBClient(tableName string) Option {
	return func(clientSet *ClientSet) {
		DB, err := client.InitMySQL(tableName)
		if err != nil {
			logger.Fatalf("init database failed, err: %v", err)
		}

		// TODO: currently our service only deploy on one or some servers, we do not need specific datacenterID
		sf, err := utils.NewSnowflake(config.Snowflake.DatancenterID, config.Snowflake.WorkerID)
		if err != nil {
			logger.Fatalf("init Snowflake object error: %v", err.Error())
		}
		clientSet.SFClient = sf
		clientSet.DBClient = db.NewDatabase(DB, sf)
		// gorm maintains a connection pool.
		// After initialization, all connections are managed by gorm and do not need to be closed manually.
	}
}

func WithElasticSearch() Option {
	return func(clientSet *ClientSet) {
		es, err := client.NewEsClient()
		if err != nil {
			logger.Fatalf("init elastic search client error: %v", err.Error())
		}
		clientSet.ESClient = es
	}
}

func WithHzClient() Option {
	return func(clientSet *ClientSet) {
		hz, err := cli.NewClient()
		if err != nil {
			logger.Fatalf("init Hertz client error: %v", err)
		}
		clientSet.HzClient = hz
	}
}
