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
	"math"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/order/domain/model"
	"github.com/west2-online/DomTok/app/order/domain/repository"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/constants"
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

func TestOrderDB_CreateOrderGoods(t *testing.T) {
	if !initConfig() {
		return
	}
	ctx := context.Background()
	order := buildTestModelOrder(t)
	orderGoods := buildTestModelOrderGoods(t, order.Id)

	Convey("TestOrderDB_CreateOrderGoods", t, func() {
		Convey("TestOrderDB_CreateOrderGoods_empty", func() {
			err := _db.CreateOrderGoods(ctx, []*model.OrderGoods{})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "empty slice found")
		})

		Convey("TestOrderDB_CreateOrderGoods_invalid_order", func() {
			invalidOrderGoods := buildTestModelOrderGoods(t, -1)
			err := _db.CreateOrderGoods(ctx, []*model.OrderGoods{invalidOrderGoods[0]})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "failed to create order goods")
		})

		Convey("TestOrderDB_CreateOrderGoods_normal", func() {
			err := _db.CreateOrderGoods(ctx, []*model.OrderGoods{orderGoods[0]})
			So(err, ShouldBeNil)

			goods, err := _db.GetOrderGoodsByOrderID(ctx, order.Id)
			So(err, ShouldBeNil)
			So(goods, ShouldNotBeNil)
			So(len(goods), ShouldEqual, 1)
			So(goods[0], ShouldNotBeNil)

			// 基本信息
			So(goods[0].OrderID, ShouldEqual, orderGoods[0].OrderID)
			So(goods[0].MerchantID, ShouldEqual, orderGoods[0].MerchantID)
			So(goods[0].GoodsID, ShouldEqual, orderGoods[0].GoodsID)
			So(goods[0].GoodsName, ShouldEqual, orderGoods[0].GoodsName)
			So(goods[0].StyleID, ShouldEqual, orderGoods[0].StyleID)
			So(goods[0].StyleName, ShouldEqual, orderGoods[0].StyleName)
			So(goods[0].GoodsVersion, ShouldEqual, orderGoods[0].GoodsVersion)
			So(goods[0].StyleHeadDrawing, ShouldEqual, orderGoods[0].StyleHeadDrawing)

			// 价格相关
			So(goods[0].OriginPrice.Equal(orderGoods[0].OriginPrice), ShouldBeTrue)
			So(goods[0].SalePrice.Equal(orderGoods[0].SalePrice), ShouldBeTrue)
			So(goods[0].SingleFreightPrice.Equal(orderGoods[0].SingleFreightPrice), ShouldBeTrue)
			So(goods[0].PurchaseQuantity, ShouldEqual, orderGoods[0].PurchaseQuantity)
			So(goods[0].TotalAmount.Equal(orderGoods[0].TotalAmount), ShouldBeTrue)
			So(goods[0].FreightAmount.Equal(orderGoods[0].FreightAmount), ShouldBeTrue)
			So(goods[0].DiscountAmount.Equal(orderGoods[0].DiscountAmount), ShouldBeTrue)
			So(goods[0].PaymentAmount.Equal(orderGoods[0].PaymentAmount), ShouldBeTrue)
			So(goods[0].SinglePrice.Equal(orderGoods[0].SinglePrice), ShouldBeTrue)

			// 优惠券信息
			So(goods[0].CouponId, ShouldEqual, orderGoods[0].CouponId)
			So(goods[0].CouponName, ShouldEqual, orderGoods[0].CouponName)
		})

		Convey("TestOrderDB_CreateOrderGoods_repeat", func() {
			err := _db.CreateOrderGoods(ctx, []*model.OrderGoods{orderGoods[0]})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "failed to create order goods")
		})

		Convey("TestOrderDB_CreateOrderGoods_clear_order", func() {
			err := _db.DeleteOrder(ctx, order.Id)
			So(err, ShouldBeNil)
		})
	})
}

