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
	ctx "context"
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/olivere/elastic/v7"
	"github.com/redis/go-redis/v9"
	"github.com/smartystreets/goconvey/convey"
	"gorm.io/gorm"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/domain/service"
	"github.com/west2-online/DomTok/app/commodity/infrastructure/es"
	"github.com/west2-online/DomTok/app/commodity/infrastructure/mysql"
	redisCommodity "github.com/west2-online/DomTok/app/commodity/infrastructure/redis"
	"github.com/west2-online/DomTok/kitex_gen/commodity"
	"github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/utils"
)

func TestUseCase_CreateSpu(t *testing.T) {
	type TestCase struct {
		Name            string
		MockUid         int64
		MockCtxError    error
		MockVerifyError error
		MockCreateError error
		MockSpuId       int64
		ExpectedError   error
		ExpectedSpuId   int64
	}

	testcase := []TestCase{
		{
			Name:          "GetUidError",
			MockUid:       1,
			ExpectedSpuId: 0,
			MockCtxError:  errors.New("GetUidError"),
			ExpectedError: errors.New("usecase.CreateSpu failed: GetUidError"),
		},
		{
			Name:            "VerifyError",
			MockVerifyError: errors.New("VerifyError"),
			ExpectedError:   errors.New("usecase.CreateSpu verify failed: VerifyError"),
			ExpectedSpuId:   0,
		},
		{
			Name:            "CreateError",
			MockCreateError: errors.New("CreateError"),
			ExpectedError:   errors.New("usecase.CreateSpu failed: CreateError"),
			ExpectedSpuId:   0,
		},
		{
			Name:            "CreateSpuSuccessfully",
			MockUid:         0,
			MockCtxError:    nil,
			MockVerifyError: nil,
			MockCreateError: nil,
			MockSpuId:       1,
			ExpectedError:   nil,
			ExpectedSpuId:   1,
		},
	}

	spu := &model.Spu{
		Name: "OPPO A93s",
	}

	defer mockey.UnPatchAll()
	for _, tc := range testcase {
		mockey.PatchConvey(tc.Name, t, func() {
			mockey.Mock(context.GetStreamLoginData).Return(tc.MockUid, tc.MockCtxError).Build()
			mockey.Mock((*service.CommodityService).Verify).Return(tc.MockVerifyError).Build()
			mockey.Mock((*service.CommodityService).CreateSpu).Return(tc.MockSpuId, tc.MockCreateError).Build()
			us := new(useCase)
			svc := new(service.CommodityService)
			us.svc = svc

			id, err := us.CreateSpu(ctx.Background(), spu)
			if err != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
			convey.So(id, convey.ShouldEqual, tc.ExpectedSpuId)
		})
	}
}

func TestUseCase_CreateSpuImage(t *testing.T) {
	type TestCase struct {
		Name            string
		MockGetSpuError error
		MockCreateImage error
		MockImageId     int64
		ExpectedError   error
		ExpectedImageId int64
	}

	testcase := []TestCase{
		{
			Name:            "GetSpuError",
			MockGetSpuError: errors.New("GetSpuError"),
			ExpectedError:   errors.New("usecase.CreateSpuImage failed: GetSpuError"),
			ExpectedImageId: 0,
		},
		{
			Name:            "CreateSpuImageError",
			MockCreateImage: errors.New("CreateSpuImageError"),
			ExpectedError:   errors.New("usecase.CreateSpuImage failed: CreateSpuImageError"),
			ExpectedImageId: 0,
		},
		{
			Name:            "CreateSpuImageSuccessfully",
			MockCreateImage: nil,
			MockGetSpuError: nil,
			MockImageId:     1,
			ExpectedError:   nil,
			ExpectedImageId: 1,
		},
	}
	img := &model.SpuImage{
		SpuID: 1,
		Url:   "http://example.jpg",
	}

	defer mockey.UnPatchAll()

	for _, tc := range testcase {
		mockey.PatchConvey(tc.Name, t, func() {
			svc := new(service.CommodityService)

			gormDB := new(gorm.DB)
			db := mysql.NewCommodityDB(gormDB)
			us := &useCase{
				db:  db,
				svc: svc,
			}

			mockey.Mock(mockey.GetMethod(us.db, "GetSpuBySpuId")).Return(nil, tc.MockGetSpuError).Build()
			mockey.Mock((*service.CommodityService).CreateSpuImage).Return(tc.MockImageId, tc.MockCreateImage).Build()

			id, err := us.CreateSpuImage(ctx.Background(), img)
			if err != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
			convey.So(id, convey.ShouldEqual, tc.ExpectedImageId)
		})
	}
}

