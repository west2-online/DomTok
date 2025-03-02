package mysql

import (
	"context"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand/v2"
	"testing"

	"github.com/west2-online/DomTok/app/payment/domain/model"
	"github.com/west2-online/DomTok/app/payment/domain/repository"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

var _paymentDB repository.PaymentDB

func initDB() {
	gormDB, err := client.InitMySQL()
	if err != nil {
		panic(err)
	}
	_paymentDB = NewPaymentDB(gormDB)
}

func initConfig() bool {
	if !utils.EnvironmentEnable() {
		return false
	}
	logger.Ignore()
	config.Init("payment-test")
	initDB()
	return true
}

// 构建测试用的支付订单数据
func buildTestModelPaymentOrder(t *testing.T) *model.PaymentOrder {
	t.Helper()
	return &model.PaymentOrder{
		ID:                        rand.Int64(),
		OrderID:                   rand.Int64(),
		UserID:                    rand.Int64(),
		Amount:                    decimal.NewFromFloat(100.50), //nolint
		Status:                    constants.PaymentStatusPendingCode,
		MaskedCreditCardNumber:    "****1234",
		CreditCardExpirationYear:  2026,
		CreditCardExpirationMonth: 6,
		Description:               "Test payment order",
	}
}

// 构建测试用的退款数据
func buildTestModelPaymentRefund(t *testing.T, orderID int64) *model.PaymentRefund {
	t.Helper()
	return &model.PaymentRefund{
		ID:                        rand.Int64(),
		OrderID:                   orderID,
		UserID:                    rand.Int64(),
		RefundAmount:              decimal.NewFromFloat(50.25), //nolint
		RefundReason:              "Test refund reason",
		Status:                    constants.RefundStatusPendingCode,
		MaskedCreditCardNumber:    "****1234",
		CreditCardExpirationYear:  2026,
		CreditCardExpirationMonth: 6,
	}
}

// 测试创建支付订单和查询支付信息
func TestPaymentDB_CreateAndGetPayment(t *testing.T) {
	if !initConfig() {
		return
	}

	ctx := context.Background()
	paymentOrder := buildTestModelPaymentOrder(t)

	Convey("TestPaymentDB_CreateAndGetPayment", t, func() {
		Convey("TestPaymentDB_CreatePayment_normal", func() {
			// 测试创建支付
			err := _paymentDB.CreatePayment(ctx, paymentOrder)
			So(err, ShouldBeNil)

			// 测试检查支付是否存在
			exist, err := _paymentDB.CheckPaymentExist(ctx, paymentOrder.OrderID)
			So(err, ShouldBeNil)
			So(exist, ShouldEqual, constants.PaymentExist)

			// 测试获取支付信息
			getPayment, err := _paymentDB.GetPaymentInfo(ctx, paymentOrder.OrderID)
			So(err, ShouldBeNil)
			So(getPayment.ID, ShouldEqual, paymentOrder.ID)
			So(getPayment.OrderID, ShouldEqual, paymentOrder.OrderID)
			So(getPayment.UserID, ShouldEqual, paymentOrder.UserID)
			So(getPayment.Amount.Equal(paymentOrder.Amount), ShouldBeTrue)
			So(getPayment.Status, ShouldEqual, paymentOrder.Status)
			So(getPayment.MaskedCreditCardNumber, ShouldEqual, paymentOrder.MaskedCreditCardNumber)
			So(getPayment.CreditCardExpirationYear, ShouldEqual, paymentOrder.CreditCardExpirationYear)
			So(getPayment.CreditCardExpirationMonth, ShouldEqual, paymentOrder.CreditCardExpirationMonth)
			So(getPayment.Description, ShouldEqual, paymentOrder.Description)
		})

		Convey("TestPaymentDB_CreatePayment_repeat_create", func() {
			// 测试重复创建支付
			err := _paymentDB.CreatePayment(ctx, paymentOrder)
			So(err, ShouldNotBeNil)
		})

		Convey("TestPaymentDB_CheckPayment_not_exist", func() {
			// 测试检查不存在的支付
			nonExistOrderID := rand.Int64()
			exist, err := _paymentDB.CheckPaymentExist(ctx, nonExistOrderID)
			So(err, ShouldBeNil)
			So(exist, ShouldEqual, constants.PaymentNotExist)

			// 测试获取不存在的支付信息
			_, err = _paymentDB.GetPaymentInfo(ctx, nonExistOrderID)
			So(err, ShouldNotBeNil)
		})
	})
}

// 测试创建退款
func TestPaymentDB_CreateRefund(t *testing.T) {
	if !initConfig() {
		return
	}

	ctx := context.Background()
	paymentOrder := buildTestModelPaymentOrder(t)

	Convey("TestPaymentDB_CreateRefund", t, func() {
		Convey("TestPaymentDB_CreateRefund_normal", func() {
			// 先创建支付订单
			err := _paymentDB.CreatePayment(ctx, paymentOrder)
			So(err, ShouldBeNil)

			// 创建退款申请
			refundOrder := buildTestModelPaymentRefund(t, paymentOrder.OrderID)
			err = _paymentDB.CreateRefund(ctx, refundOrder)
			So(err, ShouldBeNil)
		})

		Convey("TestPaymentDB_CreateRefund_non_exist_payment", func() {
			// 测试为不存在的支付创建退款
			nonExistOrderID := rand.Int64()
			refundOrder := buildTestModelPaymentRefund(t, nonExistOrderID)
			err := _paymentDB.CreateRefund(ctx, refundOrder)
			So(err, ShouldNotBeNil)
		})
	})
}

// 测试转换函数
func TestConvertFunctions(t *testing.T) {
	Convey("TestConvertFunctions", t, func() {
		Convey("TestConvertToDBModel", func() {
			// 测试正常转换
			paymentOrder := buildTestModelPaymentOrder(t)
			dbModel, err := ConvertToDBModel(paymentOrder)
			So(err, ShouldBeNil)
			So(dbModel.ID, ShouldEqual, paymentOrder.ID)
			So(dbModel.OrderID, ShouldEqual, paymentOrder.OrderID)
			So(dbModel.UserID, ShouldEqual, paymentOrder.UserID)
			So(dbModel.Amount.Equal(paymentOrder.Amount), ShouldBeTrue)

			// 测试传入nil
			dbModel, err = ConvertToDBModel(nil)
			So(err, ShouldNotBeNil)
			So(dbModel, ShouldBeNil)
		})

		Convey("TestConvertRefundToDBModel", func() {
			// 测试正常转换
			refundOrder := buildTestModelPaymentRefund(t, rand.Int64())
			dbModel, err := ConvertRefundToDBModel(refundOrder)
			So(err, ShouldBeNil)
			So(dbModel.ID, ShouldEqual, refundOrder.ID)
			So(dbModel.OrderID, ShouldEqual, refundOrder.OrderID)
			So(dbModel.RefundAmount.Equal(refundOrder.RefundAmount), ShouldBeTrue)

			// 测试传入nil
			dbModel, err = ConvertRefundToDBModel(nil)
			So(err, ShouldNotBeNil)
			So(dbModel, ShouldBeNil)
		})
	})
}
