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

	"github.com/west2-online/DomTok/app/order/domain/model"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

// SkuLockStockRollback 要注册到 mq 中 body 是从 queue 中读取到的数据, 返回值是 是否 rollback 成功
// 一般来说被监听到这个方法里的 消息, 只有两种可能, 一种是支付模块已经消费成功提交过了, 那从 cache 中读一下状态就可以直接返回了
// 如果状态还是失败, 那就说明支付失败, 直接回滚即可,
func (svc *OrderService) SkuLockStockRollback(ctx context.Context, body []byte) (sucRollback bool) {
	orderStock, err := svc.decodeStocks(body)
	if err != nil {
		logger.Errorf(err.Error())
		return false
	}

	// 此时因为订单已经到期, 所以检查一下订单的状态, 如果还是未支付, 那就调用库存增加接口
	// 如果订单已经完成, 那就确认已经更新成功
	_, exist, err := svc.cache.GetPaymentResultRecord(ctx, orderStock.OrderID)
	if err != nil {
		logger.Errorf(err.Error())
		return false
	}

	// 大部分正常的流程到这都应该结束了, 因为支付模块会发送 mq 来试图删除这个 key
	if !exist { // 说明已经被处理过了
		return true
	}

	// 能到这说明键没被删除, payment 没有发送 mq, 也就说明这个订单还没有被回滚, 那就回滚
	if err = svc.rpc.IncrSkuLockStock(ctx, orderStock); err != nil {
		logger.Errorf(err.Error())
		return false
	}

	return true
}

// PaymentResultProcess 用于处理 payment 模块发送的 msg
//
// 按照业务上的设计, 这个方法应该比上一个更先调用, 所以上一个的实现大多是在进行一个兜底操作
func (svc *OrderService) PaymentResultProcess(ctx context.Context, body []byte) bool {
	pm, err := decodePaymentResult(body)
	if err != nil {
		logger.Errorf(err.Error())
		return true // 这里解码失败重试也解决不了, 直接丢弃
	}

	data, exist, err := svc.cache.GetPaymentResultRecord(ctx, pm.OrderID)
	if err != nil {
		logger.Errorf(err.Error())
		return false
	}

	// 因为有这个 mq 存在, 所以业务上这个 kv 是一定被注册了的
	// 过期了说明重试次数太多直到这个 kv 都过期了, 那直接跳过即可
	if !exist {
		return true
	}

	orderStock, err := svc.decodeStocks(data)
	if err != nil {
		logger.Errorf(err.Error())
		return true
	}

	// 支付成功的话
	if pm.PaymentStatus == constants.PaymentStatusSuccess {
		if err = svc.rpc.DescSkuStock(ctx, orderStock); err != nil {
			logger.Errorf(err.Error())
			return false
		}
	} else {
		// 支付失败的话
		if err = svc.rpc.IncrSkuLockStock(ctx, orderStock); err != nil {
			logger.Errorf(err.Error())
			return false
		}
	}

	// 如果没有 return说明两个 rpc 操作都成功了, 那可以删除这个 key 了, 让后续 rollback 的时候直接略过
	if err = svc.cache.DelPaymentResultRecord(ctx, pm.OrderID); err != nil {
		logger.Errorf(err.Error())
	}

	// 更新 db 状态
	if err = svc.db.UpdatePaymentStatus(ctx, pm); err != nil {
		logger.Errorf(err.Error())
	}

	return true
}

func decodePaymentResult(data []byte) (*model.PaymentResultMessage, error) {
	var pm model.PaymentResultMessage
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&pm); err != nil {
		return nil, errno.NewErrNo(errno.InternalServiceErrorCode, fmt.Sprintf("failed decode data to model.PaymentResultMessage, err: %v", err))
	}
	return &pm, nil
}