func TestUseCase_DeleteSpu(t *testing.T) {
	type TestCase struct {
		Name                          string
		MockMatchDeleteConditionError error
		MockIdentifyError             error
		MockDeleteSpuError            error
		MockDeleteAllSpuImagesError   error
		MockSpuInfo                   *model.Spu
		ExpectedError                 error
	}

	testcase := []TestCase{
		{
			Name:                          "MatchConditionFailed",
			MockMatchDeleteConditionError: errors.New("MatchConditionFailed"),
			MockSpuInfo:                   nil,
			ExpectedError:                 errors.New("usecase.DeleteSpu failed: MatchConditionFailed"),
		},
		{
			Name:                          "IdentifyError",
			MockMatchDeleteConditionError: nil,
			MockIdentifyError:             errors.New("IdentifyError"),
			MockSpuInfo:                   &model.Spu{SpuId: 1, CreatorId: 1},
			ExpectedError:                 errors.New("usecase.DeleteSpu identify user failed: IdentifyError"),
		},
		{
			Name:               "DeleteSpuError",
			MockDeleteSpuError: errors.New("DeleteSpuError"),
			MockSpuInfo:        &model.Spu{SpuId: 1, CreatorId: 1},
			ExpectedError:      errors.New("usecase.DeleteSpu failed: DeleteSpuError"),
		},
		{
			Name:                        "DeleteAllSpuImagesError",
			MockDeleteAllSpuImagesError: errors.New("DeleteAllSpuImagesError"),
			MockSpuInfo:                 &model.Spu{SpuId: 1, CreatorId: 1},
			ExpectedError:               errors.New("usecase.DeleteSpu failed: DeleteAllSpuImagesError"),
		},
		{
			Name:                          "DeleteSpuSuccessfully",
			MockSpuInfo:                   &model.Spu{SpuId: 1, CreatorId: 1},
			MockDeleteSpuError:            nil,
			MockDeleteAllSpuImagesError:   nil,
			MockIdentifyError:             nil,
			MockMatchDeleteConditionError: nil,
			ExpectedError:                 nil,
		},
	}

	defer mockey.UnPatchAll()

	var spuId int64 = 1
	for _, tc := range testcase {
		mockey.PatchConvey(tc.Name, t, func() {
			svc := new(service.CommodityService)
			us := &useCase{
				svc: svc,
			}

			mockey.Mock((*service.CommodityService).MatchDeleteSpuCondition).Return(tc.MockSpuInfo, tc.MockMatchDeleteConditionError).Build()
			mockey.Mock((*service.CommodityService).IdentifyUser).Return(tc.MockIdentifyError).Build()
			mockey.Mock((*service.CommodityService).DeleteSpu).Return(tc.MockDeleteSpuError).Build()
			mockey.Mock((*service.CommodityService).DeleteAllSpuImages).Return(tc.MockDeleteAllSpuImagesError).Build()

			err := us.DeleteSpu(ctx.Background(), spuId)
			if err != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
		})
	}
}

