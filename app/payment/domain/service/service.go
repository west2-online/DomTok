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
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	paymentStatus "github.com/west2-online/DomTok/pkg/constants"
)

// sf可以生成id,详见user/domain/service/service.go
// TODO 这个函数的逻辑不知道要怎么写
func (svc *PaymentService) CreatePaymentInfo(ctx context.Context, paramToken string) (int64, error) {
	return 0, nil
}

// GeneratePaymentToken HMAC生成支付令牌
func (svc *PaymentService) GeneratePaymentToken(ctx context.Context, paramToken string) (string, int64, error) {
	// 1. 设定过期时间为15分钟后
	expirationTime := time.Now().Add(paymentStatus.ExpirationTime * time.Minute).Unix()

	// 2. 获取 HMAC 密钥（可以从环境变量或配置文件获取）
	secretKey := []byte(paymentStatus.PaymentSecretKey)

	// 3. 计算 HMAC-SHA256 哈希
	h := hmac.New(sha256.New, secretKey)
	_, err := h.Write([]byte(fmt.Sprintf("%s:%d", paramToken, expirationTime)))
	if err != nil {
		return paymentStatus.ErrorToken, paymentStatus.ErrorExpirationTime, fmt.Errorf("failed to generate HMAC: %w", err)
	}

	// 4. 生成十六进制编码的 HMAC 签名
	token := hex.EncodeToString(h.Sum(nil))

	// 5. 返回令牌和过期时间
	return token, expirationTime, nil
}

// StorePaymentToken 这里的返回值还没有想好，是返回状态码还是消息字段？
func (svc *PaymentService) StorePaymentToken(ctx context.Context, paramToken string, expTime int64) (int, error) {
	// 1. 计算令牌的过期时间（转换成 Duration）
	expirationDuration := time.Until(time.Unix(expTime, 0))

	// 2. 存储到 Redis（key: "payment_token:<token>"，value: token）
	redisKey := fmt.Sprintf("payment_token:%s", paramToken)
	err := svc.redis.SetPaymentToken(ctx, redisKey, paramToken, expirationDuration)
	if err != nil {
		return paymentStatus.RedisStoreFailed, fmt.Errorf("failed to store payment token in redis: %w", err)
	}

	// 3. 返回成功状态码
	return paymentStatus.RedisStoreSuccess, nil
}
