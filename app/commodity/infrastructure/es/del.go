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
	"context"
	"fmt"
)

// DeleteSpu deletes an SPU document from Elasticsearch by ID
func (r *Elasticsearch) DeleteSpu(ctx context.Context, indexName string, spuId int64) error {
	resp, err := r.client.Delete(
		indexName,
		fmt.Sprintf("%d", spuId), // Document ID as part of the URL
		r.client.Delete.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("es.DeleteSpu: failed to delete document: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("es.DeleteSpu: delete document error, status=%s", resp.Status())
	}

	return nil
}
