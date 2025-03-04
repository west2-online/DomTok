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
	"time"

	"github.com/bytedance/mockey"
	"github.com/smartystreets/goconvey/convey"
	"gorm.io/gorm"

	"github.com/west2-online/DomTok/app/payment/domain/model"
	"github.com/west2-online/DomTok/app/payment/domain/repository"
	"github.com/west2-online/DomTok/app/payment/domain/service"
	"github.com/west2-online/DomTok/app/payment/infrastructure/mysql"
	paymentStatus "github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
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
			ExpectedError:   errors.New("[4011] payment is processing or has already done"),
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
			mockey.Mock((*service.PaymentService).CheckOrderExist).Return(paymentStatus.PaymentExist, nil).Build()

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
			MockGetUserIDError:       nil,
			MockUserID:               1,
			MockRateLimitError:       nil,
			MockFrequencyValid:       false,
			MockTimeValid:            true,
			ExpectedRefundStatus:     0,
			ExpectedRefundID:         0,
			ExpectedError:            errors.New("check order existence failed: CheckOrderExistError"),
		},
		{
			Name:                     "OrderNotExist",
			MockCheckOrderExistError: nil,
			MockOrderExists:          false,
			MockGetUserIDError:       nil,
			MockUserID:               1,
			MockRateLimitError:       nil,
			MockFrequencyValid:       false,
			MockTimeValid:            true,
			ExpectedRefundStatus:     0,
			ExpectedRefundID:         0,
			ExpectedError:            errors.New("[4000] order does not exist"),
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

func TestPaymentUseCase_RefundReview(t *testing.T) {
	type _DB struct {
		repository.PaymentDB
	}
	mockey.Mock((*service.PaymentService).GetOrderStatus).Return(true, false, nil).Build()
	mockey.Mock((*service.PaymentService).GetUserID).Return(int64(1), nil).Build()
	mockey.Mock((*service.PaymentService).CheckAdminPermission).Return(true, nil).Build()
	mockey.Mock((*service.PaymentService).Refund).Return(int64(1), "test", nil).Build()
	mockey.Mock((*service.PaymentService).CancelOrder).Return(nil).Build()
	mockey.Mock((*_DB).GetRefundInfoByOrderID).Return(&model.PaymentRefund{}, nil).Build()
	mockey.Mock((*_DB).UpdateRefundStatusByOrderIDAndStatus).Return(nil).Build()
	mockey.Mock((*_DB).UpdateRefundStatusToSuccessAndCreateLedgerAsTransaction).Return(nil).Build()

	defer mockey.UnPatchAll()
	uc := &paymentUseCase{
		db:  &_DB{},
		svc: new(service.PaymentService),
	}
	bg := ctx.Background()
	orderID := int64(1)
	testErr := errno.NewErrNo(-1, "")

	mockey.PatchConvey("RefundReview", t, func() {
		mockey.PatchConvey("success", func() {
			err := uc.RefundReview(bg, orderID, true)
			convey.So(err, convey.ShouldBeNil)
			err = uc.RefundReview(bg, orderID, false)
			convey.So(err, convey.ShouldBeNil)
		})
		mockey.PatchConvey("GetOrderStatusError", func() {
			mockey.Mock((*service.PaymentService).GetOrderStatus).Return(false, false, testErr).Build()
			err := uc.RefundReview(bg, orderID, true)
			convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
		})
		mockey.PatchConvey("GetOrderStatusNotExist", func() {
			mockey.Mock((*service.PaymentService).GetOrderStatus).Return(false, false, nil).Build()
			err := uc.RefundReview(bg, orderID, true)
			convey.So(errno.ConvertErr(err).ErrorCode, convey.ShouldEqual, errno.ServiceOrderNotFound)
		})
		mockey.PatchConvey("GetOrderStatusExpired", func() {
			mockey.Mock((*service.PaymentService).GetOrderStatus).Return(true, true, nil).Build()
			err := uc.RefundReview(bg, orderID, true)
			convey.So(errno.ConvertErr(err).ErrorCode, convey.ShouldEqual, errno.ServiceOrderExpired)
		})
		mockey.PatchConvey("GetUserIDError", func() {
			mockey.Mock((*service.PaymentService).GetUserID).Return(int64(0), testErr).Build()
			err := uc.RefundReview(bg, orderID, true)
			convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
		})
		mockey.PatchConvey("CheckAdminPermissionError", func() {
			mockey.Mock((*service.PaymentService).CheckAdminPermission).Return(true, testErr).Build()
			err := uc.RefundReview(bg, orderID, true)
			convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
		})
		mockey.PatchConvey("CheckAdminPermissionNoPermission", func() {
			mockey.Mock((*service.PaymentService).CheckAdminPermission).Return(false, nil).Build()
			err := uc.RefundReview(bg, orderID, true)
			convey.So(errno.ConvertErr(err).ErrorCode, convey.ShouldEqual, errno.AuthNoOperatePermission.ErrorCode)
		})
		mockey.PatchConvey("GetRefundInfoByOrderIDError", func() {
			mockey.Mock((*_DB).GetRefundInfoByOrderID).Return(&model.PaymentRefund{}, testErr).Build()
			err := uc.RefundReview(bg, orderID, true)
			convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
		})
		mockey.PatchConvey("GetRefundInfoByOrderIDNotExist", func() {
			mockey.Mock((*_DB).GetRefundInfoByOrderID).Return(nil, nil).Build()
			err := uc.RefundReview(bg, orderID, true)
			convey.So(errno.ConvertErr(err).ErrorCode, convey.ShouldEqual, errno.ServicePaymentRefundNotExist)
		})
		mockey.PatchConvey("UpdateRefundStatusByOrderIDAndStatusError", func() {
			mockey.Mock((*_DB).UpdateRefundStatusByOrderIDAndStatus).Return(testErr).Build()
			err := uc.RefundReview(bg, orderID, true)
			convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
		})
		mockey.PatchConvey("UpdateRefundStatusByOrderIDAndStatusErrorInPassEqualFalse", func() {
			mockey.Mock((*_DB).UpdateRefundStatusByOrderIDAndStatus).Return(mockey.Sequence(nil).Then(testErr)).Build()
			err := uc.RefundReview(bg, orderID, false)
			convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
		})
		mockey.PatchConvey("RefundError", func() {
			mockey.Mock((*service.PaymentService).Refund).Return(int64(0), "", testErr).Build()
			err := uc.RefundReview(bg, orderID, true)
			convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
		})
		mockey.PatchConvey("CancelOrderError", func() {
			mockey.Mock((*service.PaymentService).CancelOrder).Return(testErr).Build()
			err := uc.RefundReview(bg, orderID, true)
			convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
		})
		mockey.PatchConvey("UpdateRefundStatusToSuccessAndCreateLedgerAsTransactionError", func() {
			mockey.Mock((*_DB).UpdateRefundStatusToSuccessAndCreateLedgerAsTransaction).Return(testErr).Build()
			err := uc.RefundReview(bg, orderID, true)
			convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
		})
	})
}

func TestPaymentUseCase_PaymentCheckout(t *testing.T) {
	type _DB struct {
		repository.PaymentDB
	}
	mockey.Mock((*service.PaymentService).GetOrderStatus).Return(true, false, nil).Build()
	mockey.Mock((*service.PaymentService).GetUserID).Return(int64(1), nil).Build()
	mockey.Mock((*service.PaymentService).GetExpiredAtAndDelPaymentToken).Return(true, time.Time{}, nil).Build()
	mockey.Mock((*service.PaymentService).GetPayInfo).Return(0, "", nil).Build()
	mockey.Mock((*service.PaymentService).ConfirmOrder).Return(nil).Build()
	mockey.Mock((*service.PaymentService).Pay).Return(0, "", nil).Build()
	mockey.Mock((*service.PaymentService).PutBackPaymentToken).Return(nil).Build()
	mockey.Mock((*_DB).GetPaymentInfo).Return(&model.PaymentOrder{Status: paymentStatus.PaymentStatusPendingCode}, nil).Build()
	mockey.Mock((*_DB).UpdatePaymentStatus).Return(nil).Build()
	mockey.Mock((*_DB).UpdatePaymentStatusToSuccessAndCreateLedgerAsTransaction).Return(nil).Build()
	uc := &paymentUseCase{
		db:  &_DB{},
		svc: new(service.PaymentService),
	}
	bg := ctx.Background()
	orderID := int64(1)
	testErr := errno.NewErrNo(-1, "")
	mockey.PatchConvey("PaymentCheckout", t, func() {
		mockey.PatchConvey("success", func() {
			err := uc.PaymentCheckout(bg, orderID, "test")
			convey.So(err, convey.ShouldBeNil)
		})
		mockey.PatchConvey("GetOrderStatusError", func() {
			mockey.Mock((*service.PaymentService).GetOrderStatus).Return(true, false, testErr).Build()
			err := uc.PaymentCheckout(bg, orderID, "test")
			convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
		})
		mockey.PatchConvey("GetOrderStatusNotExist", func() {
			mockey.Mock((*service.PaymentService).GetOrderStatus).Return(false, false, nil).Build()
			err := uc.PaymentCheckout(bg, orderID, "test")
			convey.So(errno.ConvertErr(err).ErrorCode, convey.ShouldEqual, errno.ServiceOrderNotFound)
		})
		mockey.PatchConvey("GetOrderStatusExpired", func() {
			mockey.Mock((*service.PaymentService).GetOrderStatus).Return(true, true, nil).Build()
			err := uc.PaymentCheckout(bg, orderID, "test")
			convey.So(errno.ConvertErr(err).ErrorCode, convey.ShouldEqual, errno.ServiceOrderExpired)
		})
		mockey.PatchConvey("GetUserIDError", func() {
			mockey.Mock((*service.PaymentService).GetUserID).Return(int64(0), testErr).Build()
			err := uc.PaymentCheckout(bg, orderID, "test")
			convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
		})
		mockey.PatchConvey("GetExpiredAtAndDelPaymentTokenErrorSoSecondCheck", func() {
			mockey.Mock((*service.PaymentService).GetExpiredAtAndDelPaymentToken).Return(false, time.Time{}, testErr).Build()
			err := uc.PaymentCheckout(bg, orderID, "test")
			convey.So(err, convey.ShouldEqual, nil)

			mockey.PatchConvey("GetPayInfoError", func() {
				mockey.Mock((*_DB).GetPaymentInfo).
					Return(&model.PaymentOrder{Status: paymentStatus.PaymentStatusPendingCode}, testErr).Build()
				err := uc.PaymentCheckout(bg, orderID, "test")
				convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
			})
			mockey.PatchConvey("PaymentAlreadyProcessing", func() {
				mockey.Mock((*_DB).GetPaymentInfo).
					Return(&model.PaymentOrder{Status: paymentStatus.PaymentStatusProcessingCode}, nil).Build()
				err := uc.PaymentCheckout(bg, orderID, "test")
				convey.So(errno.ConvertErr(err).ErrorCode, convey.ShouldEqual, errno.IllegalOperatorCode)
			})
			mockey.PatchConvey("UpdatePaymentStatusError", func() {
				mockey.Mock((*_DB).GetPaymentInfo).
					Return(&model.PaymentOrder{Status: paymentStatus.PaymentStatusPendingCode}, nil).Build()
				mockey.Mock((*_DB).UpdatePaymentStatus).Return(testErr).Build()
				err := uc.PaymentCheckout(bg, orderID, "test")
				convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
			})
		})
		mockey.PatchConvey("GetExpiredAtAndDelPaymentTokenNotExist", func() {
			mockey.Mock((*service.PaymentService).GetExpiredAtAndDelPaymentToken).Return(false, time.Time{}, nil).Build()
			err := uc.PaymentCheckout(bg, orderID, "test")
			convey.So(errno.ConvertErr(err).ErrorCode, convey.ShouldEqual, errno.IllegalOperatorCode)
		})
		mockey.PatchConvey("GetPayInfoError", func() {
			mockey.Mock((*service.PaymentService).GetPayInfo).Return(0, "", testErr).Build()
			err := uc.PaymentCheckout(bg, orderID, "test")
			convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
		})
		mockey.PatchConvey("ConfirmOrderError", func() {
			mockey.Mock((*service.PaymentService).ConfirmOrder).Return(testErr).Build()
			err := uc.PaymentCheckout(bg, orderID, "test")
			convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)

			mockey.PatchConvey("PutBackPaymentTokenError", func() {
				mockey.Mock((*service.PaymentService).PutBackPaymentToken).Return(testErr).Build()
				err := uc.PaymentCheckout(bg, orderID, "test")
				convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
			})
			mockey.PatchConvey("UpdatePaymentStatusError", func() {
				mockey.Mock((*service.PaymentService).GetExpiredAtAndDelPaymentToken).
					Return(true, time.Time{}, errno.NewErrNo(-2, "")).Build()
				err := uc.PaymentCheckout(bg, orderID, "test")
				convey.So(errno.ConvertErr(err).ErrorCode, convey.ShouldEqual, -1)
				mockey.Mock((*_DB).UpdatePaymentStatus).Return(mockey.Sequence(nil).Then(testErr)).Build()
				err = uc.PaymentCheckout(bg, orderID, "test")
				convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
			})
		})
		mockey.PatchConvey("PayError", func() {
			mockey.Mock((*service.PaymentService).Pay).Return(0, "", testErr).Build()
			err := uc.PaymentCheckout(bg, orderID, "test")
			convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)

			mockey.PatchConvey("UpdatePaymentStatusError", func() {
				mockey.Mock((*_DB).UpdatePaymentStatus).Return(testErr).Build()
				err := uc.PaymentCheckout(bg, orderID, "test")
				convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
			})
		})
		mockey.PatchConvey("RedisSuccessAndGetPayInfoError", func() {
			mockey.Mock((*_DB).GetPaymentInfo).Return(&model.PaymentOrder{Status: paymentStatus.PaymentStatusPendingCode}, testErr).Build()
			err := uc.PaymentCheckout(bg, orderID, "test")
			convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)

			mockey.PatchConvey("UpdatePaymentStatusToSuccessAndCreateLedgerAsTransaction", func() {
				mockey.Mock((*_DB).GetPaymentInfo).Return(&model.PaymentOrder{Status: paymentStatus.PaymentStatusPendingCode}, nil).Build()
				mockey.Mock((*_DB).UpdatePaymentStatusToSuccessAndCreateLedgerAsTransaction).Return(testErr).Build()
				err := uc.PaymentCheckout(bg, orderID, "test")
				convey.So(errno.ConvertErr(err), convey.ShouldEqual, testErr)
			})
		})
	})
}
