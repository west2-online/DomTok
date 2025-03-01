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

package tools

import (
	"reflect"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/cloudwego/eino/schema"
	. "github.com/smartystreets/goconvey/convey"
)

// func TreePrinter(p *schema.ParameterInfo, level int, printer func(...interface{})) {
// 	 if p == nil {
//		 return
//	 }
//
//	 indent := strings.Repeat("  ", level)
//	 printer(indent, p)
//	 if len(p.SubParams) != 0 {
//		 for _, v := range p.SubParams {
//			 TreePrinter(v, level+1, printer)
//		 }
//	 } else if p.ElemInfo != nil {
//		 TreePrinter(p.ElemInfo, level+1, printer)
//	 }
// }

func TestReflectInteger(t *testing.T) {
	PatchConvey("Test ReflectInteger", t, func(c C) {
		type Integer struct {
			I1 int   `desc:"a1" required:"true"`
			I2 int16 `desc:"b2" required:"false"`
			I3 int32 `          required:"true"`
			I4 int64 `desc:"d4"`
		}

		i := reflect.TypeOf(Integer{})

		So(reflectInteger(i.Field(0)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.Integer,
			Desc:     "a1",
			Required: true,
		})
		So(reflectInteger(i.Field(1)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.Integer,
			Desc:     "b2",
			Required: false,
		})
		So(reflectInteger(i.Field(2)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.Integer,
			Required: true,
		})
		So(reflectInteger(i.Field(3)), ShouldResemble, &schema.ParameterInfo{
			Type: schema.Integer,
			Desc: "d4",
		})
	})
}

func TestReflectNumber(t *testing.T) {
	PatchConvey("Test ReflectNumber", t, func(c C) {
		type Number struct {
			N1 float64 `desc:"a1" required:"true"`
			N2 float32 `desc:"b2" required:"false"`
			N3 float64 `          required:"true"`
			N4 float64 `desc:"d4"`
		}

		n := reflect.TypeOf(Number{})

		So(reflectNumber(n.Field(0)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.Number,
			Desc:     "a1",
			Required: true,
		})
		So(reflectNumber(n.Field(1)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.Number,
			Desc:     "b2",
			Required: false,
		})
		So(reflectNumber(n.Field(2)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.Number,
			Required: true,
		})
		So(reflectNumber(n.Field(3)), ShouldResemble, &schema.ParameterInfo{
			Type: schema.Number,
			Desc: "d4",
		})
	})
}

func TestReflectString(t *testing.T) {
	PatchConvey("Test ReflectString", t, func(c C) {
		type String struct {
			S1 string `desc:"a1"              required:"true"`
			S2 string `desc:"b2"              required:"false"`
			S3 string `                       required:"true"`
			S4 string `desc:"d4"`
			S5 string `desc:"e5" enum:"a,b,c"`
		}

		s := reflect.TypeOf(String{})

		So(reflectString(s.Field(0)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.String,
			Desc:     "a1",
			Required: true,
		})
		So(reflectString(s.Field(1)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.String,
			Desc:     "b2",
			Required: false,
		})
		So(reflectString(s.Field(2)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.String,
			Required: true,
		})
		So(reflectString(s.Field(3)), ShouldResemble, &schema.ParameterInfo{
			Type: schema.String,
			Desc: "d4",
		})
		So(reflectString(s.Field(4)), ShouldResemble, &schema.ParameterInfo{
			Type: schema.String,
			Desc: "e5",
			Enum: []string{"a", "b", "c"},
		})
	})
}

func TestReflectBoolean(t *testing.T) {
	PatchConvey("Test ReflectBoolean", t, func(c C) {
		type Boolean struct {
			B1 bool `desc:"a1" required:"true"`
			B2 bool `desc:"b2" required:"false"`
			B3 bool `          required:"true"`
			B4 bool `desc:"d4"`
		}

		b := reflect.TypeOf(Boolean{})

		So(reflectBoolean(b.Field(0)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.Boolean,
			Desc:     "a1",
			Required: true,
		})
		So(reflectBoolean(b.Field(1)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.Boolean,
			Desc:     "b2",
			Required: false,
		})
		So(reflectBoolean(b.Field(2)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.Boolean,
			Required: true,
		})
		So(reflectBoolean(b.Field(3)), ShouldResemble, &schema.ParameterInfo{
			Type: schema.Boolean,
			Desc: "d4",
		})
	})
}

