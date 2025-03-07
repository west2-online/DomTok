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
	"github.com/bits-and-blooms/bloom/v3"

	"github.com/west2-online/DomTok/pkg/constants"
)

type BF struct {
	BloomFilter []*bloom.BloomFilter
	maxSize     int
}

func NewBloomFilter() *BF {
	return &BF{
		BloomFilter: []*bloom.BloomFilter{
			bloom.NewWithEstimates(constants.BloomFilterSize, constants.FalsePositiveRate),
		},
		maxSize: constants.BFMaxSize,
	}
}

func (b *BF) Add(item []byte) {
	index := len(b.BloomFilter) - 1
	if b.BloomFilter[index].ApproximatedSize() >= constants.BloomFilterSize {
		b.BloomFilter = append(b.BloomFilter,
			bloom.NewWithEstimates(constants.BloomFilterSize, constants.FalsePositiveRate))
	}

	if len(b.BloomFilter) > constants.BFMaxSize {
		b.BloomFilter = b.BloomFilter[1:]
	}

	b.BloomFilter[index].Add(item)
}

func (b *BF) Test(item []byte) bool {
	for _, filter := range b.BloomFilter {
		if filter.Test(item) {
			return true
		}
	}
	return false
}

func (b *BF) Remove(item []byte) {
	for _, filter := range b.BloomFilter {
		if filter.Test(item) {
			filter.ClearAll()
			return
		}
	}
}
