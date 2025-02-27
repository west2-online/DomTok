package usecase

import (
	ctx "context"
	"errors"
	"github.com/west2-online/DomTok/app/payment/infrastructure/mysql"
	"gorm.io/gorm"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/payment/domain/model"
	"github.com/west2-online/DomTok/app/payment/domain/service"
	paymentStatus "github.com/west2-online/DomTok/pkg/constants"
)

func TestPaymentUseCase_GetPaymentToken(t *testing.T) {
	type TestCase struct {
		Name                       string
		MockGetUserIDError         error
		MockUserID                 int64
		MockCheckPaymentExistError error
		MockPaymentExistStatus     bool
		MockCreatePaymentInfoError error
		MockPaymentInfo            *model.PaymentOrder
		MockGetPaymentInfoError    error
		MockGenTokenError          error
		MockToken                  string
		MockExpTime                int64
		MockStoreTokenError        error
		MockRedisStoreStatus       bool
		ExpectedToken              string
		ExpectedExpTime            int64
		ExpectedError              error
	}

	testCases := []TestCase{
		{
			Name:               "GetUserIDError",
			MockGetUserIDError: errors.New("GetUserIDError"),
			ExpectedToken:      "",
			ExpectedExpTime:    0,
			ExpectedError:      errors.New("get user id failed:GetUserIDError"),
		},
		{
			Name:                       "CheckPaymentExistError",
			MockGetUserIDError:         nil,
			MockUserID:                 1,
			MockCheckPaymentExistError: errors.New("CheckPaymentExistError"),
			ExpectedToken:              "",
			ExpectedExpTime:            0,
			ExpectedError:              errors.New("check payment existed failed:CheckPaymentExistError"),
		},
		{
			Name:                       "CreatePaymentInfoError",
			MockGetUserIDError:         nil,
			MockUserID:                 1,
			MockCheckPaymentExistError: nil,
			MockPaymentExistStatus:     paymentStatus.PaymentNotExist,
			MockCreatePaymentInfoError: errors.New("CreatePaymentInfoError"),
			ExpectedToken:              "",
			ExpectedExpTime:            0,
			ExpectedError:              errors.New("create payment info failed:CreatePaymentInfoError"),
		},
		{
			Name:                       "GetPaymentInfoError",
			MockGetUserIDError:         nil,
			MockUserID:                 1,
			MockCheckPaymentExistError: nil,
			MockPaymentExistStatus:     paymentStatus.PaymentExist,
			MockGetPaymentInfoError:    errors.New("GetPaymentInfoError"),
			ExpectedToken:              "",
			ExpectedExpTime:            0,
			ExpectedError:              errors.New("get payment info failed:GetPaymentInfoError"),
		},
		{
			Name:                       "PaymentAlreadyProcessing",
			MockGetUserIDError:         nil,
			MockUserID:                 1,
			MockCheckPaymentExistError: nil,
			MockPaymentExistStatus:     paymentStatus.PaymentExist,
			MockGetPaymentInfoError:    nil,
			MockPaymentInfo: &model.PaymentOrder{
				Status: paymentStatus.PaymentStatusProcessingCode,
			},
			ExpectedToken:   "",
			ExpectedExpTime: 0,
			ExpectedError:   errors.New("payment is processing or has already done:%!w(<nil>)"),
		},
		{
			Name:                       "GenerateTokenError",
			MockGetUserIDError:         nil,
			MockUserID:                 1,
			MockCheckPaymentExistError: nil,
			MockPaymentExistStatus:     paymentStatus.PaymentExist,
			MockGetPaymentInfoError:    nil,
			MockPaymentInfo: &model.PaymentOrder{
				Status: paymentStatus.PaymentStatusFailedCode,
			},
			MockGenTokenError: errors.New("GenerateTokenError"),
			ExpectedToken:     "",
			ExpectedExpTime:   0,
			ExpectedError:     errors.New("generate payment token failed:GenerateTokenError"),
		},
		{
			Name:                       "StoreTokenError",
			MockGetUserIDError:         nil,
			MockUserID:                 1,
			MockCheckPaymentExistError: nil,
			MockPaymentExistStatus:     paymentStatus.PaymentExist,
			MockGetPaymentInfoError:    nil,
			MockPaymentInfo: &model.PaymentOrder{
				Status: paymentStatus.PaymentStatusFailedCode,
			},
			MockGenTokenError:    nil,
			MockToken:            "test-token",
			MockExpTime:          1000,
			MockStoreTokenError:  errors.New("StoreTokenError"),
			MockRedisStoreStatus: false,
			ExpectedToken:        "",
			ExpectedExpTime:      0,
			ExpectedError:        errors.New("store payment token failed:StoreTokenError"),
		},
		{
			Name:                       "SuccessfulGetPaymentToken",
			MockGetUserIDError:         nil,
			MockUserID:                 1,
			MockCheckPaymentExistError: nil,
			MockPaymentExistStatus:     paymentStatus.PaymentExist,
			MockGetPaymentInfoError:    nil,
			MockPaymentInfo: &model.PaymentOrder{
				Status: paymentStatus.PaymentStatusFailedCode,
			},
			MockGenTokenError:    nil,
			MockToken:            "test-token",
			MockExpTime:          1000,
			MockStoreTokenError:  nil,
			MockRedisStoreStatus: paymentStatus.RedisStoreSuccess,
			ExpectedToken:        "test-token",
			ExpectedExpTime:      1000,
			ExpectedError:        nil,
		},
	}

	defer mockey.UnPatchAll()
	var orderID int64 = 1

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			svc := new(service.PaymentService)
			gormDB := new(gorm.DB)
			db := mysql.NewPaymentDB(gormDB)
			uc := &paymentUseCase{
				svc: svc,
				db:  db,
			}

			mockey.Mock((*service.PaymentService).GetUserID).Return(tc.MockUserID, tc.MockGetUserIDError).Build()
			mockey.Mock(mockey.GetMethod(uc.db, "CheckPaymentExist")).Return(tc.MockPaymentExistStatus, tc.MockCheckPaymentExistError).Build()
			mockey.Mock((*service.PaymentService).CreatePaymentInfo).Return(int64(1), tc.MockCreatePaymentInfoError).Build()
			mockey.Mock(mockey.GetMethod(uc.db, "GetPaymentInfo")).Return(tc.MockPaymentInfo, tc.MockGetPaymentInfoError).Build()
			mockey.Mock((*service.PaymentService).GeneratePaymentToken).Return(tc.MockToken, tc.MockExpTime, tc.MockGenTokenError).Build()
			mockey.Mock((*service.PaymentService).StorePaymentToken).Return(tc.MockRedisStoreStatus, tc.MockStoreTokenError).Build()

			token, expTime, err := uc.GetPaymentToken(ctx.Background(), orderID)
			if err != nil && tc.ExpectedError != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
			convey.So(token, convey.ShouldEqual, tc.ExpectedToken)
			convey.So(expTime, convey.ShouldEqual, tc.ExpectedExpTime)
		})
	}
}

