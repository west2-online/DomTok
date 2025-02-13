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

package usecase

import (
	"context"
	"fmt"
	"github.com/west2-online/DomTok/app/cart/domain/model"
	"github.com/west2-online/DomTok/pkg/constants"
)

func (u *UseCase) AddGoodsIntoCart(ctx context.Context, goods *model.GoodInfo) (err error) {
	if err = u.svc.Verify(u.svc.VerifyCount(goods.Count)); err != nil {
		return
	}

	// todo: 开启metainfo透传
	/*
		loginData, err := metainfoContext.GetLoginData(ctx)
		if err != nil {
			return fmt.Errorf("cartCase.AddGoodsIntoCart metainfo unmarshal error:%w", err)
		}

	*/
	err = u.MQ.SendAddGoods(ctx, constants.UserTestId, goods)
	if err != nil {
		return fmt.Errorf("cartCase.AddGoodsIntoCart send mq error:%w", err)
	}

	return nil
}
