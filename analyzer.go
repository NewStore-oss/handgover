package handgover

import (
	"reflect"
)

type analyzeResult struct {
	Value      reflect.Value
	Properties []analyzeProperty
}

type analyzeProperty struct {
	Name  string
	Index int
	Kind  reflect.Kind
	Tag   reflect.StructTag
}

func analyze(obj interface{}) analyzeResult {
	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	result := analyzeResult{
		Value:      v,
		Properties: make([]analyzeProperty, v.NumField()),
	}

	t := result.Value.Type()
	for i := range result.Properties {
		field := t.Field(i)

		result.Properties[i] = analyzeProperty{
			Index: i,
			Name:  field.Name,
			Kind:  field.Type.Kind(),
			Tag:   field.Tag,
		}
	}
	return result
}
