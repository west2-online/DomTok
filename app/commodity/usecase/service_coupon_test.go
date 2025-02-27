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
	"fmt"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/smartystreets/goconvey/convey"
	"gorm.io/gorm"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/domain/service"
	"github.com/west2-online/DomTok/app/commodity/infrastructure/mysql"
	contextLogin "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/errno"
)

func TestUseCase_CreateCoupon(t *testing.T) {
	type TestCase struct {
		Name             string
		MockVerifyError  error
		MockInitError    error
		MockCreateError  error
		ExpectedError    error
		ExpectedCouponId int64
	}

	testCases := []TestCase{
		{
			Name:             "VerifyError",
			MockVerifyError:  errors.New("verify error"),
			ExpectedError:    errors.New("verify error"),
			ExpectedCouponId: -1,
		},
		{
			Name:             "InitCouponError",
			MockVerifyError:  nil,
			MockInitError:    errors.New("init error"),
			ExpectedError:    fmt.Errorf("usecase.CreateCoupon error: init error"),
			ExpectedCouponId: -1,
		},
		{
			Name:             "CreateCouponDBError",
			MockVerifyError:  nil,
			MockInitError:    nil,
			MockCreateError:  errors.New("db error"),
			ExpectedError:    fmt.Errorf("usecase.CreateCoupon error: db error"),
			ExpectedCouponId: -1,
		},
		{
			Name:             "CreateCouponSuccess",
			MockVerifyError:  nil,
			MockInitError:    nil,
			MockCreateError:  nil,
			ExpectedError:    nil,
			ExpectedCouponId: 1001,
		},
	}

	coupon := &model.Coupon{
		Id:             0,
		Uid:            101,
		Name:           "Test Coupon",
		TypeInfo:       1,
		ConditionCost:  99.99,
		DiscountAmount: 20.0,
		Discount:       0.0,
		RangeType:      0,
		RangeId:        0,
		Description:    "Testing coupon creation",
		ExpireTime:     time.Now().Add(24 * time.Hour),
		DeadlineForGet: time.Now().Add(12 * time.Hour),
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			gormDB := new(gorm.DB)
			db := mysql.NewCommodityDB(gormDB)
			uc := &useCase{
				svc: new(service.CommodityService),
				db:  db,
			}
			// Mock the svc.Verify and svc.InitCoupon calls
			mockey.Mock((*service.CommodityService).Verify).Return(tc.MockVerifyError).Build()
			mockey.Mock((*service.CommodityService).InitCoupon).Return(tc.MockInitError).Build()

			// Mock the db.CreateCoupon call
			mockey.Mock(mockey.GetMethod(uc.db, "CreateCoupon")).Return(tc.ExpectedCouponId, tc.MockCreateError).Build()

			couponId, err := uc.CreateCoupon(ctx.Background(), coupon)

			if err != nil && tc.ExpectedError != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				// Either no error was expected or none was returned
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}

			convey.So(couponId, convey.ShouldEqual, tc.ExpectedCouponId)
		})
	}
}

