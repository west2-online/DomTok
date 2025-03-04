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

import "github.com/west2-online/DomTok/pkg/errno"

type CartVerifyOps func() error

// Verify 通过传来的参数进行一系列的校验
func (svc *CartService) Verify(opts ...CartVerifyOps) error {
	for _, opt := range opts {
		if err := opt(); err != nil {
			return err
		}
	}
	return nil
}

func (svc *CartService) VerifyCount(cnt int64) CartVerifyOps {
	return func() error {
		if cnt < 1 {
			return errno.NewErrNo(errno.ParamVerifyErrorCode, "wrong goods count format")
		}
		return nil
	}
}

func (svc *CartService) VerifyPageNum(p int64) CartVerifyOps {
	return func() error {
		if p < 1 {
			return errno.NewErrNo(errno.ParamVerifyErrorCode, "wrong goods count format")
		}
		return nil
	}
}
