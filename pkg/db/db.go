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

package db

import (
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/db/course"
	"github.com/west2-online/fzuhelper-server/pkg/db/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/db/notice"
	"github.com/west2-online/fzuhelper-server/pkg/db/user"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

type Database struct {
	client       *gorm.DB
	sf           *utils.Snowflake
	Course       *course.DBCourse
	LaunchScreen *launch_screen.DBLaunchScreen
	Notice       *notice.DBNotice
	User         *user.DBUser
}

func NewDatabase(client *gorm.DB, sf *utils.Snowflake) *Database {
	return &Database{
		client:       client,
		sf:           sf,
		Course:       course.NewDBCourse(client, sf),
		LaunchScreen: launch_screen.NewDBLaunchScreen(client, sf),
		Notice:       notice.NewDBNotice(client, sf),
		User:         user.NewDBUser(client, sf),
	}
}