func TestReflectObject(t *testing.T) {
	PatchConvey("Test ReflectObject(no nested)", t, func(c C) {
		type Object struct {
			O1 struct {
				A int `json:"kk" desc:"a1" required:"true"`
				B int `          desc:"b2" required:"false"`
				C int `                    required:"true"`
				D int `          desc:"d4"`
			} `desc:"o1" required:"true"`
		}

		o := reflect.TypeOf(Object{})

		So(reflectObject(o.Field(0)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.Object,
			Desc:     "o1",
			Required: true,
			SubParams: map[string]*schema.ParameterInfo{
				"kk": {
					Type:     schema.Integer,
					Desc:     "a1",
					Required: true,
				},
				"b": {
					Type:     schema.Integer,
					Desc:     "b2",
					Required: false,
				},
				"c": {
					Type:     schema.Integer,
					Required: true,
				},
				"d": {
					Type: schema.Integer,
					Desc: "d4",
				},
			},
		})
	})

	PatchConvey("Test ReflectObject(with nested)", t, func(c C) {
		type Object struct {
			O1 struct {
				A int `desc:"a1" required:"true"`
				B int `desc:"b2" required:"false"`
				C int `          required:"true"`
				D int `desc:"d4"`
				E struct {
					X int `desc:"x1" required:"true"`
					Y int `desc:"y2" required:"false"`
					Z int `          required:"true"`
					W int `desc:"w4"`
				} `desc:"e5" required:"true"`
			} `desc:"o1" required:"true"`
		}

		o := reflect.TypeOf(Object{})

		So(reflectObject(o.Field(0)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.Object,
			Desc:     "o1",
			Required: true,
			SubParams: map[string]*schema.ParameterInfo{
				"a": {
					Type:     schema.Integer,
					Desc:     "a1",
					Required: true,
				},
				"b": {
					Type:     schema.Integer,
					Desc:     "b2",
					Required: false,
				},
				"c": {
					Type:     schema.Integer,
					Required: true,
				},
				"d": {
					Type: schema.Integer,
					Desc: "d4",
				},
				"e": {
					Type:     schema.Object,
					Desc:     "e5",
					Required: true,
					SubParams: map[string]*schema.ParameterInfo{
						"x": {
							Type:     schema.Integer,
							Desc:     "x1",
							Required: true,
						},
						"y": {
							Type:     schema.Integer,
							Desc:     "y2",
							Required: false,
						},
						"z": {
							Type:     schema.Integer,
							Required: true,
						},
						"w": {
							Type: schema.Integer,
							Desc: "w4",
						},
					},
				},
			},
		})
	})

	PatchConvey("Test ReflectObject(json custom name)", t, func(c C) {
		type Object struct {
			O1 struct {
				A int `json:"x" desc:"a1" required:"true"`
				B int `json:"y" desc:"b2" required:"false"`
				C int `json:"z"           required:"true"`
				D int `json:"m" desc:"d4"`
			} `desc:"o1" required:"true"`
		}

		o := reflect.TypeOf(Object{})

		So(reflectObject(o.Field(0)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.Object,
			Desc:     "o1",
			Required: true,
			SubParams: map[string]*schema.ParameterInfo{
				"x": {
					Type:     schema.Integer,
					Desc:     "a1",
					Required: true,
				},
				"y": {
					Type:     schema.Integer,
					Desc:     "b2",
					Required: false,
				},
				"z": {
					Type:     schema.Integer,
					Required: true,
				},
				"m": {
					Type: schema.Integer,
					Desc: "d4",
				},
			},
		})
	})
}

func TestReflectArray(t *testing.T) {
	PatchConvey("Test ReflectArray", t, func(c C) {
		type Array struct {
			A1 []int `desc:"a1" required:"true"`
			A2 []int `desc:"b2" required:"false"`
			A3 []int `          required:"true"`
			A4 []int `desc:"d4"`
		}

		a := reflect.TypeOf(Array{})

		So(reflectArray(a.Field(0)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.Array,
			Desc:     "a1",
			Required: true,
			ElemInfo: &schema.ParameterInfo{
				Type: schema.Integer,
			},
		})
		So(reflectArray(a.Field(1)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.Array,
			Desc:     "b2",
			Required: false,
			ElemInfo: &schema.ParameterInfo{
				Type: schema.Integer,
			},
		})
		So(reflectArray(a.Field(2)), ShouldResemble, &schema.ParameterInfo{
			Type:     schema.Array,
			Required: true,
			ElemInfo: &schema.ParameterInfo{
				Type: schema.Integer,
			},
		})
		So(reflectArray(a.Field(3)), ShouldResemble, &schema.ParameterInfo{
			Type: schema.Array,
			Desc: "d4",
			ElemInfo: &schema.ParameterInfo{
				Type: schema.Integer,
			},
		})
	})
}

