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
	"strconv"
	"sync"

	"github.com/bits-and-blooms/bloom/v3"

	"github.com/west2-online/DomTok/pkg/constants"
)

// BF 结构体表示一个布隆过滤器集合，用于高效地判断元素是否存在于集合中。
// 由于单个布隆过滤器可能会达到容量上限，因此使用多个布隆过滤器组成一个集合。
type BF struct {
	BloomFilter []*bloom.BloomFilter // 存储多个布隆过滤器的切片
	maxSize     int                  // 布隆过滤器集合的最大数量
	mu          sync.Mutex           // 互斥锁，用于保证并发安全
}

// NewBloomFilter 创建并初始化一个新的布隆过滤器集合。
// 初始时，集合中只有一个布隆过滤器。
func NewBloomFilter() *BF {
	// 创建一个新的布隆过滤器，使用常量定义的大小和误判率
	bf := bloom.NewWithEstimates(constants.BloomFilterSize, constants.FalsePositiveRate)
	return &BF{
		BloomFilter: []*bloom.BloomFilter{
			bf,
		},
		maxSize: constants.BFMaxSize, // 设置布隆过滤器集合的最大数量
	}
}

// Add 方法将一个元素添加到布隆过滤器集合中。
// 如果当前布隆过滤器达到容量上限，则创建一个新的布隆过滤器。
// 如果布隆过滤器集合的数量超过最大限制，则移除最旧的布隆过滤器。
func (b *BF) Add(uid int64) {
	item := []byte(strconv.Itoa(int(uid)))

	// 获取当前布隆过滤器集合中最后一个布隆过滤器的索引
	index := len(b.BloomFilter) - 1
	// 检查最后一个布隆过滤器是否达到容量上限
	if b.BloomFilter[index].ApproximatedSize() >= constants.BloomFilterSize {
		b.mu.Lock()
		// 创建一个新的布隆过滤器，并添加到集合中
		b.BloomFilter = append(b.BloomFilter,
			bloom.NewWithEstimates(constants.BloomFilterSize, constants.FalsePositiveRate))
		b.mu.Unlock()
	}

	// 检查布隆过滤器集合的数量是否超过最大限制
	if len(b.BloomFilter) > constants.BFMaxSize {
		b.mu.Lock()
		// 移除最旧的布隆过滤器
		b.BloomFilter = b.BloomFilter[1:]
		b.mu.Unlock()
	}

	// 将元素添加到最后一个布隆过滤器中
	b.BloomFilter[index].Add(item)
}

// Test 方法检查一个元素是否可能存在于布隆过滤器集合中。
// 由于布隆过滤器存在误判的可能，返回 true 表示元素可能存在，返回 false 表示元素一定不存在。
func (b *BF) Test(uid int64) bool {
	item := []byte(strconv.Itoa(int(uid)))

	// 遍历布隆过滤器集合中的每个布隆过滤器
	for _, filter := range b.BloomFilter {
		if filter.Test(item) {
			return true
		}
	}
	return false
}