func TestOrderDB_GetOrderGoodsByOrderID(t *testing.T) {
	if !initConfig() {
		return
	}
	ctx := context.Background()
	order := buildTestModelOrder(t)
	orderGoods := buildTestModelOrderGoods(t, order.Id)

	Convey("TestOrderDB_GetOrderGoodsByOrderID", t, func() {
		Convey("TestOrderDB_GetOrderGoodsByOrderID_normal", func() {
			err := _db.CreateOrder(ctx, order, orderGoods)
			So(err, ShouldBeNil)

			goods, err := _db.GetOrderGoodsByOrderID(ctx, order.Id)
			So(err, ShouldBeNil)
			So(goods, ShouldNotBeNil)
			So(len(goods), ShouldEqual, len(orderGoods))
			So(goods[0], ShouldNotBeNil)

			// 基本信息
			So(goods[0].OrderID, ShouldEqual, orderGoods[0].OrderID)
			So(goods[0].MerchantID, ShouldEqual, orderGoods[0].MerchantID)
			So(goods[0].GoodsID, ShouldEqual, orderGoods[0].GoodsID)
			So(goods[0].GoodsName, ShouldEqual, orderGoods[0].GoodsName)
			So(goods[0].StyleID, ShouldEqual, orderGoods[0].StyleID)
			So(goods[0].StyleName, ShouldEqual, orderGoods[0].StyleName)
			So(goods[0].GoodsVersion, ShouldEqual, orderGoods[0].GoodsVersion)
			So(goods[0].StyleHeadDrawing, ShouldEqual, orderGoods[0].StyleHeadDrawing)

			// 价格相关
			So(goods[0].OriginPrice.Equal(orderGoods[0].OriginPrice), ShouldBeTrue)
			So(goods[0].SalePrice.Equal(orderGoods[0].SalePrice), ShouldBeTrue)
			So(goods[0].SingleFreightPrice.Equal(orderGoods[0].SingleFreightPrice), ShouldBeTrue)
			So(goods[0].PurchaseQuantity, ShouldEqual, orderGoods[0].PurchaseQuantity)
			So(goods[0].TotalAmount.Equal(orderGoods[0].TotalAmount), ShouldBeTrue)
			So(goods[0].FreightAmount.Equal(orderGoods[0].FreightAmount), ShouldBeTrue)
			So(goods[0].DiscountAmount.Equal(orderGoods[0].DiscountAmount), ShouldBeTrue)
			So(goods[0].PaymentAmount.Equal(orderGoods[0].PaymentAmount), ShouldBeTrue)
			So(goods[0].SinglePrice.Equal(orderGoods[0].SinglePrice), ShouldBeTrue)

			// 优惠券信息
			So(goods[0].CouponId, ShouldEqual, orderGoods[0].CouponId)
			So(goods[0].CouponName, ShouldEqual, orderGoods[0].CouponName)
		})

		Convey("TestOrderDB_GetOrderGoodsByOrderID_not_exist", func() {
			goods, err := _db.GetOrderGoodsByOrderID(ctx, -1)
			So(err, ShouldBeNil)
			So(goods, ShouldNotBeNil)
			So(len(goods), ShouldEqual, 0)
		})

		Convey("TestOrderDB_GetOrderGoodsByOrderID_clear_order", func() {
			err := _db.DeleteOrder(ctx, order.Id)
			So(err, ShouldBeNil)
		})
	})
}

func TestOrderDB_GetOrderStatus(t *testing.T) {
	if !initConfig() {
		return
	}
	ctx := context.Background()
	order := buildTestModelOrder(t)
	orderGoods := buildTestModelOrderGoods(t, order.Id)

	Convey("TestOrderDB_GetOrderStatus", t, func() {
		Convey("TestOrderDB_GetOrderStatus_invalid_order", func() {
			status, orderedAt, err := _db.GetOrderStatus(ctx, -1)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, 0)    // 空记录返回 0
			So(orderedAt, ShouldEqual, 0) // 空记录返回 0
		})

		Convey("TestOrderDB_GetOrderStatus_normal", func() {
			err := _db.CreateOrder(ctx, order, orderGoods)
			So(err, ShouldBeNil)

			status, orderedAt, err := _db.GetOrderStatus(ctx, order.Id)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, order.Status)
			So(orderedAt, ShouldEqual, order.OrderedAt)
		})

		Convey("TestOrderDB_GetOrderStatus_clear_order", func() {
			err := _db.DeleteOrder(ctx, order.Id)
			So(err, ShouldBeNil)
		})
	})
}

