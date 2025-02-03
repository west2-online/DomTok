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
	"errors"
	"fmt"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"

	"github.com/west2-online/fzuhelper-server/config"
)

func NewEsClient() (*elasticsearch.Client, error) {
	if config.Elasticsearch == nil {
		return nil, errors.New("elasticsearch config is nil")
	}
	esConn := fmt.Sprintf("http://%s", config.Elasticsearch.Addr)
	cfg := elasticsearch.Config{
		Addresses: []string{esConn},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("es clint failed,error: %w", err)
	}
	return client, nil
}

func IsESConnected(es *elasticsearch.Client) bool {
	req := esapi.PingRequest{}
	_, err := req.Do(context.Background(), es)
	return err == nil
}
