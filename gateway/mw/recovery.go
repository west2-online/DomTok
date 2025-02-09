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

package mw

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/west2-online/DomTok/pkg/errno"
	"github.com/west2-online/DomTok/pkg/logger"
)

func RecoveryMW() app.HandlerFunc {
	return recovery.Recovery(recovery.WithRecoveryHandler(recoveryHandler))
}

func recoveryHandler(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {
	logger.Errorf("[Recovery] InternalServiceError err=%v\n stack=%s\n", err, stack)
	c.JSON(consts.StatusInternalServerError, map[string]interface{}{
		"code":    errno.InternalServiceErrorCode,
		"message": "内部服务错误，请稍后再试",
	})
}
