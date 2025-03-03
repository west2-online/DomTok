package mysql

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/app/commodity/domain/repository"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

var _db repository.CommodityDB

func initDB() {
	gorm, err := client.InitMySQL()
	if err != nil {
		panic(err)
	}
	_db = NewCommodityDB(gorm)
}

func initConfig() bool {
	if !utils.EnvironmentEnable() {
		return false
	}
	logger.Ignore()
	config.Init("commodity-category-test")
	initDB()
	return true
}

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
	Convey("TestCommodityDB_ViewCategory", t, func() {
		resp, err := _db.ViewCategory(ctx, 1, 10)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeEmpty)
	})
}
