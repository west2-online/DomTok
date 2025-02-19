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

package context

import (
	"context"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/metadata"
)

func newContext(ctx context.Context, key string, value string) context.Context {
	return metainfo.WithPersistentValue(ctx, key, value)
}

func fromContext(ctx context.Context, key string) (string, bool) {
	return metainfo.GetPersistentValue(ctx, key)
}

// 流式传输ctx不能用传统方式传递，详情见https://www.cloudwego.io/zh/docs/kitex/tutorials/advanced-feature/metainfo/#kitex-grpc-metadata
// 返回的第一个值只去第一位是因为目前只传了uid，后续需要修改的话还要加以修正
func streamFromContext(ctx context.Context, key string) (string, bool) {
	md, success := metadata.FromIncomingContext(ctx)
	return md.Get(key)[0], success
}

func streamAppendContext(ctx context.Context, key string, value string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, key, value)
}
