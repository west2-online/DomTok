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

package pack

import (
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/west2-online/DomTok/pkg/errno"
)

type Base struct {
	Code string `json:"code"`
	Msg  string `json:"message"`
}

type RespWithData struct {
	Code string `json:"code"`
	Msg  string `json:"message"`
	Data any    `json:"data"`
}

func RespError(c *app.RequestContext, err error) {
	Errno := errno.ConvertErr(err)
	c.JSON(consts.StatusOK, Base{
		Code: strconv.FormatInt(Errno.ErrorCode, 10),
		Msg:  Errno.ErrorMsg,
	})
}

func RespSuccess(c *app.RequestContext) {
	Errno := errno.Success
	c.JSON(consts.StatusOK, Base{
		Code: strconv.FormatInt(Errno.ErrorCode, 10),
		Msg:  Errno.ErrorMsg,
	})
}

func RespData(c *app.RequestContext, data any) {
	c.JSON(consts.StatusOK, RespWithData{
		Code: strconv.FormatInt(errno.SuccessCode, 10),
		Msg:  "Success",
		Data: data,
	})
}

func RespList(c *app.RequestContext, items any) {
	Errno := errno.Success
	resp := RespWithData{
		Code: strconv.FormatInt(Errno.ErrorCode, 10),
		Msg:  Errno.ErrorMsg,
		Data: items,
	}
	c.JSON(consts.StatusOK, resp)
}
