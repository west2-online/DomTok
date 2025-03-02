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

	metainfoContext "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/errno"
)

func (u *UseCase) DeleteCartGoods(ctx context.Context) error {
	userID, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		return fmt.Errorf("DeleteCartGoods get user info error: %w", err)
	}
	e, _, err := u.DB.GetCartByUserId(ctx, userID)
	if err != nil {
		return fmt.Errorf("DeleteCartGoods get cart by user id error: %w", err)
	}
	if !e {
		return errno.Errorf(errno.InvalidDeleteCartCode, "cart not exist")
	}
	err = u.DB.DeleteCart(ctx, userID)
	if err != nil {
		return fmt.Errorf("DeleteCartGoods delete cart error: %w", err)
	}
	return nil
}
