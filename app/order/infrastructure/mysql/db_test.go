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

package mysql

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/order/domain/repository"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

var _db repository.OrderDB

func initDB() {
	gormDB, err := client.InitMySQL()
	if err != nil {
		panic(err)
	}
	_db = NewOrderDB(gormDB)
}

func initConfig() bool {
	if !utils.EnvironmentEnable() {
		return false
	}
	logger.Ignore()
	config.Init("order-test")
	initDB()
	return true
}

// 测试了 创建订单接口，查询订单接口，删除订单接口
func TestOrderDB_CreateOrder(t *testing.T) {
	if !initConfig() {
		return
	}
	ctx := context.Background()
	order := buildTestModelOrder(t)
	orderGoods := buildTestModelOrderGoods(t, order.Id)

	Convey("TestOrderDB_CreateOrder", t, func() {
		Convey("TestOrderDB_CreateOrder_normal", func() {
			err := _db.CreateOrder(ctx, order, orderGoods)
			So(err, ShouldBeNil)

			getOrder, err := _db.GetOrderByID(ctx, order.Id)
			So(err, ShouldBeNil)
			So(getOrder.Id, ShouldEqual, getOrder.Id)
			So(getOrder.Status, ShouldEqual, order.Status)
			So(getOrder.Uid, ShouldEqual, order.Uid)
			So(getOrder.TotalAmountOfGoods.Equal(order.TotalAmountOfGoods), ShouldBeTrue)
			So(getOrder.TotalAmountOfFreight.Equal(order.TotalAmountOfFreight), ShouldBeTrue)
			So(getOrder.TotalAmountOfDiscount.Equal(order.TotalAmountOfDiscount), ShouldBeTrue)
			So(getOrder.PaymentAmount.Equal(order.PaymentAmount), ShouldBeTrue)
			So(getOrder.PaymentStatus, ShouldEqual, order.PaymentStatus)
			So(getOrder.PaymentAt, ShouldEqual, order.PaymentAt)
			So(getOrder.PaymentStyle, ShouldEqual, order.PaymentStyle)
			So(getOrder.OrderedAt, ShouldEqual, order.OrderedAt)
			So(getOrder.DeliveryAt, ShouldEqual, order.DeliveryAt)
			So(getOrder.AddressID, ShouldEqual, order.AddressID)
			So(getOrder.AddressInfo, ShouldEqual, order.AddressInfo)
		})

		Convey("TestOrderDB_CreateOrder_repeat_create", func() {
			err := _db.CreateOrder(ctx, order, orderGoods)
			So(err, ShouldNotBeNil)
		})

		Convey("TestOrderDB_CreateOrder_clear_order", func() {
			err := _db.DeleteOrder(ctx, order.Id)
			So(err, ShouldBeNil)
		})
	})
}
