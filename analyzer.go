// Copyright (c) 2020 NewStore GmbH <tpauling@newstore.com>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
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
