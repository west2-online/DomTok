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
	"io"

	"github.com/bytedance/sonic"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/pkg/errno"
)

// SearchSpu search by weight
func (r *Elasticsearch) SearchSpu(ctx context.Context, indexName string, query string, k int) ([]*model.Spu, error) {
	body := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": query, // 'name' field
			},
		},
		"size": k, // Limit to top K results
		"sort": []map[string]interface{}{
			{"weight": map[string]interface{}{"order": "desc"}}, // Sort by weight descending
		},
	}

	data, err := sonic.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(indexName),
		r.client.Search.WithBody(bytes.NewReader(data)),
	)
	if err != nil {
		return nil, errno.Errorf(errno.InternalESErrorCode, "es.SearchSpu: search request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, errno.Errorf(errno.InternalESErrorCode, "es.SearchSpu: search request error, status=%s", resp.Status())
	}

	var respBody struct {
		Hits struct {
			Hits []struct {
				Source struct {
					SpuId  int64   `json:"spu_id"`
					Name   string  `json:"name"`
					Weight float64 `json:"weight"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errno.Errorf(errno.InternalESErrorCode, "failed to read response body: %v", err)
	}
	if err := sonic.Unmarshal(bodyBytes, &respBody); err != nil {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "es.SearchSpu: failed to unmarshal response: %v", err)
	}

	// 得到spuIdList
	var result []*model.Spu
	for _, hit := range respBody.Hits.Hits {
		result = append(result, &model.Spu{
			SpuId: hit.Source.SpuId,
			/*
				Name:  hit.Source.Name,
				Price: hit.Source.Weight,
			*/
		})
	}

	return result, nil
}
