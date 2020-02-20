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

func setValue(kind reflect.Kind, property reflect.Value, values []string) error {
	switch kind {
	case reflect.Ptr:
		return setPointer(kind, property, values)
	case reflect.Slice:
		return setSlice(kind, property, values)
	case reflect.String:
		return setString(kind, property, values)
	case reflect.Int:
		return setInt(kind, property, values)
	case reflect.Int8:
		return setInt(kind, property, values)
	case reflect.Int16:
		return setInt(kind, property, values)
	case reflect.Int32:
		return setInt(kind, property, values)
	case reflect.Int64:
		return setInt(kind, property, values)
	case reflect.Uint:
		return setUInt(kind, property, values)
	case reflect.Uint8:
		return setUInt(kind, property, values)
	case reflect.Uint16:
		return setUInt(kind, property, values)
	case reflect.Uint32:
		return setUInt(kind, property, values)
	case reflect.Uint64:
		return setUInt(kind, property, values)
	case reflect.Bool:
		return setBool(kind, property, values)
	case reflect.Float32:
		return setFloat32(kind, property, values)
	case reflect.Float64:
		return setFloat64(kind, property, values)
	case reflect.Struct:
		return setStruct(kind, property, values)
	default:
		return fmt.Errorf("unsupported property kind %q", kind)
	}
}

func setPointer(_ reflect.Kind, property reflect.Value, values []string) error {
	underlineType := property.Type().Elem()
	property.Set(reflect.New(underlineType))
	return setValue(underlineType.Kind(), property.Elem(), values)
}

func setStruct(_ reflect.Kind, property reflect.Value, values []string) error {
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

func setString(_ reflect.Kind, property reflect.Value, values []string) error {
	property.SetString(values[0])
	return nil
}

func setSlice(kind reflect.Kind, property reflect.Value, values []string) error {
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
		if err := setValue(propertyElementKind, slice.Index(i), []string{values[i]}); err != nil {
			return err
		}
	}

	property.Set(slice)
	return nil
}

func setInt(_ reflect.Kind, property reflect.Value, values []string) error {
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

func setUInt(_ reflect.Kind, property reflect.Value, values []string) error {
	ui, err := strconv.ParseUint(values[0], 10, 64)
	if err != nil {
		return err
	}
	property.SetUint(ui)
	return nil
}

func setBool(_ reflect.Kind, property reflect.Value, values []string) error {
	b, err := strconv.ParseBool(values[0])
	if err != nil {
		return err
	}
	property.SetBool(b)
	return nil
}

func setFloat32(_ reflect.Kind, property reflect.Value, values []string) error {
	f, err := strconv.ParseFloat(values[0], 32)
	if err != nil {
		return err
	}
	property.SetFloat(f)
	return nil
}

func setFloat64(_ reflect.Kind, property reflect.Value, values []string) error {
	f, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return err
	}
	property.SetFloat(f)
	return nil
}

// Source defines the source of a given struct field tag.
//
// Tag contains the field tag name
// Get is a function to get the value/values for your given field.
type Source struct {
	Tag string
	Get func(string) ([]string, error)
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

	res := analyze(obj)

	for _, reflectedProperty := range res.Properties {
		for _, source := range sources {

			tagValue, ok := reflectedProperty.Tag.Lookup(source.Tag)
			if !ok {
				continue
			}

			property := res.Value.Field(reflectedProperty.Index)
			if !property.IsValid() || !property.CanSet() {
				continue
			}
			values, err := source.Get(tagValue)
			if err != nil {
				return newError(tagValue, source.Tag, values, err)
			}

			if len(values) == 0 {
				continue
			}

			err = setValue(reflectedProperty.Kind, property, values)
			if err != nil {
				return newError(tagValue, source.Tag, values, err)
			}
		}
	}
	return nil
}
