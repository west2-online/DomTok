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

package service

import (
	"context"

	"github.com/west2-online/DomTok/pkg/errno"
)

// Login logs in the user
func (s Core) Login(ctx context.Context) error {
	_, ok := ctx.Value(CtxKeyID).(string)
	if !ok {
		return errno.NewErrNoWithStack(errno.InternalServiceErrorCode, "missing id in context")
	}
	_, ok = ctx.Value(CtxKeyAccessToken).(string)
	if !ok {
		return errno.NewErrNoWithStack(errno.InternalServiceErrorCode, "missing access token in context")
	}

	return nil
}
