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
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
)

func buildTestModelCategory(t *testing.T) *model.Category {
	t.Helper()
	return &model.Category{
		Name:      fmt.Sprintf("testcase-%d", time.Now().UnixNano()),
		CreatorId: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Now(),
	}
}

func TestCommodityDB_IsCategoryExistByName(t *testing.T) {
	if !initConfig() {
		return
	}
	ctx := context.Background()
	category := buildTestModelCategory(t)
	Convey("TestCommodityDB_CreateCategory", t, func() {
		Convey("TestCommodityDB_CreateCategory_Normal", t, func() {
			exist, err := _db.IsCategoryExistByName(ctx, category.Name)
			So(exist, ShouldBeFalse)
			So(err, ShouldBeNil)
			err = _db.CreateCategory(ctx, category)
			So(err, ShouldBeNil)
			So(category.Id, ShouldNotBeBetweenOrEqual, 0)
			So(category.Name, ShouldEqual, category.Name)
			So(category.CreatorId, ShouldNotBeEmpty)
		})
		Convey("TestCommodityDB_CreateCategory_recreate", t, func() {
			exist, _ := _db.IsCategoryExistByName(ctx, category.Name)
			So(exist, ShouldBeTrue)
		})
	})
}

func TestCommodityDB_DeleteAndUpdateCategory(t *testing.T) {
	if !initConfig() {
		return
	}
	ctx := context.Background()
	category := buildTestModelCategory(t)
	Convey("TestCommodityDB_DeleteAndUpdateCategory", t, func() {
		category, err := _db.GetCategoryById(ctx, category.Id)
		So(err, ShouldBeNil)
		err = _db.UpdateCategory(ctx, &model.Category{
			Id:   1,
			Name: "testcase2",
		})
		So(err, ShouldBeNil)
		So(category, ShouldEqual, "testcase2")
		err = _db.DeleteCategory(ctx, category)
		So(err, ShouldBeNil)
		So(category, ShouldBeEmpty)
	})
}

func TestCommodityDB_ViewCategory(t *testing.T) {
	if !initConfig() {
		return
	}
	ctx := context.Background()
	pagesize := 1
	pagenum := 10
	Convey("TestCommodityDB_ViewCategory", t, func() {
		resp, err := _db.ViewCategory(ctx, pagesize, pagenum)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeEmpty)
	})
}
