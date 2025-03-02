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

	"github.com/redis/go-redis/v9"

	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/logger"
)

var _cli *paymentRedis

func _EnvSetup() {
	if _cli != nil {
		err := _cli.client.FlushDB(context.Background()).Err()
		if err != nil {
			panic(err)
		}
		return
	}
	logger.Ignore()
	config.Init(constants.PaymentServiceName)
	r := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
	})
	err := r.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	err = r.FlushDB(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	_cli, _ = NewPaymentRedis(r).(*paymentRedis)
}