func TestOrderDB_UpdatePaymentStatus(t *testing.T) {
	if !initConfig() {
		return
	}
	ctx := context.Background()
	order := buildTestModelOrder(t)
	orderGoods := buildTestModelOrderGoods(t, order.Id)

	Convey("TestOrderDB_UpdatePaymentStatus", t, func() {
		Convey("TestOrderDB_UpdatePaymentStatus_normal", func() {
			err := _db.CreateOrder(ctx, order, orderGoods)
			So(err, ShouldBeNil)

			paymentResult := &model.PaymentResult{
				OrderID:       order.Id,
				PaymentStyle:  "支付宝",
				PaymentAt:     time.Now().UnixMilli(),
				PaymentStatus: constants.PaymentStatusSuccessCode,
			}

			err = _db.UpdatePaymentStatus(ctx, paymentResult)
			So(err, ShouldBeNil)

			updatedOrder, err := _db.GetOrderByID(ctx, order.Id)
			So(err, ShouldBeNil)
			So(updatedOrder.PaymentStatus, ShouldEqual, constants.PaymentStatusSuccessCode)
			So(updatedOrder.PaymentStyle, ShouldEqual, paymentResult.PaymentStyle)
			So(updatedOrder.PaymentAt, ShouldEqual, paymentResult.PaymentAt)
		})

		Convey("TestOrderDB_UpdatePaymentStatus_clear_order", func() {
			err := _db.DeleteOrder(ctx, order.Id)
			So(err, ShouldBeNil)
		})
	})
}

func TestOrderDB_IsOrderExist(t *testing.T) {
	if !initConfig() {
		return
	}
	ctx := context.Background()
	order := buildTestModelOrder(t)
	orderGoods := buildTestModelOrderGoods(t, order.Id)

	Convey("TestOrderDB_IsOrderExist", t, func() {
		Convey("TestOrderDB_IsOrderExist_normal", func() {
			err := _db.CreateOrder(ctx, order, orderGoods)
			So(err, ShouldBeNil)

			exist, orderedAt, err := _db.IsOrderExist(ctx, order.Id)
			So(err, ShouldBeNil)
			So(exist, ShouldBeTrue)
			So(orderedAt, ShouldEqual, order.OrderedAt)
		})

		Convey("TestOrderDB_IsOrderExist_not_exist", func() {
			exist, orderedAt, err := _db.IsOrderExist(ctx, math.MaxInt64)
			So(err, ShouldBeNil)
			So(exist, ShouldBeFalse)
			So(orderedAt, ShouldEqual, 0) // 不存在的订单应返回0
		})

		Convey("TestOrderDB_IsOrderExist_clear_order", func() {
			err := _db.DeleteOrder(ctx, order.Id)
			So(err, ShouldBeNil)

			exist, orderedAt, err := _db.IsOrderExist(ctx, order.Id)
			So(err, ShouldBeNil)
			So(exist, ShouldBeFalse)
			So(orderedAt, ShouldEqual, 0)
		})
	})
}

func TestOrderDB_GetOrdersByUserID(t *testing.T) {
	if !initConfig() {
		return
	}
	ctx := context.Background()

	Convey("TestOrderDB_GetOrdersByUserID", t, func() {
		// 清理之前的测试数据
		err := _db.DeleteOrder(ctx, 2)
		So(err, ShouldBeNil)

		// 创建多个测试订单
		orders := make([]*model.Order, 0)
		for i := 0; i < 3; i++ {
			order := buildTestModelOrder(t)
			order.Uid = 2 // 设置相同的用户ID
			orderGoods := buildTestModelOrderGoods(t, order.Id)
			err := _db.CreateOrder(ctx, order, orderGoods)
			So(err, ShouldBeNil)
			orders = append(orders, order)
		}

		Convey("TestOrderDB_GetOrdersByUserID_normal", func() {
			list, total, err := _db.GetOrdersByUserID(ctx, 2, 1, 2)
			So(err, ShouldBeNil)
			So(total, ShouldEqual, 3)     // 总数应该是3
			So(len(list), ShouldEqual, 2) //  当前页返回2条订单

			list, total, err = _db.GetOrdersByUserID(ctx, 2, 1, 10)
			So(err, ShouldBeNil)
			So(total, ShouldEqual, 3)     // 总数应该是3
			So(len(list), ShouldEqual, 3) //  当前页返回3条订单
		})

		Convey("TestOrderDB_GetOrdersByUserID_no_orders", func() {
			list, total, err := _db.GetOrdersByUserID(ctx, 999, 1, 10)
			So(err, ShouldBeNil)
			So(total, ShouldEqual, 0)
			So(len(list), ShouldEqual, 0)
		})

		Convey("TestOrderDB_GetOrdersByUserID_invalid_page", func() {
			_, _, err := _db.GetOrdersByUserID(ctx, 1, -1, -1)
			So(err, ShouldNotBeNil)
		})

		// 清理测试数据
		Convey("TestOrderDB_GetOrdersByUserID_cleanup", func() {
			for _, order := range orders {
				err := _db.DeleteOrder(ctx, order.Id)
				So(err, ShouldBeNil)
			}
		})
	})
}

