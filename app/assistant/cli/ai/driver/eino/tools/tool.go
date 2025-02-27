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
	"strings"

	"github.com/cloudwego/eino/schema"
)

const (
	TagJSON         = "json"
	TagDesc         = "desc"
	TagRequired     = "required"
	TagRequiredTrue = "true"
	TagEnum         = "enum"
	TagEnumSplit    = ","
)

func reflectString(s reflect.StructField) *schema.ParameterInfo {
	p := schema.ParameterInfo{
		Type:     schema.String,
		Desc:     s.Tag.Get(TagDesc),
		Required: s.Tag.Get(TagRequired) == TagRequiredTrue,
	}
	tagEnum := s.Tag.Get(TagEnum)
	if tagEnum != "" {
		p.Enum = strings.Split(tagEnum, TagEnumSplit)
	}
	return &p
}

func reflectInteger(i reflect.StructField) *schema.ParameterInfo {
	return &schema.ParameterInfo{
		Type:     schema.Integer,
		Desc:     i.Tag.Get(TagDesc),
		Required: i.Tag.Get(TagRequired) == TagRequiredTrue,
	}
}

func reflectNumber(n reflect.StructField) *schema.ParameterInfo {
	return &schema.ParameterInfo{
		Type:     schema.Number,
		Desc:     n.Tag.Get(TagDesc),
		Required: n.Tag.Get(TagRequired) == TagRequiredTrue,
	}
}

func reflectBoolean(b reflect.StructField) *schema.ParameterInfo {
	return &schema.ParameterInfo{
		Type:     schema.Boolean,
		Desc:     b.Tag.Get(TagDesc),
		Required: b.Tag.Get(TagRequired) == TagRequiredTrue,
	}
}

func reflectObject(o reflect.StructField) *schema.ParameterInfo {
	p := &schema.ParameterInfo{
		Type:      schema.Object,
		Desc:      o.Tag.Get(TagDesc),
		SubParams: make(map[string]*schema.ParameterInfo),
		Required:  o.Tag.Get(TagRequired) == TagRequiredTrue,
	}
	for i := 0; i < o.Type.NumField(); i++ {
		f := o.Type.Field(i)
		jsonTag := f.Tag.Get(TagJSON)
		if jsonTag == "" {
			jsonTag = strings.ToLower(f.Name)
		}
		p.SubParams[jsonTag] = reflectAny(f)
	}
	return p
}

func reflectArray(a reflect.StructField) *schema.ParameterInfo {
	p := &schema.ParameterInfo{
		Type:     schema.Array,
		Desc:     a.Tag.Get(TagDesc),
		Required: a.Tag.Get(TagRequired) == TagRequiredTrue,
	}
	elemType := a.Type.Elem()
	elemField := reflect.StructField{
		Type: elemType,
	}
	p.ElemInfo = reflectAny(elemField)

	return p
}

func reflectAny(f reflect.StructField) *schema.ParameterInfo {
	switch f.Type.Kind() {
	case reflect.String:
		return reflectString(f)
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		return reflectInteger(f)
	case reflect.Float64, reflect.Float32:
		return reflectNumber(f)
	case reflect.Bool:
		return reflectBoolean(f)
	case reflect.Struct:
		return reflectObject(f)
	case reflect.Slice, reflect.Array:
		return reflectArray(f)
	default:
		panic("unsupported type: " + f.Type.String())
	}
}

// Reflect obj should be a struct instance or a struct pointer(pointer will be dereferenced automatically)
func Reflect(obj interface{}) *map[string]*schema.ParameterInfo {
	for reflect.TypeOf(obj).Kind() == reflect.Ptr {
		obj = reflect.ValueOf(obj).Elem().Interface()
	}
	if reflect.TypeOf(obj).Kind() != reflect.Struct {
		panic("obj must be a struct instance")
	}
	p := make(map[string]*schema.ParameterInfo)
	for i := 0; i < reflect.TypeOf(obj).NumField(); i++ {
		f := reflect.TypeOf(obj).Field(i)
		jsonTag := f.Tag.Get(TagJSON)
		if jsonTag == "" {
			jsonTag = strings.ToLower(f.Name)
		}
		p[jsonTag] = reflectAny(f)
	}
	return &p
}