func TestUseCase_UpdateSpu(t *testing.T) {
	type TestCase struct {
		Name              string
		MockGetSpuError   error
		MockSpuInfo       *model.Spu
		MockIdentifyError error
		MockVerifyError   error
		MockUpdateError   error
		ExpectedError     error
	}

	testcase := []TestCase{
		{
			Name:            "GetSpuError",
			MockGetSpuError: errors.New("GetSpuError"),
			MockSpuInfo:     nil,
			ExpectedError:   errors.New("usecase.UpdateSpu failed: GetSpuError"),
		},
		{
			Name:              "IdentifyError",
			MockIdentifyError: errors.New("IdentifyError"),
			MockSpuInfo:       &model.Spu{SpuId: 1, CreatorId: 1},
			ExpectedError:     errors.New("usecase.UpdateSpu identify user failed: IdentifyError"),
		},
		{
			Name:            "VerifyError",
			MockVerifyError: errors.New("VerifyError"),
			MockSpuInfo:     &model.Spu{SpuId: 1, CreatorId: 1},
			ExpectedError:   errors.New("usecase.UpdateSpu verify failed: VerifyError"),
		},
		{
			Name:            "UpdateSpuError",
			MockUpdateError: errors.New("UpdateSpuError"),
			MockSpuInfo:     &model.Spu{SpuId: 1, CreatorId: 1},
			ExpectedError:   errors.New("usecase.UpdateSpu failed: UpdateSpuError"),
		},
		{
			Name:            "UpdateSpuSuccessfully",
			MockUpdateError: nil,
			MockSpuInfo:     &model.Spu{SpuId: 1, CreatorId: 1},
			ExpectedError:   nil,
		},
	}

	defer mockey.UnPatchAll()

	spu := &model.Spu{SpuId: 1, Description: "description"}

	for _, tc := range testcase {
		mockey.PatchConvey(tc.Name, t, func() {
			svc := new(service.CommodityService)
			gormDB := new(gorm.DB)
			db := mysql.NewCommodityDB(gormDB)
			us := &useCase{
				svc: svc,
				db:  db,
			}
			mockey.Mock(mockey.GetMethod(us.db, "GetSpuBySpuId")).Return(tc.MockSpuInfo, tc.MockGetSpuError).Build()
			mockey.Mock((*service.CommodityService).IdentifyUserInStreamCtx).Return(tc.MockIdentifyError).Build()
			mockey.Mock((*service.CommodityService).Verify).Return(tc.MockVerifyError).Build()
			mockey.Mock((*service.CommodityService).UpdateSpu).Return(tc.MockUpdateError).Build()
			mockey.Mock(utils.GenerateFileName).Return("").Build()

			err := us.UpdateSpu(ctx.Background(), spu)
			if err != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
		})
	}
}

func TestUseCase_UpdateSpuImage(t *testing.T) {
	type TestCase struct {
		Name              string
		MockGetSpuError   error
		MockSpuInfo       *model.Spu
		MockImageInfo     *model.SpuImage
		MockIdentifyError error
		MockUpdateError   error
		ExpectedError     error
	}

	testcase := []TestCase{
		{
			Name:            "GetSpuError",
			MockGetSpuError: errors.New("GetSpuError"),
			ExpectedError:   errors.New("usecase.UpdateSpuImage failed: GetSpuError"),
		},
		{
			Name:              "IdentifyError",
			MockIdentifyError: errors.New("IdentifyError"),
			ExpectedError:     errors.New("usecase.UpdateSpuImage identify user failed: IdentifyError"),
			MockSpuInfo:       &model.Spu{SpuId: 1, CreatorId: 1},
			MockImageInfo:     &model.SpuImage{ImageID: 1, SpuID: 1},
		},
		{
			Name:            "UpdateSpuError",
			MockUpdateError: errors.New("UpdateSpuError"),
			MockSpuInfo:     &model.Spu{SpuId: 1, CreatorId: 1},
			MockImageInfo:   &model.SpuImage{ImageID: 1, SpuID: 1},
			ExpectedError:   errors.New("usecase.UpdateSpuImage failed: UpdateSpuError"),
		},
		{
			Name:            "UpdateSpuSuccessfully",
			MockUpdateError: nil,
			MockSpuInfo:     &model.Spu{SpuId: 1, CreatorId: 1},
			MockImageInfo:   &model.SpuImage{ImageID: 1, SpuID: 1},
			ExpectedError:   nil,
		},
	}

	defer mockey.UnPatchAll()
	img := &model.SpuImage{ImageID: 1, SpuID: 1}
	for _, tc := range testcase {
		mockey.PatchConvey(tc.Name, t, func() {
			svc := new(service.CommodityService)
			us := &useCase{
				svc: svc,
			}

			mockey.Mock((*service.CommodityService).GetSpuFromImageId).Return(tc.MockSpuInfo, tc.MockImageInfo, tc.MockGetSpuError).Build()
			mockey.Mock((*service.CommodityService).IdentifyUserInStreamCtx).Return(tc.MockIdentifyError).Build()
			mockey.Mock((*service.CommodityService).UpdateSpuImage).Return(tc.MockUpdateError).Build()
			mockey.Mock(utils.GenerateFileName).Return("").Build()
			err := us.UpdateSpuImage(ctx.Background(), img)
			if err != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
		})
	}
}

