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
	"fmt"

	"github.com/bytedance/sonic"

	"github.com/west2-online/DomTok/app/commodity/domain/model"
	"github.com/west2-online/DomTok/pkg/errno"
)

// UpsertSPU 如果文档已存在则累加 weight，否则插入新文档。
func (r *Elasticsearch) UpsertSPU(ctx context.Context, indexName string, spu *model.Spu, incWeight float64) error {
	body := map[string]interface{}{
		"script": map[string]interface{}{
			/*
				"source": `
				   if (ctx._source.containsKey('weight')) {
				       ctx._source.weight += params.inc
				   } else {
				       ctx._source.weight = params.inc
				   }
				`,
			*/
			"source": `
    			if (ctx._source == null) {
        			ctx._source = new HashMap();
    			}
    			if (ctx._source.containsKey('weight')) {
        			ctx._source.weight += params.inc
    			} else {
        			ctx._source.weight = params.inc
    			}
			`,
			"params": map[string]interface{}{
				"inc": incWeight,
			},
		},

		"upsert": map[string]interface{}{
			"spu_id": spu.SpuId,
			"name":   spu.Name,
			"weight": incWeight,
		},
	}

	data, err := sonic.Marshal(body)
	if err != nil {
		return errno.Errorf(errno.InternalESErrorCode, "es.UpsertSPU: failed to marshal body: %v", err)
	}

	resp, err := r.client.Update(
		indexName,
		fmt.Sprintf("%d", spu.SpuId), // 以 spuId 作为 ES 文档 _id,判断该doc是否存在
		bytes.NewReader(data),
		r.client.Update.WithContext(ctx),
	)
	if err != nil {
		return errno.Errorf(errno.InternalESErrorCode, "es.UpsertSPU: failed to update document: %v", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return errno.Errorf(errno.InternalESErrorCode, "es.UpsertSPU: update document error, status=%s", resp.Status())
	}
	return nil
}
