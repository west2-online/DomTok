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
	"github.com/west2-online/DomTok/app/assistant/model"
	"github.com/west2-online/DomTok/pkg/errno"
)

type _ResponseFactory struct{}

// ResponseFactory is a global variable that points to an instance of _ResponseFactory.
var ResponseFactory _ResponseFactory

const (
	MetaType        = "type"
	MetaTypeCommand = "command"
	MetaTypeMessage = "message"
	MetaTypeError   = "error"
	MetaTypePing    = "ping"
	MetaExtra       = "extra"

	DataContent   = "content"
	DataCommandOp = "op"

	OpDialog = "dialog"

	OpContentClose = "close"
	OpContentOpen  = "open"
)

func (_ResponseFactory) ConnectSuccess(extra interface{}) []byte {
	resp := model.NewResponse()
	resp.SetMeta(MetaType, MetaTypePing)
	resp.SetMeta(MetaExtra, extra)
	return resp.MustMarshal()
}

// Error returns a response with an error message.
func (_ResponseFactory) Error(err error) []byte {
	resp := model.NewResponse()
	resp.SetMeta(MetaType, MetaTypeError)
	e := errno.ConvertErr(err)
	resp.SetData(model.NewErrorData(e.ErrorCode, e.ErrorMsg))
	return resp.MustMarshal()
}

// Command returns a response with a command.
func (_ResponseFactory) Command(params interface{}) []byte {
	resp := model.NewResponse()
	resp.SetMeta(MetaType, MetaTypeCommand)
	resp.SetData(params)
	return resp.MustMarshal()
}

// Message returns a response with a message.
func (_ResponseFactory) Message(params interface{}) []byte {
	resp := model.NewResponse()
	resp.SetMeta(MetaType, MetaTypeMessage)
	resp.SetData(params)
	return resp.MustMarshal()
}
