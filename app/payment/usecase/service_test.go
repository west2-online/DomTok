package usecase

import (
	ctx "context"
	"errors"
	"github.com/bytedance/mockey"
	"github.com/smartystreets/goconvey/convey"
	"testing"

	"github.com/west2-online/DomTok/app/payment/domain/service"
	"github.com/west2-online/DomTok/pkg/constants"
)

func TestUseCase_GetPaymentToken(t *testing.T) {
	type TestCase struct {
		Name                 string
		MockGetUserIDError   error
		MockUid              int64
		MockToken            string
		MockExpirationTime   int64
		MockCheckPaymentErr  error
		MockPaymentExist     bool
		MockGenerateTokenErr error
		MockStoreTokenErr    error
		MockRedisStatus      bool
		ExpectedError        error
		ExpectedToken        string
	}

	testCases := []TestCase{
		{
			Name:               "GetUserIDError",
			MockUid:            1,
			MockGetUserIDError: errors.New("GetUserIDError"),
			ExpectedError:      errors.New("get user id failed:GetUserIDError"),
			ExpectedToken:      "",
		},
		{
			Name:                "PaymentExist",
			MockCheckPaymentErr: nil,
			MockPaymentExist:    constants.PaymentExist,
			ExpectedError:       nil,
			ExpectedToken:       "generatedToken",
		},
		{
			Name:                 "GenerateTokenError",
			MockToken:            "generatedToken",
			MockExpirationTime:   int64(3600),
			MockGenerateTokenErr: errors.New("GenerateTokenError"),
			ExpectedError:        errors.New("generate payment token failed:GenerateTokenError"),
			ExpectedToken:        "",
		},
		{
			Name:              "StoreTokenError",
			MockStoreTokenErr: errors.New("StoreTokenError"),
			MockRedisStatus:   false,
			ExpectedError:     errors.New("store payment token failed:StoreTokenError"),
			ExpectedToken:     "",
		},
		{
			Name:            "GeneratePaymentTokenSuccessfully",
			MockRedisStatus: true,
			ExpectedError:   nil,
			ExpectedToken:   "generatedToken",
		},
	}

	/*p:=&model.PaymentOrder{
		OrderID:123,
	}*/
	orderID := int64(123)

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			mockey.Mock((*service.PaymentService).GetUserID).Return(tc.MockUid, tc.MockGetUserIDError).Build()
			mockey.Mock((*service.PaymentService).GeneratePaymentToken).Return(tc.MockToken, tc.MockExpirationTime, tc.MockGenerateTokenErr).Build()
			mockey.Mock((*service.PaymentService).StorePaymentToken).Return(tc.MockRedisStatus, tc.MockStoreTokenErr).Build()
			us := new(paymentUseCase)
			paymentSvc := new(service.PaymentService)

			token, _, err := us.GetPaymentToken(ctx.Background(), orderID)
			if tc.ExpectedError == nil {
				convey.So(err, convey.ShouldBeNil)
				convey.So(token, convey.ShouldEqual, tc.ExpectedToken)
			} else {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			}
		})
	}
}

/*func TestUseCase_CreateRefund(t *testing.T) {
	type TestCase struct {
		Name                  string
		MockGetUserIDError    error
		MockRedisRateLimitErr error
		MockCreateRefundErr   error
		ExpectedError         error
		ExpectedRefundID      int64
	}

	testCases := []TestCase{
		{
			Name:               "GetUserIDError",
			MockGetUserIDError: errors.New("GetUserIDError"),
			ExpectedError:      errors.New("get user id failed: GetUserIDError"),
			ExpectedRefundID:   0,
		},
		{
			Name:                "CreateRefundError",
			MockCreateRefundErr: errors.New("CreateRefundError"),
			ExpectedError:       errors.New("create refund info failed: CreateRefundError"),
			ExpectedRefundID:    0,
		},
		{
			Name:             "CreateRefundSuccessfully",
			ExpectedError:    nil,
			ExpectedRefundID: 1,
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			paymentSvc := new(service.PaymentService)
			us := &paymentUseCase{
				svc: paymentSvc,
			}

			mockey.Mock((*service.PaymentService).GetUserID).Return(int64(1), tc.MockGetUserIDError).Build()
			mockey.Mock((*service.PaymentService).CreateRefundInfo).Return(tc.ExpectedRefundID, tc.MockCreateRefundErr).Build()

			_, refundID, err := us.CreateRefund(ctx.Background(), 1)
			if tc.ExpectedError == nil {
				convey.So(err, convey.ShouldBeNil)
				convey.So(refundID, convey.ShouldEqual, tc.ExpectedRefundID)
			} else {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			}
		})
	}
}*/