func TestReflectAny(t *testing.T) {
	PatchConvey("Test Reflect complex struct", t, func(c C) {
		PatchConvey("Struct with struct and base type", func(c C) {
			type Complex struct { //nolint: golint,unused
				A int `desc:"a1" required:"true"`
				B struct {
					X int `desc:"x1" required:"true"`
					Y int `desc:"y2" required:"false"`
					Z int `          required:"true"`
				}
				C float32 `desc:"c3"              required:"false"`
				D string  `desc:"d4" enum:"a,b,c" required:"true"`
				E struct {
					F bool `desc:"f5" required:"true"`
				}
			}

			type S struct {
				co Complex //nolint: golint,unused
			}

			co := reflect.TypeOf(S{})

			So(reflectAny(co.Field(0)), ShouldResemble, &schema.ParameterInfo{
				Type: schema.Object,
				SubParams: map[string]*schema.ParameterInfo{
					"a": {
						Type:     schema.Integer,
						Desc:     "a1",
						Required: true,
					},
					"b": {
						Type: schema.Object,
						SubParams: map[string]*schema.ParameterInfo{
							"x": {
								Type:     schema.Integer,
								Desc:     "x1",
								Required: true,
							},
							"y": {
								Type:     schema.Integer,
								Desc:     "y2",
								Required: false,
							},
							"z": {
								Type:     schema.Integer,
								Required: true,
							},
						},
					},
					"c": {
						Type:     schema.Number,
						Desc:     "c3",
						Required: false,
					},
					"d": {
						Type:     schema.String,
						Desc:     "d4",
						Enum:     []string{"a", "b", "c"},
						Required: true,
					},
					"e": {
						Type: schema.Object,
						SubParams: map[string]*schema.ParameterInfo{
							"f": {
								Type:     schema.Boolean,
								Desc:     "f5",
								Required: true,
							},
						},
					},
				},
			})
		})

		PatchConvey("Struct with array(elem is base type)", func(c C) {
			type Complex struct { //nolint: golint,unused
				A []int `desc:"a1" required:"true"`
				B []int `desc:"b2" required:"false"`
				C []int `          required:"true"`
				D []int `desc:"d4"`
			}
			type S struct {
				co Complex //nolint: golint,unused
			}

			co := reflect.TypeOf(S{})
			So(reflectAny(co.Field(0)), ShouldResemble, &schema.ParameterInfo{
				Type: schema.Object,
				SubParams: map[string]*schema.ParameterInfo{
					"a": {
						Type:     schema.Array,
						Desc:     "a1",
						Required: true,
						ElemInfo: &schema.ParameterInfo{
							Type: schema.Integer,
						},
					},
					"b": {
						Type:     schema.Array,
						Desc:     "b2",
						Required: false,
						ElemInfo: &schema.ParameterInfo{
							Type: schema.Integer,
						},
					},
					"c": {
						Type:     schema.Array,
						Required: true,
						ElemInfo: &schema.ParameterInfo{
							Type: schema.Integer,
						},
					},
					"d": {
						Type: schema.Array,
						Desc: "d4",
						ElemInfo: &schema.ParameterInfo{
							Type: schema.Integer,
						},
					},
				},
			})
		})

		PatchConvey("Struct with array(elem is struct)", func(c C) {
			type Complex struct { //nolint: golint,unused
				A []struct {
					X int `desc:"x1" required:"true"`
				} `desc:"a1" required:"true"`
				B []struct {
					Y []int `desc:"y1" required:"true"`
				} `desc:"b2" required:"false"`
				C []struct {
					Z []struct {
						W int `desc:"w1" required:"true"`
					} `required:"true"`
				} `required:"true"`
			}
			type S struct {
				co Complex //nolint: golint,unused
			}

			co := reflect.TypeOf(S{})
			So(reflectAny(co.Field(0)), ShouldResemble, &schema.ParameterInfo{
				Type: schema.Object,
				SubParams: map[string]*schema.ParameterInfo{
					"a": {
						Type:     schema.Array,
						Desc:     "a1",
						Required: true,
						ElemInfo: &schema.ParameterInfo{
							Type: schema.Object,
							SubParams: map[string]*schema.ParameterInfo{
								"x": {
									Type:     schema.Integer,
									Desc:     "x1",
									Required: true,
								},
							},
						},
					},
					"b": {
						Type:     schema.Array,
						Desc:     "b2",
						Required: false,
						ElemInfo: &schema.ParameterInfo{
							Type: schema.Object,
							SubParams: map[string]*schema.ParameterInfo{
								"y": {
									Type:     schema.Array,
									Desc:     "y1",
									Required: true,
									ElemInfo: &schema.ParameterInfo{
										Type: schema.Integer,
									},
								},
							},
						},
					},
					"c": {
						Type:     schema.Array,
						Required: true,
						ElemInfo: &schema.ParameterInfo{
							Type: schema.Object,
							SubParams: map[string]*schema.ParameterInfo{
								"z": {
									Type:     schema.Array,
									Required: true,
									ElemInfo: &schema.ParameterInfo{
										Type: schema.Object,
										SubParams: map[string]*schema.ParameterInfo{
											"w": {
												Type:     schema.Integer,
												Desc:     "w1",
												Required: true,
											},
										},
									},
								},
							},
						},
					},
				},
			})
		})
	})
}

