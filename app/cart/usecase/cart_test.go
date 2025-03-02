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
	"errors"
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/cart/domain/model"
	"github.com/west2-online/DomTok/app/cart/domain/service"
	"github.com/west2-online/DomTok/app/cart/infrastructure/db"
	"github.com/west2-online/DomTok/app/cart/infrastructure/mq"
	"github.com/west2-online/DomTok/app/cart/infrastructure/rpc"
	metainfoContext "github.com/west2-online/DomTok/pkg/base/context"
)

func TestUseCase_AddGoodsIntoCart(t *testing.T) {
	type TestCase struct {
		Name string
		// Mock返回的错误
		MockVerifyError error
		MockLoginError  error
		MockSendError   error
		// 期望的错误
		ExpectedError error
	}

	testCases := []TestCase{
		{
			Name:            "VerifyError",
			MockVerifyError: errors.New("verify count error"),
			// 如果校验不通过，后续不会执行其他逻辑
			ExpectedError: errors.New("verify count error"),
		},
		{
			Name:           "GetLoginDataError",
			MockLoginError: errors.New("get login data error"),
			ExpectedError:  fmt.Errorf("cartCase.AddGoodsIntoCart metainfo unmarshal error:%w", errors.New("get login data error")),
		},
		{
			Name:          "SendAddGoodsError",
			MockSendError: errors.New("send mq error"),
			ExpectedError: fmt.Errorf("cartCase.AddGoodsIntoCart send mq error:%w", errors.New("send mq error")),
		},
		{
			Name:          "AddGoodsIntoCartSuccess",
			ExpectedError: nil,
		},
	}

	// 定义一个测试的 GoodInfo
	goods := &model.GoodInfo{
		SkuId:     123,
		ShopId:    456,
		VersionId: 1,
		Count:     5,
	}

	// 清理 mock
	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// 1. 构造依赖并进行 mock
			svcMock := new(service.CartService)
			d := new(db.DBAdapter)
			mqMock := new(mq.KafkaAdapter)
			uc := &UseCase{
				svc: svcMock,
				DB:  d,
				MQ:  mqMock,
			}

			mockey.Mock((*service.CartService).Verify).Return(tc.MockVerifyError).Build()
			mockey.Mock(metainfoContext.GetLoginData).Return(int64(101), tc.MockLoginError).Build()
			mockey.Mock(mockey.GetMethod(uc.MQ, "SendAddGoods")).Return(tc.MockSendError).Build()

			err := uc.AddGoodsIntoCart(context.Background(), goods)

			if err != nil && tc.ExpectedError != nil {
				// 比较错误字符串
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				// Either no error was expected or none was returned
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
		})
	}
}

