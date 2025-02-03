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

package client

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// InitMySQL 通用初始化mysql函数，传入tableName指定表
func InitMySQL(tableName string) (db *gorm.DB, err error) {
	dsn, err := utils.GetMysqlDSN()
	if err != nil {
		return nil, fmt.Errorf("dal.InitMySQL %s:get mysql DSN error: %w", tableName, err)
	}
	db, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,  // 在执行任何 SQL 时都会创建一个 prepared statement 并将其缓存，以提高后续的效率
			SkipDefaultTransaction: false, // 不禁用默认事务(即单个创建、更新、删除时使用事务)
			TranslateError:         true,  // 允许翻译错误
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 使用单数表名
			},
			Logger: glogger.New(
				logger.GetMysqlLogger(),
				glogger.Config{
					SlowThreshold:             time.Second,  // 超过一秒的查询被认为是慢查询
					LogLevel:                  glogger.Warn, // 日志等级
					IgnoreRecordNotFoundError: true,         // 当未找到(RecordNotFoundError)时候不记录
					ParameterizedQueries:      true,         // 在 SQL 中不包含参数
					Colorful:                  false,        // 禁用颜色渲染
				}),
		})
	if err != nil {
		return nil, fmt.Errorf("dal.InitMySQL %s:mysql connect error: %w", tableName, err)
	}

	sqlDB, err := db.DB() // 尝试获取 DB 实例对象
	if err != nil {
		return nil, fmt.Errorf("get generic database object error: %w", err)
	}

	sqlDB.SetMaxIdleConns(constants.MaxIdleConns)       // 最大闲置连接数
	sqlDB.SetMaxOpenConns(constants.MaxConnections)     // 最大连接数
	sqlDB.SetConnMaxLifetime(constants.ConnMaxLifetime) // 最大可复用时间
	sqlDB.SetConnMaxIdleTime(constants.ConnMaxIdleTime) // 最长保持空闲状态时间
	db = db.Table(tableName).WithContext(context.Background())

	// 进行连通性测试
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database error: %w", err)
	}

	return db, nil
}