func TestUseCase_DeleteCoupon(t *testing.T) {
	// Define test cases.
	type TestCase struct {
		Name           string
		MockLoginUID   int64
		MockLoginError error
		// For GetCouponById
		CouponFound    bool
		CouponInDB     *model.Coupon
		GetCouponError error
		// For DeleteCouponById
		MockDeleteError error
		// Expected result
		ExpectedError error
	}

	// For demonstration, if errno.ParamVerifyError and errno.AuthInvalid are not defined,
	// you could uncomment the following definitions:
	// var paramVerifyError = errors.New("param verify error")
	// var authInvalidError = errors.New("auth invalid error")

	testCases := []TestCase{
		{
			Name:           "LoginError",
			MockLoginError: errors.New("login error"),
			ExpectedError:  fmt.Errorf("usecase.DeleteCoupon get logindata error: %w", errors.New("login error")),
		},
		{
			Name:           "GetCouponByIdError",
			MockLoginUID:   101,
			MockLoginError: nil,
			// Simulate an error when fetching the coupon.
			GetCouponError: errors.New("get coupon error"),
			ExpectedError:  fmt.Errorf("usecase.DeleteCoupon error: %w", errors.New("get coupon error")),
		},
		{
			Name:           "CouponNotFound",
			MockLoginUID:   101,
			MockLoginError: nil,
			// Coupon not found: e == false.
			CouponFound:    false,
			CouponInDB:     nil,
			GetCouponError: nil,
			ExpectedError:  errno.ParamVerifyError,
			// or, if not using errno, use: errors.New("param verify error")
		},
		{
			Name:           "UIDMismatch",
			MockLoginUID:   101,
			MockLoginError: nil,
			CouponFound:    true,
			// Coupon exists but its Uid does not match the logged in user.
			CouponInDB: &model.Coupon{
				Id:  1,
				Uid: 202,
			},
			GetCouponError: nil,
			ExpectedError:  errno.AuthInvalid, // or errors.New("auth invalid")
		},
		{
			Name:           "DeleteCouponByIdError",
			MockLoginUID:   101,
			MockLoginError: nil,
			CouponFound:    true,
			CouponInDB: &model.Coupon{
				Id:  1,
				Uid: 101,
			},
			GetCouponError:  nil,
			MockDeleteError: errors.New("delete error"),
			ExpectedError:   fmt.Errorf("usecase.DeleteCoupon error: %w", errors.New("delete error")),
		},
		{
			Name:           "DeleteCouponSuccess",
			MockLoginUID:   101,
			MockLoginError: nil,
			CouponFound:    true,
			CouponInDB: &model.Coupon{
				Id:  1,
				Uid: 101,
			},
			GetCouponError:  nil,
			MockDeleteError: nil,
			ExpectedError:   nil,
		},
	}

	// The coupon to delete. Its Id should match the one we expect from GetCouponById.
	coupon := &model.Coupon{
		Id:  1,
		Uid: 101,
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// Patch the login data retrieval.
			mockey.Mock(contextLogin.GetLoginData).Return(tc.MockLoginUID, tc.MockLoginError).Build()

			// Create a dummy db instance.
			gormDB := new(gorm.DB)
			db := mysql.NewCommodityDB(gormDB)
			// Patch GetCouponById method: returns found flag, coupon details, and error.
			mockey.Mock(mockey.GetMethod(db, "GetCouponById")).Return(tc.CouponFound, tc.CouponInDB, tc.GetCouponError).Build()

			// Patch DeleteCouponById method.
			mockey.Mock(mockey.GetMethod(db, "DeleteCouponById")).Return(tc.MockDeleteError).Build()

			uc := &useCase{
				svc: new(service.CommodityService), // not used in DeleteCoupon
				db:  db,
			}

			err := uc.DeleteCoupon(ctx.Background(), coupon)
			if err != nil && tc.ExpectedError != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
		})
	}
}

// TestUseCase_GetCreatorCoupons tests the GetCreatorCoupons method.
func TestUseCase_GetCreatorCoupons(t *testing.T) {
	type TestCase struct {
		Name            string
		PageNum         int64
		MockVerifyError error
		MockLoginUID    int64
		MockLoginError  error
		// For GetCouponsByCreatorId:
		CouponsFromDB   []*model.Coupon
		GetCouponsError error
		// Expected result:
		ExpectedError   error
		ExpectedCoupons []*model.Coupon
	}

	testCases := []TestCase{
		{
			Name:            "VerifyError",
			PageNum:         1,
			MockVerifyError: errors.New("invalid page number"),
			ExpectedError:   errors.New("invalid page number"),
		},
		{
			Name:            "LoginError",
			PageNum:         1,
			MockVerifyError: nil,
			MockLoginError:  errors.New("login error"),
			ExpectedError:   fmt.Errorf("usecase.CreatorGetCoupons get logindata error: %w", errors.New("login error")),
		},
		{
			Name:            "GetCouponsError",
			PageNum:         1,
			MockVerifyError: nil,
			MockLoginUID:    101,
			MockLoginError:  nil,
			GetCouponsError: errors.New("db error"),
			ExpectedError:   fmt.Errorf("usecase.CreatorGetCoupons get coupons error: %w", errors.New("db error")),
		},
		{
			Name:            "GetCouponsSuccess",
			PageNum:         1,
			MockVerifyError: nil,
			MockLoginUID:    101,
			MockLoginError:  nil,
			CouponsFromDB: []*model.Coupon{
				{Id: 1, Uid: 101, Name: "Coupon 1"},
				{Id: 2, Uid: 101, Name: "Coupon 2"},
			},
			ExpectedError: nil,
			ExpectedCoupons: []*model.Coupon{
				{Id: 1, Uid: 101, Name: "Coupon 1"},
				{Id: 2, Uid: 101, Name: "Coupon 2"},
			},
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// First, patch the verification call.
			// In the use case the code calls:
			//    if err = uc.svc.Verify(uc.svc.VerifyPageNum(pageNum)); err != nil { ... }
			// So we patch the Verify method to return the desired error.
			mockey.Mock((*service.CommodityService).Verify).Return(tc.MockVerifyError).Build()

			// Patch login data.
			mockey.Mock(contextLogin.GetLoginData).Return(tc.MockLoginUID, tc.MockLoginError).Build()

			// Create a dummy db instance.
			gormDB := new(gorm.DB)
			db := mysql.NewCommodityDB(gormDB)
			// Patch GetCouponsByCreatorId.
			mockey.Mock(mockey.GetMethod(db, "GetCouponsByCreatorId")).Return(tc.CouponsFromDB, tc.GetCouponsError).Build()

			uc := &useCase{
				svc: new(service.CommodityService),
				db:  db,
			}

			coupons, err := uc.GetCreatorCoupons(ctx.Background(), tc.PageNum)
			if err != nil && tc.ExpectedError != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
			convey.So(coupons, convey.ShouldResemble, tc.ExpectedCoupons)
		})
	}
}