func TestUseCase_ShowCartGoods(t *testing.T) {
	type TestCase struct {
		Name string
		// 1) svc.Verify(pageNum)
		MockVerifyPageNumError error
		// 2) metainfoContext.GetLoginData
		MockLoginUID   int64
		MockLoginError error
		// 3) svc.TryGetCartFromCache
		MockCacheExist bool
		MockCacheRes   []*model.CartGoods
		MockCacheError error
		// 4) DB.GetCartByUserId
		MockDBExist  bool
		MockCartData *model.Cart
		MockDBError  error
		// 5) sonic.Unmarshal  (如需强行 mock，需要封装或给 sonic.Unmarshal 再包一层)
		MockUnmarshalError error
		// 6) model.ConvertCartJsonToCartGoods
		MockConvertRes []*model.CartGoods
		// 7) Rpc.GetGoodsInfo
		MockGoodsError error
		MockGoodsRes   []*model.CartGoods
		// 8) svc.TrySetCartCache
		//    在 go routine 中执行，只打印日志，不影响主流程
		MockSetCacheErr error

		// 期望结果
		ExpectedErr  error
		ExpectedCart []*model.CartGoods
	}

	testCases := []TestCase{
		{
			Name:                   "VerifyPageNumError",
			MockVerifyPageNumError: errors.New("page num invalid"),
			ExpectedErr:            errors.New("page num invalid"),
		},
		{
			Name:           "GetLoginDataError",
			MockLoginError: errors.New("login error"),
			ExpectedErr:    fmt.Errorf("ShowCartGoods get user info error: %w", errors.New("login error")),
		},
		{
			Name:           "TryGetCartFromCacheError",
			MockCacheError: errors.New("cache error"),
			ExpectedErr:    fmt.Errorf("ShowCartGoods get cart from cache error: %w", errors.New("cache error")),
		},
		{
			Name:           "CacheExistReturnResult",
			MockCacheExist: true,
			MockCacheRes: []*model.CartGoods{
				{GoodsID: 1001, GoodsName: "CachedGoods"},
			},
			ExpectedCart: []*model.CartGoods{
				{GoodsID: 1001, GoodsName: "CachedGoods"},
			},
		},
		{
			Name:        "DBGetCartError",
			MockDBError: errors.New("db error"),
			ExpectedErr: fmt.Errorf("ShowCartGoods DB get cart error: %w", errors.New("db error")),
		},
		{
			Name:        "DBNotExist",
			MockDBExist: false,
			// 数据不存在，返回空切片
			ExpectedCart: []*model.CartGoods{},
		},
		{
			Name:        "GetGoodsInfoError",
			MockDBExist: true,
			MockCartData: &model.Cart{
				SkuJson: `{"store":[{"store_id":1,"sku":[{"sku_id":1002,"version_id":2,"count":2}]}]}`,
			},
			// Convert 成功后
			MockConvertRes: []*model.CartGoods{
				{GoodsID: 1002, PurchaseQuantity: 2},
			},
			// RPC 出错
			MockGoodsError: errors.New("rpc error"),
			ExpectedErr:    fmt.Errorf("ShowCartGoods RPC error: %w", errors.New("rpc error")),
		},
		{
			Name:        "SuccessButSetCacheError",
			MockDBExist: true,
			MockCartData: &model.Cart{
				SkuJson: `{"store":[{"store_id":2,"sku":[{"sku_id":2001,"version_id":1,"count":3}]}]}`,
			},
			MockConvertRes: []*model.CartGoods{
				{GoodsID: 2001, PurchaseQuantity: 3},
			},
			MockGoodsRes: []*model.CartGoods{
				{GoodsID: 2001, GoodsName: "FinalGoods"},
			},
			MockSetCacheErr: errors.New("cache set error"),
			// 不影响主流程，仍旧成功返回
			ExpectedCart: []*model.CartGoods{
				{GoodsID: 2001, GoodsName: "FinalGoods"},
			},
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// 1. 构造依赖实例
			svcMock := new(service.CartService)
			dbMock := new(db.DBAdapter)
			rpcMock := new(rpc.CartRpcImpl)
			uc := &UseCase{
				svc: svcMock,
				DB:  dbMock,
				Rpc: rpcMock,
			}

			// 2. mock svc.Verify(pageNum)
			mockey.Mock(mockey.GetMethod(svcMock, "Verify")).
				Return(tc.MockVerifyPageNumError).Build()

			// 3. mock metainfoContext.GetLoginData
			mockey.Mock(metainfoContext.GetLoginData).Return(int64(101), tc.MockLoginError).Build()

			// 4. mock svc.TryGetCartFromCache
			mockey.Mock(mockey.GetMethod(svcMock, "TryGetCartFromCache")).
				Return(tc.MockCacheExist, tc.MockCacheRes, tc.MockCacheError).Build()

			// 5. mock db.GetCartByUserId
			mockey.Mock(mockey.GetMethod(dbMock, "GetCartByUserId")).
				Return(tc.MockDBExist, tc.MockCartData, tc.MockDBError).Build()

			// 8. mock rpc.GetGoodsInfo
			mockey.Mock(mockey.GetMethod(rpcMock, "GetGoodsInfo")).
				Return(tc.MockGoodsRes, tc.MockGoodsError).Build()

			// 9. mock svc.TrySetCartCache (在 go routine 中执行)
			mockey.Mock(mockey.GetMethod(svcMock, "TrySetCartCache")).
				Return(tc.MockSetCacheErr).Build()

			// 开始测试
			res, err := uc.ShowCartGoods(context.Background(), 1)

			// 如果预期或实际存在错误，就仅对错误做断言
			if err != nil || tc.ExpectedErr != nil {
				if err != nil && tc.ExpectedErr != nil {
					convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedErr.Error())
				} else {
					convey.So(err, convey.ShouldEqual, tc.ExpectedErr)
				}
				return
			}
			// 如果无错误，则对结果做断言
			convey.So(res, convey.ShouldResemble, tc.ExpectedCart)
		})
	}
}
