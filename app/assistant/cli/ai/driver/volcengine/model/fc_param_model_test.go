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

package model_test

import (
	"encoding/json"
	"testing"

	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/assistant/cli/ai/driver/volcengine/model"
)

func TestMarshalFcParamModelStruct(t *testing.T) {
	PatchConvey("Test if the fc param model can be marshaled correctly", t, func() {
		PatchConvey("Test root parameter with no properties", func() {
			p := model.RootParameter{
				Type:       model.RootType,
				Properties: make(map[string]model.PropertyField),
				Required:   []string{},
			}

			b, err := json.Marshal(p)

			So(err, ShouldEqual, nil)
			So(string(b), ShouldEqual, `{"type":"object","properties":{},"required":[]}`)
		})

		PatchConvey("Test base property", func() {
			types := []string{model.StringType, model.IntegerType, model.NumberType, model.BooleanType}
			for _, c := range types {
				p := model.BaseProperty{
					Type:        c,
					Description: "test",
				}

				b, err := json.Marshal(p)

				So(err, ShouldEqual, nil)
				So(string(b), ShouldEqual, `{"type":"`+c+`","description":"test"}`)
			}
		})

		PatchConvey("Test array property", func() {
			p := model.ArrayProperty{
				Type:        model.ArrayType,
				Description: "test",
				Items: []model.PropertyField{
					model.BaseProperty{
						Type:        model.StringType,
						Description: "test",
					},
					model.BaseProperty{
						Type:        model.IntegerType,
						Description: "test",
					},
				},
			}

			b, err := json.Marshal(p)

			expected := `{"type":"array"` +
				`,"description":"test"` +
				`,"items":[{"type":"string","description":"test"},{"type":"integer","description":"test"}]}`

			So(err, ShouldEqual, nil)
			So(string(b), ShouldEqual, expected)
		})

		PatchConvey("Test object property", func() {
			p := model.ObjectProperty{
				Type:        model.ObjectType,
				Description: "test",
				Properties: map[string]model.PropertyField{
					"string": model.BaseProperty{
						Type:        model.StringType,
						Description: "test",
					},
					"integer": model.BaseProperty{
						Type:        model.IntegerType,
						Description: "test",
					},
				},
				Required: []string{"string"},
			}

			b, err := json.Marshal(p)

			expected := `{"type":"object"` +
				`,"description":"test"` +
				`,"properties":{"integer":{"type":"integer","description":"test"}` +
				`,"string":{"type":"string","description":"test"}}` +
				`,"required":["string"]}`

			So(err, ShouldEqual, nil)
			So(string(b), ShouldEqual, expected)
		})

		PatchConvey("Test root parameter with properties", func() {
			p := model.RootParameter{
				Type: model.RootType,
				Properties: map[string]model.PropertyField{
					"string": model.BaseProperty{
						Type:        model.StringType,
						Description: "test",
					},
					"integer": model.BaseProperty{
						Type:        model.IntegerType,
						Description: "test",
					},
				},
				Required: []string{"string"},
			}

			b, err := json.Marshal(p)

			expected := `{"type":"object"` +
				`,"properties":{"integer":{"type":"integer","description":"test"}` +
				`,"string":{"type":"string","description":"test"}}` +
				`,"required":["string"]}`

			So(err, ShouldEqual, nil)
			So(string(b), ShouldEqual, expected)
		})

		PatchConvey("Test root parameter with properties and nested object", func() {
			p := model.RootParameter{
				Type: model.RootType,
				Properties: map[string]model.PropertyField{
					"string": model.BaseProperty{
						Type:        model.StringType,
						Description: "test",
					},
					"object": model.ObjectProperty{
						Type:        model.ObjectType,
						Description: "test",
						Properties: map[string]model.PropertyField{
							"string": model.BaseProperty{
								Type:        model.StringType,
								Description: "test",
							},
							"integer": model.BaseProperty{
								Type:        model.IntegerType,
								Description: "test",
							},
						},
						Required: []string{"string"},
					},
				},
				Required: []string{"string"},
			}

			b, err := json.Marshal(p)

			expected := `{"type":"object"` +
				`,"properties":{"object":{"type":"object","description":"test"` +
				`,"properties":{"integer":{"type":"integer","description":"test"}` +
				`,"string":{"type":"string","description":"test"}}` +
				`,"required":["string"]}` +
				`,"string":{"type":"string","description":"test"}}` +
				`,"required":["string"]}`

			So(err, ShouldEqual, nil)
			So(string(b), ShouldEqual, expected)
		})
	})
}
