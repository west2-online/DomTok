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

type PropertyField interface{}

const (
	StringType  = "string"
	IntegerType = "integer"
	NumberType  = "number"
	BooleanType = "boolean"
	ObjectType  = "object"
	ArrayType   = "array"

	RootType = ObjectType
)

type RootParameter struct {
	Type       string                   `json:"type"`
	Properties map[string]PropertyField `json:"properties"`
	Required   []string                 `json:"required"`
}

type BaseProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ObjectProperty struct {
	Type        string                   `json:"type"`
	Description string                   `json:"description"`
	Properties  map[string]PropertyField `json:"properties"`
	Required    []string                 `json:"required"`
}

type ArrayProperty struct {
	Type        string          `json:"type"`
	Description string          `json:"description"`
	Items       []PropertyField `json:"items"`
}