func TestUseCase_CreateUserCoupon(t *testing.T) {
	type TestCase struct {
		Name string
		// For svc.Verify(VerifyRemainUses)
		MockVerifyError error
		// For login data
		MockLoginUID   int64
		MockLoginError error
		// For db.GetCouponById(coupon.CouponId)
		GetCouponFound bool
		GetCouponError error
		// For db.CreateUserCoupon(coupon)
		MockCreateUserCouponError error
		// Expected error returned from CreateUserCoupon
		ExpectedError error
	}

	// Create a sample user coupon input.
	userCoupon := &model.UserCoupon{
		CouponId:      1,
		RemainingUses: 5,
		// Uid will be set by the use case.
	}

	testCases := []TestCase{
		{
			Name:            "VerifyError",
			MockVerifyError: errors.New("verify error"),
			ExpectedError:   errors.New("verify error"),
		},
		{
			Name:            "LoginError",
			MockVerifyError: nil,
			MockLoginError:  errors.New("login error"),
			ExpectedError:   fmt.Errorf("usecase.UserGetCoupons get logindata error: %w", errors.New("login error")),
		},
		{
			Name:            "GetCouponByIdError",
			MockVerifyError: nil,
			MockLoginUID:    101,
			MockLoginError:  nil,
			GetCouponError:  errors.New("get coupon error"),
			ExpectedError:   fmt.Errorf("usecase.DeleteCoupon error: %w", errors.New("get coupon error")),
		},
		{
			Name:            "CouponNotFound",
			MockVerifyError: nil,
			MockLoginUID:    101,
			MockLoginError:  nil,
			GetCouponFound:  false,
			GetCouponError:  nil,
			ExpectedError:   errno.ParamVerifyError,
		},
		{
			Name:                      "CreateUserCouponError",
			MockVerifyError:           nil,
			MockLoginUID:              101,
			MockLoginError:            nil,
			GetCouponFound:            true,
			GetCouponError:            nil,
			MockCreateUserCouponError: errors.New("create user coupon error"),
			ExpectedError:             fmt.Errorf("usecase.UserGetCoupon error: %w", errors.New("create user coupon error")),
		},
		{
			Name:            "CreateUserCouponSuccess",
			MockVerifyError: nil,
			MockLoginUID:    101,
			MockLoginError:  nil,
			GetCouponFound:  true,
			GetCouponError:  nil,
			// No error when creating user coupon.
			MockCreateUserCouponError: nil,
			ExpectedError:             nil,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// Patch the verification call for remaining uses.
			mockey.Mock((*service.CommodityService).Verify).Return(tc.MockVerifyError).Build()

			// Patch login data retrieval.
			mockey.Mock(contextLogin.GetLoginData).Return(tc.MockLoginUID, tc.MockLoginError).Build()

			// Create a dummy db instance.
			db := mysql.NewCommodityDB(new(gorm.DB))
			// Patch GetCouponById: returns (found, coupon, error).
			mockey.Mock(mockey.GetMethod(db, "GetCouponById")).
				Return(tc.GetCouponFound, &model.Coupon{Id: userCoupon.CouponId}, tc.GetCouponError).Build()
			// Patch CreateUserCoupon.
			mockey.Mock(mockey.GetMethod(db, "CreateUserCoupon")).
				Return(tc.MockCreateUserCouponError).Build()

			uc := &useCase{
				svc: new(service.CommodityService),
				db:  db,
			}

			// Create a copy of userCoupon to avoid side effects.
			couponCopy := *userCoupon
			err := uc.CreateUserCoupon(ctx.Background(), &couponCopy)
			if err != nil && tc.ExpectedError != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
			// On success, the coupon's Uid should be set.
			if err == nil {
				convey.So(couponCopy.Uid, convey.ShouldEqual, tc.MockLoginUID)
			}
		})
	}
}

