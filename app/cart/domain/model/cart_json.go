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

import (
	"sort"
	"time"
)

// CartJson DBSkuJson
type CartJson struct {
	// []store
	Store []Store `json:"store"`
}

type Store struct {
	// []goods
	Goods     []Sku     `json:"sku"`
	StoreID   int64     `json:"store_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Sku struct {
	SkuID int64 `json:"sku_id"`
	Count int64 `json:"count"`
}

// SortStoresByUpdatedAt 对CartJson进行降序排序（最近的时间在前）
func (cart *CartJson) SortStoresByUpdatedAt() {
	sort.Slice(cart.Store, func(i, j int) bool {
		return cart.Store[i].UpdatedAt.After(cart.Store[j].UpdatedAt)
	})
}

// InsertSku 将sku插入json，如果sku已存在则增加count
func (cart *CartJson) InsertSku(info *GoodInfo) {
	index := -1

	// 遍历查找是否已存在 shopID
	for i, store := range cart.Store {
		if store.StoreID == info.ShopId {
			index = i
			break
		}
	}

	// shopID 存在
	if index != -1 {
		store := cart.Store[index]

		// 查找是否已存在 skuID
		skuIndex := -1
		for i, sku := range store.Goods {
			if sku.SkuID == info.SkuId {
				skuIndex = i
				break
			}
		}

		// skuID 存在，增加 count
		if skuIndex != -1 {
			store.Goods[skuIndex].Count += info.Count
		} else {
			// skuID 不存在，插入新的 sku
			store.Goods = append([]Sku{
				{
					SkuID: info.SkuId,
					Count: info.Count,
				},
			}, store.Goods...)
		}

		// 删除旧位置
		cart.Store = append(cart.Store[:index], cart.Store[index+1:]...)
		// 插到最前面
		cart.Store = append([]Store{store}, cart.Store...)
	} else {
		// shopID 不存在，创建新的 store 并插入 sku
		newStore := Store{
			StoreID: info.ShopId,
			Goods: []Sku{
				{
					SkuID: info.SkuId,
					Count: info.Count,
				},
			},
			UpdatedAt: time.Now(),
		}
		// 插到最前面
		cart.Store = append([]Store{newStore}, cart.Store...)
	}
}

// GetRecentNStores 获得CartJson中前n个店铺
func (cart *CartJson) GetRecentNStores(n int) *CartJson {
	cartJson := new(CartJson)
	if len(cart.Store) > n {
		cartJson.Store = cart.Store[:n]
	} else {
		cartJson.Store = cart.Store
	}
	return cartJson
}