func TestPaymentUseCase_CreateRefund(t *testing.T) {
	type TestCase struct {
		Name                      string
		MockCheckOrderExistError  error
		MockOrderExists           bool
		MockGetUserIDError        error
		MockUserID                int64
		MockRateLimitError        error
		MockFrequencyValid        bool
		MockTimeValid             bool
		MockCreateRefundInfoError error
		MockRefundID              int64
		ExpectedRefundStatus      int64
		ExpectedRefundID          int64
		ExpectedError             error
	}

	testCases := []TestCase{
		{
			Name:                     "CheckOrderExistError",
			MockCheckOrderExistError: errors.New("CheckOrderExistError"),
			// TODO 这个错误信息后面可能要改动，因为今晚跑太多次了
			MockGetUserIDError: nil,
			MockUserID:         1,
			MockRateLimitError: nil,
			// TODO这个可能要改
			MockFrequencyValid:   false,
			MockTimeValid:        true,
			ExpectedRefundStatus: 0,
			ExpectedRefundID:     0,
			ExpectedError:        errors.New("too many refund requests in a short time"),
		},
		{
			Name:                     "OrderNotExist",
			MockCheckOrderExistError: nil,
			MockOrderExists:          false,
			// TODO 这个错误信息后面可能要改动，因为今晚跑太多次了
			MockGetUserIDError:   nil,
			MockUserID:           1,
			MockRateLimitError:   nil,
			MockFrequencyValid:   false,
			MockTimeValid:        true,
			ExpectedRefundStatus: 0,
			ExpectedRefundID:     0,
			ExpectedError:        errors.New("too many refund requests in a short time"),
		},
		{
			Name:                     "GetUserIDError",
			MockCheckOrderExistError: nil,
			MockOrderExists:          true,
			MockGetUserIDError:       errors.New("GetUserIDError"),
			ExpectedRefundStatus:     0,
			ExpectedRefundID:         0,
			ExpectedError:            errors.New("get user id failed: GetUserIDError"),
		},
		{
			Name:                     "CheckRateLimitError",
			MockCheckOrderExistError: nil,
			MockOrderExists:          true,
			MockGetUserIDError:       nil,
			MockUserID:               1,
			MockRateLimitError:       errors.New("RateLimitError"),
			ExpectedRefundStatus:     0,
			ExpectedRefundID:         0,
			ExpectedError:            errors.New("check redis rate limiting failed: RateLimitError"),
		},
		{
			Name:                     "FrequencyInvalid",
			MockCheckOrderExistError: nil,
			MockOrderExists:          true,
			MockGetUserIDError:       nil,
			MockUserID:               1,
			MockRateLimitError:       nil,
			MockFrequencyValid:       false,
			ExpectedRefundStatus:     0,
			ExpectedRefundID:         0,
			ExpectedError:            errors.New("too many refund requests in a short time"),
		},
		{
			Name:                     "TimeInvalid",
			MockCheckOrderExistError: nil,
			MockOrderExists:          true,
			MockGetUserIDError:       nil,
			MockUserID:               1,
			MockRateLimitError:       nil,
			MockFrequencyValid:       true,
			MockTimeValid:            false,
			ExpectedRefundStatus:     0,
			ExpectedRefundID:         0,
			ExpectedError:            errors.New("refund already requested for this order in the last 24 hours"),
		},
		{
			Name:                      "CreateRefundInfoError",
			MockCheckOrderExistError:  nil,
			MockOrderExists:           true,
			MockGetUserIDError:        nil,
			MockUserID:                1,
			MockRateLimitError:        nil,
			MockFrequencyValid:        true,
			MockTimeValid:             true,
			MockCreateRefundInfoError: errors.New("CreateRefundInfoError"),
			ExpectedRefundStatus:      0,
			ExpectedRefundID:          0,
			ExpectedError:             errors.New("create refund info failed: CreateRefundInfoError"),
		},
		{
			Name:                      "SuccessfulCreateRefund",
			MockCheckOrderExistError:  nil,
			MockOrderExists:           true,
			MockGetUserIDError:        nil,
			MockUserID:                1,
			MockRateLimitError:        nil,
			MockFrequencyValid:        true,
			MockTimeValid:             true,
			MockCreateRefundInfoError: nil,
			MockRefundID:              123,
			ExpectedRefundStatus:      paymentStatus.RefundStatusProcessingCode,
			ExpectedRefundID:          123,
			ExpectedError:             nil,
		},
	}

	defer mockey.UnPatchAll()
	var orderID int64 = 1

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			svc := new(service.PaymentService)
			uc := &paymentUseCase{
				svc: svc,
			}
			mockey.Mock((*service.PaymentService).CheckOrderExist).Return(tc.MockOrderExists, tc.MockCheckOrderExistError).Build()
			mockey.Mock((*service.PaymentService).GetUserID).Return(tc.MockUserID, tc.MockGetUserIDError).Build()
			mockey.Mock((*service.PaymentService).CheckRedisRateLimiting).Return(tc.MockFrequencyValid, tc.MockTimeValid, tc.MockRateLimitError).Build()
			mockey.Mock((*service.PaymentService).CreateRefundInfo).Return(tc.MockRefundID, tc.MockCreateRefundInfoError).Build()

			status, refundID, err := uc.CreateRefund(ctx.Background(), orderID)
			if err != nil && tc.ExpectedError != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
			convey.So(status, convey.ShouldEqual, tc.ExpectedRefundStatus)
			convey.So(refundID, convey.ShouldEqual, tc.ExpectedRefundID)
		})
	}
}
