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

package es

import (
	"bytes"
	"context"
	"net/http"

	"github.com/bytedance/sonic"

	"github.com/west2-online/DomTok/pkg/errno"
)

// CreateIndexForSearch CreateIndex 根据给定 indexName 创建带IK分词器的索引
func (r *Elasticsearch) CreateIndexForSearch(ctx context.Context, indexName string) error {
	// 判断索引是否存在
	resp, err := r.client.Indices.Exists([]string{indexName})
	if err != nil {
		return errno.Errorf(errno.InternalESErrorCode, "es.CreateIndex: failed to check index existence: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		// 只要不是 404，说明索引已经存在或出错
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		return errno.Errorf(errno.InternalESErrorCode, "es.CreateIndex: failed to check index existence, status code: %d", resp.StatusCode)
	}

	// 构造索引的设置和映射
	body := map[string]interface{}{
		"settings": map[string]interface{}{
			"analysis": map[string]interface{}{
				"tokenizer": map[string]interface{}{
					"ik_max_word": map[string]interface{}{
						"type": "ik_max_word", // IK Analyzer for Chinese
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "long",
				},
				"name": map[string]interface{}{
					"type":     "text",
					"analyzer": "ik_max_word", // Use the custom Chinese tokenizer
				},
				"weight": map[string]interface{}{
					"type": "float", // Store weight as float for sorting purposes
				},
			},
		},
	}

	// 将 body 序列化为 JSON
	data, err := sonic.Marshal(body)
	if err != nil {
		return errno.Errorf(errno.InternalServiceErrorCode, "es.CreateIndex: failed to marshal: %v", err)
	}

	// 创建索引
	createResp, err := r.client.Indices.Create(
		indexName,
		r.client.Indices.Create.WithContext(ctx),
		r.client.Indices.Create.WithBody(bytes.NewReader(data)),
	)
	if err != nil {
		return errno.Errorf(errno.InternalESErrorCode, "es.CreateIndex: failed to create index: %v", err)
	}
	defer createResp.Body.Close()
	if createResp.IsError() {
		return errno.Errorf(errno.InternalESErrorCode, "es.CreateIndex: failed to create index %q, status: %s", indexName, createResp.Status())
	}
	return nil
}
