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
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func setValue(property reflect.Value, values ...string) error {
	switch kind := property.Kind(); kind {
	case reflect.Ptr:
		return setPointer(property, values)
	case reflect.Slice:
		return setSlice(property, values)
	case reflect.String:
		return setString(property, values)
	case reflect.Int:
		return setInt(property, values)
	case reflect.Int8:
		return setInt(property, values)
	case reflect.Int16:
		return setInt(property, values)
	case reflect.Int32:
		return setInt(property, values)
	case reflect.Int64:
		return setInt(property, values)
	case reflect.Uint:
		return setUInt(property, values)
	case reflect.Uint8:
		return setUInt(property, values)
	case reflect.Uint16:
		return setUInt(property, values)
	case reflect.Uint32:
		return setUInt(property, values)
	case reflect.Uint64:
		return setUInt(property, values)
	case reflect.Bool:
		return setBool(property, values)
	case reflect.Float32:
		return setFloat32(property, values)
	case reflect.Float64:
		return setFloat64(property, values)
	case reflect.Struct:
		return setStruct(property, values)
	default:
		return fmt.Errorf("unsupported property kind %q", kind)
	}
}

func setPointer(property reflect.Value, values []string) error {
	property.Set(reflect.New(property.Type().Elem()))
	return setValue(property.Elem(), values...)
}

func setStruct(property reflect.Value, values []string) error {
	switch property.Interface().(type) {
	case time.Time:
		t, err := time.Parse(time.RFC3339, values[0])
		if err != nil {
			return err
		}
		property.Set(reflect.ValueOf(t))
		return nil
	default:
		s := reflect.New(property.Type())
		err := json.Unmarshal([]byte(values[0]), s.Interface())
		if err != nil {
			return err
		}
		property.Set(s.Elem())
	}
	return nil
}

func setString(property reflect.Value, values []string) error {
	property.SetString(values[0])
	return nil
}

func setSlice(property reflect.Value, values []string) error {
	var (
		propertyType        = property.Type()
		propertyElementKind = propertyType.Elem().Kind()
	)

	switch propertyElementKind {
	// case of a byte array
	case reflect.Uint8:
		values = strings.Split(values[0], "")
		for i, c := range values {
			values[i] = strconv.FormatUint(uint64([]byte(c)[0]), 10)
		}
	}

	var (
		lenVals = len(values)
		slice   = reflect.MakeSlice(propertyType, lenVals, lenVals)
	)

	for i := 0; i < lenVals; i++ {
		if err := setValue(slice.Index(i), values[i]); err != nil {
			return err
		}
	}

	property.Set(slice)
	return nil
}

func setInt(property reflect.Value, values []string) error {
	switch property.Interface().(type) {
	case time.Duration:
		d, err := time.ParseDuration(values[0])
		if err != nil {
			return err
		}
		property.SetInt(int64(d))
	default:
		v, err := strconv.ParseInt(values[0], 10, 64)
		if err != nil {
			return err
		}
		property.SetInt(v)
	}
	return nil
}

func setUInt(property reflect.Value, values []string) error {
	ui, err := strconv.ParseUint(values[0], 10, 64)
	if err != nil {
		return err
	}
	property.SetUint(ui)
	return nil
}

func setBool(property reflect.Value, values []string) error {
	b, err := strconv.ParseBool(values[0])
	if err != nil {
		return err
	}
	property.SetBool(b)
	return nil
}

func setFloat32(property reflect.Value, values []string) error {
	f, err := strconv.ParseFloat(values[0], 32)
	if err != nil {
		return err
	}
	property.SetFloat(f)
	return nil
}

func setFloat64(property reflect.Value, values []string) error {
	f, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return err
	}
	property.SetFloat(f)
	return nil
}

type Valuer interface {
	values() []string
}

// Values that converts a slice of strings to a Valuer interface.
func Value(v string) Valuer {
	return values([]string{v})
}

// Values that converts a slice of strings to a Valuer interface.
func Values(v []string) Valuer {
	return values(v)
}

type values []string

func (v values) values() []string {
	return v
}

// Source defines the source of a given struct field tag.
//
// Tag contains the field tag name
// Get is a function to get the value/values for your given field.
type Source struct {
	Tag string
	Get func(string) (Valuer, error)
}

type From []Source

// To takes the given sources and try to fill the fields of the given struct.
func (sources From) To(obj interface{}) error {
	if obj == nil {
		return errors.New("given struct to fill is nil")
	}

	if len(sources) == 0 {
		return nil
	}

	valueOf := reflect.ValueOf(obj)
	for valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}

	t := valueOf.Type()
	for i := 0; i < valueOf.NumField(); i++ {
		for _, source := range sources {
			field := t.Field(i)

			tagValue, ok := field.Tag.Lookup(source.Tag)
			if !ok {
				continue
			}

			property := valueOf.Field(i)
			if !property.IsValid() || !property.CanSet() {
				continue
			}

			var values []string
			v, err := source.Get(tagValue)

			if v != nil {
				values = v.values()
			}

			if err != nil {
				return newError(tagValue, source.Tag, values, err)
			}

			if len(values) == 0 {
				continue
			}

			err = setValue(property, values...)
			if err != nil {
				return newError(tagValue, source.Tag, values, err)
			}
		}
	}
	return nil
}