func TestOrderDB_UpdateOrderStatus(t *testing.T) {
	if !initConfig() {
		return
	}
	ctx := context.Background()
	order := buildTestModelOrder(t)
	orderGoods := buildTestModelOrderGoods(t, order.Id)

	Convey("TestOrderDB_UpdateOrderStatus", t, func() {
		Convey("TestOrderDB_UpdateOrderStatus_normal", func() {
			err := _db.CreateOrder(ctx, order, orderGoods)
			So(err, ShouldBeNil)

			newStatus := int32(2)
			err = _db.UpdateOrderStatus(ctx, order.Id, newStatus)
			So(err, ShouldBeNil)

			updatedOrder, err := _db.GetOrderByID(ctx, order.Id)
			So(err, ShouldBeNil)
			So(updatedOrder.Status, ShouldEqual, newStatus)
		})

		Convey("TestOrderDB_UpdateOrderStatus_not_exist", func() {
			err := _db.UpdateOrderStatus(ctx, math.MaxInt64, 2)
			So(err, ShouldBeNil) // GORM 在记录不存在时不会返回错误
		})

		Convey("TestOrderDB_UpdateOrderStatus_cleanup", func() {
			err := _db.DeleteOrder(ctx, order.Id)
			So(err, ShouldBeNil)
		})
	})
}

func TestOrderDB_UpdateOrderAddress(t *testing.T) {
	if !initConfig() {
		return
	}
	ctx := context.Background()
	order := buildTestModelOrder(t)
	orderGoods := buildTestModelOrderGoods(t, order.Id)

	Convey("TestOrderDB_UpdateOrderAddress", t, func() {
		Convey("TestOrderDB_UpdateOrderAddress_normal", func() {
			err := _db.CreateOrder(ctx, order, orderGoods)
			So(err, ShouldBeNil)

			newAddressID := int64(999)
			newAddressInfo := "新地址信息"
			err = _db.UpdateOrderAddress(ctx, order.Id, newAddressID, newAddressInfo)
			So(err, ShouldBeNil)

			updatedOrder, err := _db.GetOrderByID(ctx, order.Id)
			So(err, ShouldBeNil)
			So(updatedOrder.AddressID, ShouldEqual, newAddressID)
			So(updatedOrder.AddressInfo, ShouldEqual, newAddressInfo)
		})

		Convey("TestOrderDB_UpdateOrderAddress_not_exist", func() {
			err := _db.UpdateOrderAddress(ctx, math.MaxInt64, 999, "新地址")
			So(err, ShouldBeNil)
		})

		Convey("TestOrderDB_UpdateOrderAddress_cleanup", func() {
			err := _db.DeleteOrder(ctx, order.Id)
			So(err, ShouldBeNil)
		})
	})
}

