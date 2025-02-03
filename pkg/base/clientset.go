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

package base

import (
	"sync"

	"github.com/cloudwego/hertz/pkg/app/client"
	elastic "github.com/elastic/go-elasticsearch"

	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

var (
	instance *ClientSet
	once     sync.Once
)

// ClientSet storage various client objects
// Notice: some or all of them maybe nil, we should check obj when use
type ClientSet struct {
	CacheClient *cache.Cache     // Redis
	ESClient    *elastic.Client  // ElasticSearch
	DBClient    *db.Database     // Database
	SFClient    *utils.Snowflake // Snowflake(DB initialize together)
	cleanups    []func()         // Functions to clean resources
	HzClient    *client.Client   // Hertz client
}

type Option func(clientSet *ClientSet)

// NewClientSet will be protected by sync.Once for ensure only 1 instance could be created in 1 lifecycle
func NewClientSet(opt ...Option) *ClientSet {
	once.Do(func() {
		var options []Option
		instance = &ClientSet{}
		options = append(options, opt...)
		for _, opt := range options {
			opt(instance)
		}
	})
	return instance
}

// Close iterates over all cleanup functions and calls them.
func (cs *ClientSet) Close() {
	for _, cleanup := range cs.cleanups {
		cleanup()
	}
}
