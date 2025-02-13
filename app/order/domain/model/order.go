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

package model

type Order struct {
	ID          int64  `gorm:"primarykey"`
	UserID      int64  `gorm:"not null;index"`
	AddressID   int64  `gorm:"not null"`
	AddressInfo string `gorm:"type:text"`
	Status      int32  `gorm:"not null;default:0"` // 0:待支付 1:已支付 2:已完成 3:已取消 4:未知状态
	CreatedAt   int64  `gorm:"autoCreateTime"`
	UpdatedAt   int64  `gorm:"autoUpdateTime"`
}

type OrderGoods struct {
	ID        int64   `gorm:"primarykey"`
	OrderID   int64   `gorm:"not null;index"`
	GoodsID   int64   `gorm:"not null"`
	Quantity  int32   `gorm:"not null"`
	Price     float64 `gorm:"not null"`
	CreatedAt int64   `gorm:"autoCreateTime"`
}
