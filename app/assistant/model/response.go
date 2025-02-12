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
	Meta map[string]interface{} `thrift:"meta,1" form:"meta" json:"meta" query:"meta"`
	Data map[string]interface{} `thrift:"data,2" form:"data" json:"data" query:"data"`
}

// NewResponse creates a new Response.
func NewResponse() *Response {
	return &Response{
		Meta: make(map[string]interface{}),
		Data: make(map[string]interface{}),
	}
}

// SetMeta sets the value of the meta field.
func (r *Response) SetMeta(key string, value interface{}) {
	r.Meta[key] = value
}

// SetData sets the value of the data field.
func (r *Response) SetData(key string, value interface{}) {
	r.Data[key] = value
}

// GetMeta returns the value of the meta field.
func (r *Response) GetMeta(key string) interface{} {
	return r.Meta[key]
}

// GetData returns the value of the data field.
func (r *Response) GetData(key string) interface{} {
	return r.Data[key]
}

func (r *Response) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Response) MustMarshal() []byte {
	b, _ := r.Marshal()
	return b
}
