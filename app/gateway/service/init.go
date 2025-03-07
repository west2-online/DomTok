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

import "context"

type GateWayService struct {
	Re *RedisService
	Bf *BF
}

func NewGateWayService() *GateWayService {
	svc := &GateWayService{
		Re: NewRedisService(),
		Bf: NewBloomFilter(),
	}
	svc.initBloomFilter()
	return svc
}

func (svc *GateWayService) initBloomFilter() {
	keys := svc.Re.GetAllBanedUser(context.Background())
	for _, key := range keys {
		svc.Bf.Add(key)
	}
}
