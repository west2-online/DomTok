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
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"github.com/west2-online/DomTok/app/order/domain/model"
	"github.com/west2-online/DomTok/app/order/domain/repository"
	"github.com/west2-online/DomTok/pkg/errno"
)

type orderDB struct {
	client *gorm.DB
}

func NewOrderDB(client *gorm.DB) repository.OrderDB {
	return &orderDB{client: client}
}

// IsOrderExist 检查订单是否存在
func (db *orderDB) IsOrderExist(ctx context.Context, orderID int64) (bool, error) {
	var count int64
	if err := db.client.WithContext(ctx).Model(&model.Order{}).
		Where("id = ?", orderID).
		Count(&count).Error; err != nil {
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to check order exist: %v", err)
	}
	return count > 0, nil
}

// CreateOrder 创建订单
func (db *orderDB) CreateOrder(ctx context.Context, order *model.Order) error {
	if err := db.client.WithContext(ctx).Create(order).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create order: %v", err)
	}
	return nil
}

// CreateOrderGoods 创建订单商品
func (db *orderDB) CreateOrderGoods(ctx context.Context, goods []*model.OrderGoods) error {
	if err := db.client.WithContext(ctx).Create(goods).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create order goods: %v", err)
	}
	return nil
}

// GetOrderByID 根据ID获取订单
func (db *orderDB) GetOrderByID(ctx context.Context, orderID int64) (*model.Order, error) {
	var order model.Order
	if err := db.client.WithContext(ctx).First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NewErrNo(errno.ServiceOrderNotFound, "order not found")
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get order: %v", err)
	}
	return &order, nil
}

// GetOrderGoodsByOrderID 获取订单商品列表
func (db *orderDB) GetOrderGoodsByOrderID(ctx context.Context, orderID int64) ([]*model.OrderGoods, error) {
	var goods []*model.OrderGoods
	if err := db.client.WithContext(ctx).Where("order_id = ?", orderID).Find(&goods).Error; err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get order goods: %v", err)
	}
	return goods, nil
}

// GetOrdersByUserID 分页获取用户订单列表
func (db *orderDB) GetOrdersByUserID(ctx context.Context, userID int64, page, size int32) ([]*model.Order, int32, error) {
	var orders []*model.Order
	var total int64

	// 1. 获取总数
	if err := db.client.WithContext(ctx).Model(&model.Order{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to count orders: %v", err)
	}

	// 2. 分页查询
	if err := db.client.WithContext(ctx).
		Where("user_id = ?", userID).
		Offset(int((page - 1) * size)).
		Limit(int(size)).
		Find(&orders).Error; err != nil {
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get orders: %v", err)
	}

	return orders, int32(total), nil
}

// UpdateOrderStatus 更新订单状态
func (db *orderDB) UpdateOrderStatus(ctx context.Context, orderID int64, status int32) error {
	if err := db.client.WithContext(ctx).Model(&model.Order{}).
		Where("id = ?", orderID).
		Update("status", status).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update order status: %v", err)
	}
	return nil
}

// UpdateOrderAddress 更新订单地址
func (db *orderDB) UpdateOrderAddress(ctx context.Context, orderID int64, addressID int64, addressInfo string) error {
	if err := db.client.WithContext(ctx).Model(&model.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"address_id":   addressID,
			"address_info": addressInfo,
		}).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update order address: %v", err)
	}
	return nil
}

// DeleteOrder 删除订单
func (db *orderDB) DeleteOrder(ctx context.Context, orderID int64) error {
	return db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 删除订单商品
		if err := tx.Where("order_id = ?", orderID).Delete(&model.OrderGoods{}).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete order goods: %v", err)
		}

		// 2. 删除订单
		if err := tx.Delete(&model.Order{}, orderID).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete order: %v", err)
		}

		return nil
	})
}

func (db *orderDB) GetOrderWithGoods(ctx context.Context, orderID int64) (*model.Order, []*model.OrderGoods, error) {
	var order Order
	var goods []OrderGoods

	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 查询订单
		if err := tx.Model(&order).Where("id=?", orderID).Find(&order).Error; err != nil {
			return errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("can't find order by id: %d", orderID))
		}

		// 查询订单商品
		if err := tx.Model(&OrderGoods{}).Where("order_id=?", orderID).Find(&goods).Error; err != nil {
			return errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("can't find order_goods by order_id: %d", orderID))
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	// 转换 goods 切片
	modelGoods := make([]*model.OrderGoods, len(goods))
	for i, g := range goods {
		modelGoods[i] = &model.OrderGoods{
			MerchantID:       g.MerchantID,
			GoodsID:          g.GoodsID,
			GoodsName:        g.GoodsName,
			GoodsHeadDrawing: g.GoodsHeadDrawing,
			StyleID:          int64(g.StyleID),
			StyleName:        g.StyleName,
			StyleHeadDrawing: g.StyleHeadDrawing,
			OriginCast:       g.OriginCast,
			SaleCast:         g.SaleCast,
			PurchaseQuantity: g.PurchaseQuantity,
			PaymentAmount:    g.PaymentAmount,
			FreightAmount:    g.FreightAmount,
			SettlementAmount: g.SettlementAmount,
			DiscountAmount:   g.DiscountAmount,
			SingleCast:       g.SingleCast,
			CouponID:         g.CouponID,
		}
	}

	return db.convertOrder(&order), modelGoods, nil
}

func (db *orderDB) convertOrder(order *Order) *model.Order {
	return &model.Order{
		Id:                    order.Id,
		Status:                order.Status,
		Uid:                   order.Uid,
		TotalAmountOfGoods:    order.TotalAmountOfGoods,
		TotalAmountOfFreight:  order.TotalAmountOfFreight,
		TotalAmountOfDiscount: order.TotalAmountOfDiscount,
		PaymentAmount:         order.PaymentAmount,
		PaymentStatus:         strconv.Itoa(int(order.PaymentStatus)),
		PaymentAt:             order.PaymentAt,
		PaymentStyle:          order.PaymentStyle,
		OrderedAt:             order.OrderedAt,
		DeletedAt:             order.DeletedAt,
		DeliveryAt:            order.DeliveryAt,
		AddressID:             order.AddressID,
		AddressInfo:           order.AddressInfo,
	}
}
