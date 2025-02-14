package model_test

import (
	"encoding/json"
	"github.com/west2-online/DomTok/app/assistant/cli/ai/driver/volcengine/model"
	"testing"

	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
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

			So(err, ShouldEqual, nil)
			So(string(b), ShouldEqual, `{"type":"array","description":"test","items":[{"type":"string","description":"test"},{"type":"integer","description":"test"}]}`)
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

			So(err, ShouldEqual, nil)
			So(string(b), ShouldEqual, `{"type":"object","description":"test","properties":{"integer":{"type":"integer","description":"test"},"string":{"type":"string","description":"test"}},"required":["string"]}`)
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

			So(err, ShouldEqual, nil)
			So(string(b), ShouldEqual, `{"type":"object","properties":{"integer":{"type":"integer","description":"test"},"string":{"type":"string","description":"test"}},"required":["string"]}`)
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

			So(err, ShouldEqual, nil)
			So(string(b), ShouldEqual, `{"type":"object","properties":{"object":{"type":"object","description":"test","properties":{"integer":{"type":"integer","description":"test"},"string":{"type":"string","description":"test"}},"required":["string"]},"string":{"type":"string","description":"test"}},"required":["string"]}`)
		})
	})

}