func TestOrderDB_GetOrderAndGoods(t *testing.T) {
	if !initConfig() {
		return
	}
	ctx := context.Background()
	order := buildTestModelOrder(t)
	orderGoods := buildTestModelOrderGoods(t, order.Id)

	Convey("TestOrderDB_GetOrderAndGoods", t, func() {
		Convey("TestOrderDB_GetOrderAndGoods_normal", func() {
			err := _db.CreateOrder(ctx, order, orderGoods)
			So(err, ShouldBeNil)

			gotOrder, gotGoods, err := _db.GetOrderAndGoods(ctx, order.Id)
			So(err, ShouldBeNil)
			So(gotOrder, ShouldNotBeNil)
			So(gotGoods, ShouldNotBeNil)

			So(gotOrder.Id, ShouldEqual, order.Id)
			So(gotOrder.Status, ShouldEqual, order.Status)
			So(gotOrder.Uid, ShouldEqual, order.Uid)
			So(gotOrder.OrderedAt, ShouldEqual, order.OrderedAt)
			So(gotOrder.DeliveryAt, ShouldEqual, order.DeliveryAt)

			// 验证订单金额信息
			So(gotOrder.TotalAmountOfGoods.Equal(order.TotalAmountOfGoods), ShouldBeTrue)
			So(gotOrder.TotalAmountOfFreight.Equal(order.TotalAmountOfFreight), ShouldBeTrue)
			So(gotOrder.TotalAmountOfDiscount.Equal(order.TotalAmountOfDiscount), ShouldBeTrue)
			So(gotOrder.PaymentAmount.Equal(order.PaymentAmount), ShouldBeTrue)

			// 验证支付信息
			So(gotOrder.PaymentStatus, ShouldEqual, order.PaymentStatus)
			So(gotOrder.PaymentAt, ShouldEqual, order.PaymentAt)
			So(gotOrder.PaymentStyle, ShouldEqual, order.PaymentStyle)

			// 验证地址信息
			So(gotOrder.AddressID, ShouldEqual, order.AddressID)
			So(gotOrder.AddressInfo, ShouldEqual, order.AddressInfo)

			// 验证商品信息
			So(len(gotGoods), ShouldEqual, len(orderGoods))
			for i, goods := range gotGoods {
				// 基本信息
				So(goods.OrderID, ShouldEqual, orderGoods[i].OrderID)
				So(goods.MerchantID, ShouldEqual, orderGoods[i].MerchantID)
				So(goods.GoodsID, ShouldEqual, orderGoods[i].GoodsID)
				So(goods.GoodsName, ShouldEqual, orderGoods[i].GoodsName)
				So(goods.StyleID, ShouldEqual, orderGoods[i].StyleID)
				So(goods.StyleName, ShouldEqual, orderGoods[i].StyleName)
				So(goods.GoodsVersion, ShouldEqual, orderGoods[i].GoodsVersion)
				So(goods.StyleHeadDrawing, ShouldEqual, orderGoods[i].StyleHeadDrawing)

				// 价格信息
				So(goods.OriginPrice.Equal(orderGoods[i].OriginPrice), ShouldBeTrue)
				So(goods.SalePrice.Equal(orderGoods[i].SalePrice), ShouldBeTrue)
				So(goods.SingleFreightPrice.Equal(orderGoods[i].SingleFreightPrice), ShouldBeTrue)
				So(goods.SinglePrice.Equal(orderGoods[i].SinglePrice), ShouldBeTrue)
				So(goods.PurchaseQuantity, ShouldEqual, orderGoods[i].PurchaseQuantity)
				So(goods.TotalAmount.Equal(orderGoods[i].TotalAmount), ShouldBeTrue)
				So(goods.FreightAmount.Equal(orderGoods[i].FreightAmount), ShouldBeTrue)
				So(goods.DiscountAmount.Equal(orderGoods[i].DiscountAmount), ShouldBeTrue)
				So(goods.PaymentAmount.Equal(orderGoods[i].PaymentAmount), ShouldBeTrue)

				// 优惠券信息
				So(goods.CouponId, ShouldEqual, orderGoods[i].CouponId)
				So(goods.CouponName, ShouldEqual, orderGoods[i].CouponName)
			}
		})

		Convey("TestOrderDB_GetOrderAndGoods_not_exist", func() {
			_, _, err := _db.GetOrderAndGoods(ctx, math.MaxInt64)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "can't find order by id")
		})

		Convey("TestOrderDB_GetOrderAndGoods_deleted", func() {
			err := _db.DeleteOrder(ctx, order.Id)
			So(err, ShouldBeNil)
		})
	})
}