func TestUseCase_SearchUserCoupons(t *testing.T) {
	type TestCase struct {
		Name    string
		PageNum int64
		// For svc.Verify(VerifyPageNum)
		MockVerifyError error
		// For login data retrieval.
		MockLoginUID   int64
		MockLoginError error
		// For db.GetUserCouponsByUId.
		UserCouponList               []*model.UserCoupon
		MockGetUserCouponsByUIdError error
		// For svc.GetCouponsByUserCoupons.
		ExpectedCoupons                  []*model.Coupon
		MockGetCouponsByUserCouponsError error
		// Expected overall error.
		ExpectedError error
	}

	testCases := []TestCase{
		{
			Name:            "VerifyError",
			PageNum:         1,
			MockVerifyError: errors.New("invalid page number"),
			ExpectedError:   errors.New("invalid page number"),
		},
		{
			Name:            "LoginError",
			PageNum:         1,
			MockVerifyError: nil,
			MockLoginError:  errors.New("login error"),
			ExpectedError:   fmt.Errorf("usecase.UserGetCoupons get logindata error: %w", errors.New("login error")),
		},
		{
			Name:                         "GetUserCouponsByUIdError",
			PageNum:                      1,
			MockVerifyError:              nil,
			MockLoginUID:                 101,
			MockLoginError:               nil,
			MockGetUserCouponsByUIdError: errors.New("db error"),
			ExpectedError:                fmt.Errorf("usecase.UserGetCoupons error: %w", errors.New("db error")),
		},
		{
			Name:            "GetCouponsByUserCouponsError",
			PageNum:         1,
			MockVerifyError: nil,
			MockLoginUID:    101,
			MockLoginError:  nil,
			UserCouponList: []*model.UserCoupon{
				{CouponId: 1},
				{CouponId: 2},
			},
			MockGetCouponsByUserCouponsError: errors.New("svc error"),
			ExpectedError:                    fmt.Errorf("usecase.UserGetCoupons error: %w", errors.New("svc error")),
		},
		{
			Name:            "SearchUserCouponsSuccess",
			PageNum:         1,
			MockVerifyError: nil,
			MockLoginUID:    101,
			MockLoginError:  nil,
			UserCouponList: []*model.UserCoupon{
				{CouponId: 1},
				{CouponId: 2},
			},
			// svc.GetCouponsByUserCoupons returns the list of coupons.
			ExpectedCoupons: []*model.Coupon{
				{Id: 1, Name: "Coupon 1"},
				{Id: 2, Name: "Coupon 2"},
			},
			ExpectedError: nil,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// Patch the verification for page number.
			mockey.Mock((*service.CommodityService).Verify).Return(tc.MockVerifyError).Build()
			// Patch login data retrieval.
			mockey.Mock(contextLogin.GetLoginData).Return(tc.MockLoginUID, tc.MockLoginError).Build()

			// Create a dummy db instance.
			db := mysql.NewCommodityDB(new(gorm.DB))
			// Patch GetUserCouponsByUId.
			mockey.Mock(mockey.GetMethod(db, "GetUserCouponsByUId")).
				Return(tc.UserCouponList, tc.MockGetUserCouponsByUIdError).Build()

			// Patch svc.GetCouponsByUserCoupons.
			mockey.Mock((*service.CommodityService).GetCouponsByUserCoupons).
				Return(tc.ExpectedCoupons, tc.MockGetCouponsByUserCouponsError).Build()

			uc := &useCase{
				svc: new(service.CommodityService),
				db:  db,
			}

			coupons, err := uc.SearchUserCoupons(ctx.Background(), tc.PageNum)
			if err != nil && tc.ExpectedError != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldEqual, tc.ExpectedError)
			}
			convey.So(coupons, convey.ShouldResemble, tc.ExpectedCoupons)
		})
	}
}
