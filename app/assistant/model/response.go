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

package model

import "encoding/json"

// Response is a struct that represents the response of the websocket.
type Response struct {
	Meta map[string]interface{} `json:"meta" form:"meta" query:"meta"`
	Data interface{}            `json:"data" form:"data" query:"data"`
}

// NewResponse creates a new Response.
func NewResponse() *Response {
	return &Response{
		Meta: make(map[string]interface{}),
	}
}

// SetMeta sets the value of the meta field.
func (r *Response) SetMeta(key string, value interface{}) {
	r.Meta[key] = value
}

// SetData sets the value of the data field.
func (r *Response) SetData(value interface{}) {
	r.Data = value
}

// GetMeta returns the value of the meta field.
func (r *Response) GetMeta(key string) interface{} {
	return r.Meta[key]
}

// GetData returns the value of the data field.
func (r *Response) GetData() interface{} {
	return r.Data
}

// Marshal marshals the response into a byte slice.
func (r *Response) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// MustMarshal marshals the response into a byte slice, ignoring errors.
func (r *Response) MustMarshal() []byte {
	// prevent nil pointer dereference
	if r.Meta == nil {
		r.Meta = make(map[string]interface{})
	}
	// prevent nil pointer dereference
	if r.Data == nil {
		r.Data = struct{}{}
	}
	b, _ := r.Marshal()
	return b
}

type ConnectSuccess struct {
	DialogID string `json:"dialog_id" form:"dialog_id" query:"dialog_id"`
	TZ       string `json:"tz"        form:"tz"        query:"tz"`
}

func NewConnectSuccess(dialogID string, tz string) ConnectSuccess {
	return ConnectSuccess{
		DialogID: dialogID,
		TZ:       tz,
	}
}

type ErrorData struct {
	Code  int64  `json:"code"  form:"code"  query:"code"`
	Error string `json:"error" form:"error" query:"error"`
}

func NewErrorData(code int64, err string) ErrorData {
	return ErrorData{
		Code:  code,
		Error: err,
	}
}

type DeltaContent struct {
	Delta string `json:"delta" form:"delta" query:"delta"`
	Index int64  `json:"index" form:"index" query:"index"`
	Turn  int64  `json:"turn"  form:"turn"  query:"turn"`
}

func NewDeltaContent(delta string, index int64, turn int64) DeltaContent {
	return DeltaContent{
		Delta: delta,
		Index: index,
		Turn:  turn,
	}
}

type DialogOp struct {
	Content string `json:"content" form:"content" query:"content"`
	Turn    int64  `json:"turn"    form:"turn"    query:"turn"`
}

func NewDialogOp(content string, turn int64) DialogOp {
	return DialogOp{
		Content: content,
		Turn:    turn,
	}
}
