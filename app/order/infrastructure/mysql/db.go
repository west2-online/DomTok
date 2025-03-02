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

	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/west2-online/DomTok/app/order/domain/model"
	"github.com/west2-online/DomTok/app/order/domain/repository"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

type orderDB struct {
	client *gorm.DB
}

func NewOrderDB(client *gorm.DB) repository.OrderDB {
	return &orderDB{client: client}
}

// IsOrderExist 检查订单是否存在
func (db *orderDB) IsOrderExist(ctx context.Context, orderID int64) (bool, int64, error) {
	var t int64
	if err := db.client.WithContext(ctx).Model(&Order{}).
		Select("ordered_at").Where("id = ?", orderID).Find(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, t, nil
		}
		return false, t, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to check order exist: %v", err)
	}

	return true, t, nil
}

// CreateOrder 创建订单
func (db *orderDB) CreateOrder(ctx context.Context, o *model.Order, gs []*model.OrderGoods) error {
	order := db.model2Order(o)
	goods := lo.Map(gs, func(item *model.OrderGoods, index int) *OrderGoods {
		return db.model2Goods(item)
	})

	return db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create order: %v", err)
		}
		if err := tx.Create(goods).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create goods: %v", err)
		}
		return nil
	})
}

// CreateOrderGoods 创建订单商品
func (db *orderDB) CreateOrderGoods(ctx context.Context, goods []*model.OrderGoods) error {
	gs := lo.Map(goods, func(item *model.OrderGoods, index int) *OrderGoods {
		return db.model2Goods(item)
	})

	if err := db.client.WithContext(ctx).Create(gs).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create order goods: %v", err)
	}
	return nil
}

// GetOrderByID 根据ID获取订单
func (db *orderDB) GetOrderByID(ctx context.Context, orderID int64) (*model.Order, error) {
	order := &Order{Id: orderID}

	if err := db.client.WithContext(ctx).Model(order).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NewErrNo(errno.ServiceOrderNotFound, "order not found")
		}
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get order: %v", err)
	}

	return db.order2Model(order), nil
}

// GetOrderGoodsByOrderID 获取订单商品列表
func (db *orderDB) GetOrderGoodsByOrderID(ctx context.Context, orderID int64) ([]*model.OrderGoods, error) {
	var goods []*OrderGoods
	if err := db.client.WithContext(ctx).Table("order_goods").
		Where("order_id = ?", orderID).
		Find(&goods).Error; err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get order goods: %v", err)
	}

	gs := lo.Map(goods, func(item *OrderGoods, index int) *model.OrderGoods {
		return db.goods2Model(item)
	})
	return gs, nil
}

// GetOrdersByUserID 分页获取用户订单列表
func (db *orderDB) GetOrdersByUserID(ctx context.Context, userID int64, page, size int32) ([]*model.Order, int32, error) {
	var orders []*Order
	var total int64

	// 1. 获取总数
	if err := db.client.WithContext(ctx).Model(&Order{}).
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

	os := lo.Map(orders, func(item *Order, index int) *model.Order {
		return db.order2Model(item)
	})
	return os, int32(total), nil
}

func (db *orderDB) GetOrderStatus(ctx context.Context, id int64) (int8, int64, error) {
	o := Order{Id: id}
	if err := db.client.WithContext(ctx).Model(&o).Select("status,ordered_at").Scan(&o).Error; err != nil {
		return 0, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get order status: %v", err)
	}
	return o.Status, o.OrderedAt, nil
}

// UpdateOrderStatus 更新订单状态
func (db *orderDB) UpdateOrderStatus(ctx context.Context, orderID int64, status int32) error {
	if err := db.client.WithContext(ctx).Model(&Order{Id: orderID}).
		Update("status", status).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update order status: %v", err)
	}
	return nil
}

// UpdateOrderAddress 更新订单地址
func (db *orderDB) UpdateOrderAddress(ctx context.Context, orderID int64, addressID int64, addressInfo string) error {
	if err := db.client.WithContext(ctx).Model(&Order{Id: orderID}).
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
	if err := db.client.WithContext(ctx).Delete(&Order{Id: orderID}).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to delete order: %v", err)
	}
	return nil
}

func (db *orderDB) GetOrderAndGoods(ctx context.Context, orderID int64) (*model.Order, []*model.OrderGoods, error) {
	var order Order
	var goods []OrderGoods

	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 先查询订单
		if err := tx.Model(&Order{}).First(&order, orderID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errno.NewErrNo(errno.ServiceOrderNotFound, fmt.Sprintf("can't find order by id: %d", orderID))
			}
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get order: %v", err)
		}

		// 2. 再查询订单商品
		if err := tx.Table("order_goods").Where("order_id = ?", orderID).Find(&goods).Error; err != nil {
			return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to get order goods: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	// 转换 goods 切片
	modelGoods := make([]*model.OrderGoods, len(goods))
	for i, g := range goods {
		modelGoods[i] = db.goods2Model(&g)
	}

	return db.order2Model(&order), modelGoods, nil
}

func (db *orderDB) IsOrderPaid(ctx context.Context, orderID int64) (bool, error) {
	var status int8
	if err := db.client.WithContext(ctx).
		Model(&Order{Id: orderID}).
		Select("").Where("id = ?", orderID).Scan(&status); err != nil {
		return false, errno.NewErrNo(errno.InternalDatabaseErrorCode,
			fmt.Sprintf("Failed to query the payment status of an order with order id %d, err: %v", orderID, err))
	}
	return status == constants.PaymentStatusSuccessCode, nil
}

func (db *orderDB) UpdatePaymentStatus(ctx context.Context, message *model.PaymentResult) error {
	if err := db.client.WithContext(ctx).Model(&Order{Id: message.OrderID}).
		Updates(map[string]interface{}{
			"status":         message.PaymentStatus,
			"payment_status": message.PaymentStatus,
			"payment_at":     message.PaymentAt,
			"payment_style":  message.PaymentStyle,
		}).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update payment status: %v", err)
	}
	return nil
}
