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

package service

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"

	"github.com/west2-online/DomTok/app/order/domain/model"
	basecontext "github.com/west2-online/DomTok/pkg/base/context"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

func (svc *OrderService) CreateOrder(ctx context.Context, order *model.Order, goods []*model.OrderGoods) error {
	if err := svc.db.CreateOrder(ctx, order, goods); err != nil {
		return err
	}

	if err := svc.cache.SetPaymentStatus(ctx, &model.CachePaymentStatus{
		OrderID:       order.Id,
		OrderExpire:   svc.calcOrderExpireTime(order.OrderedAt),
		PaymentStatus: order.PaymentStatus,
	}); err != nil {
		return err
	}

	return nil
}

func (svc *OrderService) MakeOrderByGoods(ctx context.Context, addressID int64, addressInfo string, goods []*model.OrderGoods) (*model.Order, error) {
	userID, err := basecontext.GetLoginData(ctx)
	if err != nil {
		return nil, err
	}

	order := &model.Order{
		Id:                    svc.nextVal(),
		Status:                constants.OrderStatusUnpaidCode,
		Uid:                   userID,
		TotalAmountOfGoods:    decimal.NewFromFloat(0), // 待计算
		TotalAmountOfFreight:  decimal.NewFromFloat(0), // 待计算
		TotalAmountOfDiscount: decimal.NewFromFloat(0), // 待计算
		PaymentAmount:         decimal.NewFromFloat(0), // 待计算
		PaymentStatus:         constants.PaymentStatusPendingCode,
		PaymentAt:             0,  // 这个值等后续被支付后才进行更新
		PaymentStyle:          "", // 等更新
		OrderedAt:             time.Now().UnixMilli(),
		DeletedAt:             0, // 默认为 null
		DeliveryAt:            0, // 等后续发货更新
		AddressID:             addressID,
		AddressInfo:           addressInfo,
		CouponId:              0,  // TODO 订单级别的优惠券应该为全局活动, 考虑在优惠券接口进行实现
		CouponName:            "", // TODO 同上
	}
	lo.ForEach(goods, func(item *model.OrderGoods, index int) {
		item.OrderID = order.Id
	})

	if err = svc.CalculateTheAmount(goods, order); err != nil {
		return nil, err
	}

	return order, nil
}

// TODO 优惠券接口完善后可以考虑接入优惠券的接口来实现这个方法的功能
func (svc *OrderService) CalculateTheAmount(goods []*model.OrderGoods, order *model.Order) error {
	// orderGoods 的  DiscountAmount PaymentAmount SinglePrice, couponName 还未赋值
	lo.ForEach(goods, func(item *model.OrderGoods, index int) {
		item.DiscountAmount = decimal.NewFromInt(0)
		item.PaymentAmount = item.TotalAmount.Add(item.DiscountAmount)
		item.SinglePrice = item.PaymentAmount.Div(decimal.NewFromInt(item.PurchaseQuantity))

		order.TotalAmountOfGoods = order.TotalAmountOfGoods.Add(item.TotalAmount)
		order.TotalAmountOfDiscount = order.TotalAmountOfDiscount.Add(item.DiscountAmount)
		order.TotalAmountOfFreight = order.TotalAmountOfFreight.Add(item.FreightAmount)
		order.PaymentAmount = order.PaymentAmount.Add(item.PaymentAmount)
	})

	// 全局活动, 应该调用 coupon 接口实现
	order.CouponId = 0
	order.CouponName = ""
	return nil
}

// DescSkuLockStock 预扣商品
func (svc *OrderService) DescSkuLockStock(ctx context.Context, orderID int64, goods []*model.OrderGoods) error {
	stocks := lo.Map(goods, func(item *model.OrderGoods, index int) *model.Stock {
		return &model.Stock{SkuID: item.StyleID, Count: item.PurchaseQuantity}
	})
	orderStock := &model.OrderStock{
		OrderID: orderID,
		Stocks:  stocks,
	}

	// 尝试预扣库存
	if err := svc.rpc.DescSkuLockStock(ctx, orderStock); err != nil {
		return err
	}

	var err error
	defer func() { // 如果操作出错了尝试进行回滚, 由于刚刚调用 DescSkuLockStock成功, 所以这里回滚大概率是成功的
		if err != nil {
			if e := svc.rpc.IncrSkuLockStock(ctx, orderStock); e != nil {
				logger.Errorf("failed to rollback for incr sku lock stock,rollbackErr: %v, caused_err: %v", e, err)
			}
		}
	}()

	var data []byte
	if data, err = svc.encodeStocks(orderStock); err != nil {
		return err
	}

	err = svc.mq.SendSyncMsg(ctx, constants.SkuStockRollbackTopic, &model.MqMessage{
		Body:       data,
		DelayLevel: constants.SkuStockRollbackTopicDelayTimeLevel,
	})
	return err
}

// encodeStocks 会把 model.OrderStock 编码成二进制丢进 msg, 可以少一次 db 查询
func (svc *OrderService) encodeStocks(stocks *model.OrderStock) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(stocks); err != nil {
		return nil, errno.NewErrNo(errno.InternalServiceErrorCode, fmt.Sprintf("failed when try encode []*model.OrderGoods to []byte, err: %v", err))
	}
	return buf.Bytes(), nil
}

// decodeStocks 从 msg 中 unmarshal 出 model.OrderStock
func (svc *OrderService) decodeStocks(data []byte) (*model.OrderStock, error) {
	var stocks *model.OrderStock
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&stocks); err != nil {
		return nil, errno.NewErrNo(errno.InternalServiceErrorCode, fmt.Sprintf("failed when try decode stocks data, err: %v", err))
	}
	return stocks, nil
}