func TestUseCase_DeleteSpuImage(t *testing.T) {
	type TestCase struct {
		Name                    string
		MockGetSpuError         error
		MockIdentifyError       error
		MockDeleteSpuImageError error
		MockSpuInfo             *model.Spu
		MockImageInfo           *model.SpuImage
		ExpectedError           error
	}

	testcase := []TestCase{
		{
			Name:            "GetSpuError",
			MockGetSpuError: errors.New("GetSpuError"),
			ExpectedError:   errors.New("usecase.DeleteSpuImage failed: GetSpuError"),
		},
		{
			Name:              "IdentifyError",
			MockIdentifyError: errors.New("IdentifyError"),
			MockSpuInfo:       &model.Spu{SpuId: 1, CreatorId: 1},
			MockImageInfo:     &model.SpuImage{ImageID: 1, SpuID: 1, Url: "http://example.jpg"},
			ExpectedError:     errors.New("usecase.DeleteSpuImage identify user failed: IdentifyError"),
		},
		{
			Name:                    "DeleteSpuImageError",
			MockImageInfo:           &model.SpuImage{ImageID: 1, SpuID: 1, Url: "http://example.jpg"},
			MockSpuInfo:             &model.Spu{SpuId: 1, CreatorId: 1},
			MockDeleteSpuImageError: errors.New("DeleteSpuImageError"),
			ExpectedError:           errors.New("usecase.DeleteSpuImage failed: DeleteSpuImageError"),
		},
		{
			Name:              "DeleteSpuImageSuccessfully",
			MockSpuInfo:       &model.Spu{SpuId: 1, CreatorId: 1},
			MockImageInfo:     &model.SpuImage{ImageID: 1, SpuID: 1, Url: "http://example.jpg"},
			MockIdentifyError: nil,
			ExpectedError:     nil,
		},
	}

	defer mockey.UnPatchAll()

	var imgId int64 = 1

	for _, tc := range testcase {
		mockey.PatchConvey(tc.Name, t, func() {
			svc := new(service.CommodityService)
			us := &useCase{
				svc: svc,
			}

			mockey.Mock((*service.CommodityService).GetSpuFromImageId).Return(tc.MockSpuInfo, tc.MockImageInfo, tc.MockGetSpuError).Build()
			mockey.Mock((*service.CommodityService).IdentifyUser).Return(tc.MockIdentifyError).Build()
			mockey.Mock((*service.CommodityService).DeleteSpuImage).Return(tc.MockDeleteSpuImageError).Build()

			err := us.DeleteSpuImage(ctx.Background(), imgId)
			if err != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
		})
	}
}

func TestUseCase_ViewSpuImages(t *testing.T) {
	type TestCase struct {
		Name                  string
		MockGetSpuImagesError error
		MockSpuInfo           []*model.SpuImage
		ExpectedError         error
		ExpectedInfo          []*model.SpuImage
		ExpectedTotal         int64
	}

	infos := []*model.SpuImage{
		{
			ImageID: 1,
			SpuID:   1,
			Url:     "http://example.jpg",
		},
		{
			ImageID: 2,
			SpuID:   1,
			Url:     "http://example1.jpg",
		},
	}

	testcase := []TestCase{
		{
			Name:                  "GetSpuImagesError",
			MockSpuInfo:           []*model.SpuImage{},
			MockGetSpuImagesError: errors.New("GetSpuImagesError"),
			ExpectedError:         errors.New("GetSpuImagesError"),
			ExpectedInfo:          []*model.SpuImage{},
			ExpectedTotal:         0,
		},
		{
			Name:          "GetSpuImagesSuccess",
			MockSpuInfo:   infos,
			ExpectedInfo:  infos,
			ExpectedTotal: int64(len(infos)),
		},
	}

	defer mockey.UnPatchAll()

	offset := 0
	limit := 5

	var spuId int64 = 1

	for _, tc := range testcase {
		mockey.PatchConvey(tc.Name, t, func() {
			svc := new(service.CommodityService)
			us := &useCase{
				svc: svc,
			}

			mockey.Mock((*service.CommodityService).GetSpuImages).Return(tc.MockSpuInfo, len(tc.MockSpuInfo), tc.MockGetSpuImagesError).Build()
			infos, total, err := us.ViewSpuImages(ctx.Background(), spuId, offset, limit)
			if err != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
			convey.So(total, convey.ShouldEqual, tc.ExpectedTotal)
			convey.So(infos, convey.ShouldEqual, tc.ExpectedInfo)
		})
	}
}