func TestReflect(t *testing.T) {
	PatchConvey("Test Reflect", t, func(c C) {
		PatchConvey("A simple struct", func(c C) {
			type Args struct {
				Integer int64   `json:"int"  desc:"integer"              required:"true"`
				Number  float64 `json:"num"  desc:"number"               required:"false"`
				String  string  `json:"str"  desc:"string"  enum:"a,b,c" required:"true"`
				Boolean bool    `json:"bool" desc:"boolean"              required:"false"`
			}

			p := *Reflect(Args{})
			So(p, ShouldResemble, map[string]*schema.ParameterInfo{
				"int": {
					Type:     schema.Integer,
					Desc:     "integer",
					Required: true,
				},
				"num": {
					Type:     schema.Number,
					Desc:     "number",
					Required: false,
				},
				"str": {
					Type:     schema.String,
					Desc:     "string",
					Enum:     []string{"a", "b", "c"},
					Required: true,
				},
				"bool": {
					Type:     schema.Boolean,
					Desc:     "boolean",
					Required: false,
				},
			})
		})

		PatchConvey("A struct with array", func(c C) {
			type Args struct {
				Integer []int     `json:"int"  desc:"integer" required:"true"`
				Number  []float64 `json:"num"  desc:"number"  required:"false"`
				String  []string  `json:"str"  desc:"string"  required:"true"`
				Boolean []bool    `json:"bool" desc:"boolean" required:"false"`
			}

			p := *Reflect(Args{})
			So(p, ShouldResemble, map[string]*schema.ParameterInfo{
				"int": {
					Type:     schema.Array,
					Desc:     "integer",
					Required: true,
					ElemInfo: &schema.ParameterInfo{
						Type: schema.Integer,
					},
				},
				"num": {
					Type:     schema.Array,
					Desc:     "number",
					Required: false,
					ElemInfo: &schema.ParameterInfo{
						Type: schema.Number,
					},
				},
				"str": {
					Type:     schema.Array,
					Desc:     "string",
					Required: true,
					ElemInfo: &schema.ParameterInfo{
						Type: schema.String,
					},
				},
				"bool": {
					Type:     schema.Array,
					Desc:     "boolean",
					Required: false,
					ElemInfo: &schema.ParameterInfo{
						Type: schema.Boolean,
					},
				},
			})
		})

		PatchConvey("A struct with struct", func(c C) {
			type Args struct {
				Integer struct {
					A int     `json:"a" desc:"a"              required:"true"`
					B float64 `json:"b" desc:"b"              required:"false"`
					C string  `json:"c" desc:"c" enum:"a,b,c" required:"true"`
				} `desc:"integer" required:"true"`
			}

			p := *Reflect(Args{})
			So(p, ShouldResemble, map[string]*schema.ParameterInfo{
				"integer": {
					Type:     schema.Object,
					Desc:     "integer",
					Required: true,
					SubParams: map[string]*schema.ParameterInfo{
						"a": {
							Type:     schema.Integer,
							Desc:     "a",
							Required: true,
						},
						"b": {
							Type:     schema.Number,
							Desc:     "b",
							Required: false,
						},
						"c": {
							Type:     schema.String,
							Desc:     "c",
							Enum:     []string{"a", "b", "c"},
							Required: true,
						},
					},
				},
			})
		})

		PatchConvey("A struct with array of struct", func(c C) {
			type Args struct {
				Integer []struct {
					A int `json:"a" desc:"a" required:"true"`
				} `json:"int" desc:"integer" required:"true"`
			}

			p := *Reflect(Args{})
			So(p, ShouldResemble, map[string]*schema.ParameterInfo{
				"int": {
					Type:     schema.Array,
					Desc:     "integer",
					Required: true,
					ElemInfo: &schema.ParameterInfo{
						Type: schema.Object,
						SubParams: map[string]*schema.ParameterInfo{
							"a": {
								Type:     schema.Integer,
								Desc:     "a",
								Required: true,
							},
						},
					},
				},
			})
		})

		PatchConvey("A pointer of struct", func(c C) {
			type Args struct {
				Value int `json:"value" desc:"value" required:"true"`
			}
			p1 := &Args{}
			p2 := &p1
			p3 := &p2
			p4 := &p3
			reflectRes := map[string]*schema.ParameterInfo{
				"value": {
					Type:     schema.Integer,
					Desc:     "value",
					Required: true,
				},
			}
			So(*Reflect(p1), ShouldResemble, reflectRes)
			So(*Reflect(p2), ShouldResemble, reflectRes)
			So(*Reflect(p3), ShouldResemble, reflectRes)
			So(*Reflect(p4), ShouldResemble, reflectRes)
		})
	})
}
