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
)

// sf可以生成id,详见user/domain/service/service.go

func (svc *PaymentService) CreatePaymentInfo(ctx context.Context, paramToken string) (int64, error) {
	return 0, nil
}
func (svc *PaymentService) GeneratePaymentToken(ctx context.Context, paramToken string) (string, int64, error) {
	return "", 0, nil
}

// StorePaymentToken 这里的返回值还没有想好，是返回状态码还是消息字段？
func (svc *PaymentService) StorePaymentToken(ctx context.Context, paramToken string) (int, error) {
	return 0, nil
}