func TestUseCase_ViewSpu(t *testing.T) {
	type TestCase struct {
		Name                string
		MockSearchItemError error
		MockGetSpuError     error
		MockIds             []int64
		MockSpuInfo         []*model.Spu
		ExpectedError       error
		ExpectedInfo        []*model.Spu
		ExpectedTotal       int64
	}

	ids := []int64{1, 2}

	infos := []*model.Spu{
		{
			SpuId: 1,
			Name:  "OppO phone",
		},
		{
			SpuId: 2,
			Name:  "Vivo phone",
		},
	}

	testcase := []TestCase{
		{
			Name:                "GetSpuError",
			MockSearchItemError: errors.New("SearchItemError"),
			ExpectedError:       errors.New("usecase.ViewSpus failed: SearchItemError"),
			ExpectedInfo:        nil,
			ExpectedTotal:       0,
		},
		{
			Name:                "GetSpuError",
			MockSearchItemError: nil,
			MockIds:             ids,
			MockGetSpuError:     errors.New("GetSpuError"),
			ExpectedInfo:        nil,
			ExpectedError:       errors.New("usecase.ViewSpus failed: GetSpuError"),
		},
		{
			Name:          "GetSpusSuccess",
			MockSpuInfo:   infos,
			MockIds:       ids,
			ExpectedInfo:  infos,
			ExpectedTotal: int64(len(infos)),
		},
	}

	defer mockey.UnPatchAll()

	keyword := "phone"

	for _, tc := range testcase {
		mockey.PatchConvey(tc.Name, t, func() {
			svc := new(service.CommodityService)

			gormDB := new(gorm.DB)
			db := mysql.NewCommodityDB(gormDB)

			ela := new(elastic.Client)
			e := es.NewCommodityElastic(ela)

			us := &useCase{
				svc: svc,
				es:  e,
				db:  db,
			}
			mockey.Mock(mockey.GetMethod(us.es, "SearchItems")).Return(tc.MockIds, len(tc.MockIds), tc.MockSearchItemError).Build()
			mockey.Mock(mockey.GetMethod(us.db, "GetSpuByIds")).Return(tc.MockSpuInfo, tc.MockGetSpuError).Build()

			res, total, err := us.ViewSpus(ctx.Background(), &commodity.ViewSpuReq{
				KeyWord: &keyword,
			})
			if err != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
			convey.So(total, convey.ShouldEqual, tc.ExpectedTotal)
			convey.So(res, convey.ShouldEqual, tc.ExpectedInfo)
		})
	}
}

func TestUseCase_ListSpuImages(t *testing.T) {
	type TestCase struct {
		Name            string
		MockGetSpuError error
		MockSpuInfo     []*model.Spu
		ExpectedSpuInfo []*model.Spu
		ExpectedError   error
	}

	ids := []int64{1, 2}

	infos := []*model.Spu{
		{
			SpuId: 1,
			Name:  "OppO phone",
		},
		{
			SpuId: 2,
			Name:  "Vivo phone",
		},
	}

	testcase := []TestCase{
		{
			Name:            "GetSpuError",
			MockGetSpuError: errors.New("GetSpuError"),
			ExpectedError:   errors.New("GetSpuError"),
			ExpectedSpuInfo: nil,
		},
		{
			Name:            "GetSpuSuccessfully",
			MockGetSpuError: nil,
			MockSpuInfo:     infos,
			ExpectedSpuInfo: infos,
			ExpectedError:   nil,
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testcase {
		mockey.PatchConvey(tc.Name, t, func() {
			svc := new(service.CommodityService)
			gormDB := new(gorm.DB)
			db := mysql.NewCommodityDB(gormDB)

			us := &useCase{
				svc: svc,
				db:  db,
			}

			mockey.Mock(mockey.GetMethod(us.db, "GetSpuByIds")).Return(tc.ExpectedSpuInfo, tc.ExpectedError).Build()

			res, err := us.ListSpuInfo(ctx.Background(), ids)
			if err != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
			convey.So(res, convey.ShouldEqual, tc.ExpectedSpuInfo)
		})
	}
}

func TestUseCase_DecrStock(t *testing.T) {
	type TestCase struct {
		Name                 string
		MockIsHealthy        bool
		MockServiceDecrError error
		MockDBDecrError      error
		ExpectedError        error
	}

	input := []*model.SkuBuyInfo{
		{
			SkuID: 1,
			Count: 1,
		},
		{
			SkuID: 2,
			Count: 2,
		},
	}

	testCases := []TestCase{
		{
			Name:            "decrStockInDBError",
			MockIsHealthy:   false,
			MockDBDecrError: errors.New("ServiceDecrError"),
			ExpectedError:   errors.New("usecase.DecrStock failed: ServiceDecrError"),
		},
		{
			Name:                 "ServiceDecrStockError",
			MockIsHealthy:        true,
			MockServiceDecrError: errors.New("ServiceDecrError"),
			ExpectedError:        errors.New("ServiceDecrError"),
		},
		{
			Name:          "decrStockSuccess",
			ExpectedError: nil,
		},
	}

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			svc := new(service.CommodityService)
			gormDB := new(gorm.DB)
			db := mysql.NewCommodityDB(gormDB)
			redisCache := new(redis.Client)
			cache := redisCommodity.NewCommodityCache(redisCache)

			us := &useCase{
				db:    db,
				cache: cache,
				svc:   svc,
			}
			mockey.Mock(mockey.GetMethod(us.db, "DecrStock")).Return(tc.MockDBDecrError).Build()
			mockey.Mock((*service.CommodityService).DecrStockInNX).Return(tc.MockServiceDecrError).Build()
			mockey.Mock((*service.CommodityService).IsHealthy).Return(tc.MockIsHealthy).Build()
			err := us.DecrStock(ctx.Background(), input)
			if err != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
		})
	}
}

