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
	"testing"

	"github.com/smarty/assertions"
	"github.com/smarty/assertions/assert"
)

func TestAddAndTest(t *testing.T) {
	////arr := []string{"sova", "cypher", "jett", "sage", "KO"}
	bf := NewBloomFilter()
	////
	////for _, v := range arr {
	////	bf.Add([]byte(v))
	////}
	testItem := "KO"
	success := bf.Test([]byte(testItem))
	assert.So(success, assertions.ShouldEqual, true)

	testItem = "jett"
	failure := bf.Test([]byte(testItem))
	assert.So(failure, assertions.ShouldEqual, false)
}
