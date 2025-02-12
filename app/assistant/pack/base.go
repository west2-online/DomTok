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

import "github.com/west2-online/DomTok/app/assistant/model"

type _ResponseFactory struct{}

// ResponseFactory is a global variable that points to an instance of _ResponseFactory.
var ResponseFactory _ResponseFactory

func (_ResponseFactory) Error(err error) []byte {
	resp := model.NewResponse()
	resp.SetMeta("type", "error")
	resp.SetData("message", err.Error())
	return resp.MustMarshal()
}

func (_ResponseFactory) Command(command string) []byte {
	resp := model.NewResponse()
	resp.SetMeta("type", "command")
	resp.SetData("command", command)
	return resp.MustMarshal()
}

func (_ResponseFactory) Message(data map[string]interface{}) []byte {
	resp := model.NewResponse()
	resp.SetMeta("type", "message")
	resp.Data = data
	return resp.MustMarshal()
}
