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

package service

import (
	"github.com/redis/go-redis/v9"

	"github.com/west2-online/DomTok/app/payment/domain/repository"
	"github.com/west2-online/DomTok/pkg/utils"
)

type PaymentService struct {
	db repository.PaymentDB
	sf *utils.Snowflake
	// 是不是还要把redis的加进去
	redisClient *redis.Client
	// emailRe *regexp.Regexp
}

func NewPaymentService(db repository.PaymentDB, sf *utils.Snowflake) *PaymentService {
	if db == nil {
		panic("paymentService`s db should not be nil")
	}
	if sf == nil {
		panic("paymentService`s sf should not be nil")
	}
	svc := &PaymentService{db: db}
	// TODO redis的初始化放在哪里？
	//svc.init() 我需要写这个吗？
	return svc
}

func (svc *PaymentService) init() {
	//TODO redis的初始化放在这里吗？
	//svc.redisClient = redis.NewClient(&redis.Options{
	//Addr: "localhost:6379", // Redis 服务器地址
	//DB:   0,                // 默认使用 0 号数据库
	//})
	return
}
