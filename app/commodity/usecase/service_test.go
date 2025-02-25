package usecase

import (
	ctx "context"
	"errors"
	"github.com/bytedance/mockey"
	"github.com/smartystreets/goconvey/convey"
	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/domain/service"
	"github.com/west2-online/DomTok/pkg/base/context"
	"testing"
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