func TestUseCase_DecrLockStock(t *testing.T) {
	type TestCase struct {
		Name                 string
		MockIsHealthy        bool
		MockServiceDecrError error
		MockDBDecrError      error
		ExpectedError        error
	}

	input := []*model.SkuBuyInfo{
		{
			SkuID: 1,
			Count: 1,
		},
		{
			SkuID: 2,
			Count: 2,
		},
	}

	testCases := []TestCase{
		{
			Name:            "decrStockInDBError",
			MockIsHealthy:   false,
			MockDBDecrError: errors.New("ServiceDecrError"),
			ExpectedError:   errors.New("usecase.DecrLockStock failed: ServiceDecrError"),
		},
		{
			Name:                 "ServiceDecrStockError",
			MockIsHealthy:        true,
			MockServiceDecrError: errors.New("ServiceDecrError"),
			ExpectedError:        errors.New("ServiceDecrError"),
		},
		{
			Name:          "decrStockSuccess",
			ExpectedError: nil,
		},
	}

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			svc := new(service.CommodityService)
			gormDB := new(gorm.DB)
			db := mysql.NewCommodityDB(gormDB)
			redisCache := new(redis.Client)
			cache := redisCommodity.NewCommodityCache(redisCache)

			us := &useCase{
				db:    db,
				cache: cache,
				svc:   svc,
			}
			mockey.Mock(mockey.GetMethod(us.db, "DecrLockStock")).Return(tc.MockDBDecrError).Build()
			mockey.Mock((*service.CommodityService).DecrLockStockInNX).Return(tc.MockServiceDecrError).Build()
			mockey.Mock((*service.CommodityService).IsHealthy).Return(tc.MockIsHealthy).Build()
			err := us.DecrLockStock(ctx.Background(), input)
			if err != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
		})
	}
}

func TestUseCase_IncrLockStock(t *testing.T) {
	type TestCase struct {
		Name                 string
		MockIsHealthy        bool
		MockServiceIncrError error
		MockDBIncrError      error
		ExpectedError        error
	}

	input := []*model.SkuBuyInfo{
		{
			SkuID: 1,
			Count: 1,
		},
		{
			SkuID: 2,
			Count: 2,
		},
	}

	testCases := []TestCase{
		{
			Name:            "IncrStockInDBError",
			MockIsHealthy:   false,
			MockDBIncrError: errors.New("ServiceIncrError"),
			ExpectedError:   errors.New("usecase.IncrLockStock failed: ServiceIncrError"),
		},
		{
			Name:                 "ServiceIncrStockError",
			MockIsHealthy:        true,
			MockServiceIncrError: errors.New("ServiceIncrError"),
			ExpectedError:        errors.New("ServiceIncrError"),
		},
		{
			Name:          "IncrStockSuccess",
			ExpectedError: nil,
		},
	}

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			svc := new(service.CommodityService)
			gormDB := new(gorm.DB)
			db := mysql.NewCommodityDB(gormDB)
			redisCache := new(redis.Client)
			cache := redisCommodity.NewCommodityCache(redisCache)

			us := &useCase{
				db:    db,
				cache: cache,
				svc:   svc,
			}
			mockey.Mock(mockey.GetMethod(us.db, "IncrLockStock")).Return(tc.MockDBIncrError).Build()
			mockey.Mock((*service.CommodityService).IncrLockStockInNX).Return(tc.MockServiceIncrError).Build()
			mockey.Mock((*service.CommodityService).IsHealthy).Return(tc.MockIsHealthy).Build()
			err := us.IncrLockStock(ctx.Background(), input)
			if err != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
		})
	}
}
