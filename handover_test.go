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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFillWithNilStruct(t *testing.T) {

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"helloworld"}, nil
			},
		},
	}
	assert.Error(t, From(sources).To(nil))
}

func TestFillWithNoSource(t *testing.T) {

	var (
		s struct {
			Pointer *string `foo:"bar"`
		}
		sources []Source
	)
	assert.NoError(t, From(sources).To(&s))
}

func TestFillPointer(t *testing.T) {

	var s struct {
		Pointer *string `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"helloworld"}, nil
			},
		},
	}
	assert.NoError(t, From(sources).To(&s))

	assert.NotNil(t, s.Pointer)
	assert.Equal(t, "helloworld", *s.Pointer)
}

func TestFillSlice(t *testing.T) {

	var s struct {
		Slice   []string         `foo:"bar"`
		Bytes   []byte           `john:"doe"`
		RawJSON *json.RawMessage `john:"doe"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"hello", "world"}, nil
			},
		},
		{
			Tag: "john",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "doe", field)
				return []string{`{ "some": "json" }`}, nil
			},
		},
	}

	assert.NoError(t, From(sources).To(&s))

	assert.Equal(t, []string{"hello", "world"}, s.Slice)
	assert.Equal(t, []byte(`{ "some": "json" }`), s.Bytes)

	assert.NotNil(t, s.RawJSON)
	assert.Equal(t, json.RawMessage(`{ "some": "json" }`), *s.RawJSON)
}

func TestFillSliceWithInvalidValue(t *testing.T) {

	var s struct {
		Slice []int `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"invalid", "value"}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parsedErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
}

func TestFillString(t *testing.T) {

	var s struct {
		String string `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"helloworld"}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.NoError(t, err)

	assert.Equal(t, "helloworld", s.String)
}

func TestFillTimeDuration(t *testing.T) {

	var s struct {
		Duration time.Duration `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"1h"}, nil
			},
		},
	}

	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, time.Minute*60, s.Duration)
}

func TestFillTimeDurationWithInvalidValue(t *testing.T) {

	var s struct {
		Duration time.Duration `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"1"}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parsedErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)

}

func TestFillInt(t *testing.T) {

	var s struct {
		Int int `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"1"}, nil
			},
		},
	}
	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, int(1), s.Int)
}

func TestFillIntWithInvalidValue(t *testing.T) {

	s := struct {
		Int int `foo:"bar"`
	}{}
	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"invalid"}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parseErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parseErr.Field)
	assert.Error(t, parseErr.InnerError)
}

func TestFillInt8(t *testing.T) {

	var s struct {
		Int8 int8 `foo:"bar"`
	}
	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"1"}, nil
			},
		},
	}
	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, int8(1), s.Int8)
}

func TestFillInt8WithInvalidValue(t *testing.T) {

	var s struct {
		Int8 int8 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"invalid"}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parsedErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
}

func TestFillInt16(t *testing.T) {

	var s struct {
		Int16 int16 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"1"}, nil
			},
		},
	}

	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, int16(1), s.Int16)
}

func TestFillInt16WithInvalidValue(t *testing.T) {

	var s struct {
		Int16 int16 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"invalid"}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parsedErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
}

func TestFillInt32(t *testing.T) {

	var s struct {
		Int32 int32 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"1"}, nil
			},
		},
	}

	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, int32(1), s.Int32)
}

func TestFillInt32WithInvalidValue(t *testing.T) {

	var s struct {
		Int32 int32 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"invalid"}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parsedErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
}

func TestFillInt64(t *testing.T) {

	var s struct {
		Int64 int64 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"1"}, nil
			},
		},
	}

	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, int64(1), s.Int64)
}

func TestFillInt64WithInvalidValue(t *testing.T) {

	var s struct {
		Int64 int64 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"invalid"}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parsedErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
}

func TestFillUInt(t *testing.T) {

	var s struct {
		UInt uint `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"1"}, nil
			},
		},
	}

	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, uint(1), s.UInt)
}

func TestFillUIntWithInvalidValue(t *testing.T) {

	var s struct {
		UInt uint `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"invalid"}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parsedErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
}

func TestFillUInt8(t *testing.T) {

	var s struct {
		UInt8 uint8 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"1"}, nil
			},
		},
	}

	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, uint8(1), s.UInt8)
}

func TestFillUInt8WithInvalidValue(t *testing.T) {

	var s struct {
		UInt8 uint8 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"invalid"}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parsedErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
}

func TestFillUInt16(t *testing.T) {

	var s struct {
		UInt16 uint16 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"1"}, nil
			},
		},
	}

	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, uint16(1), s.UInt16)
}

func TestFillUInt16WithInvalidValue(t *testing.T) {

	var s struct {
		UInt16 uint16 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"invalid"}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parsedErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
}

func TestFillUInt32(t *testing.T) {

	var s struct {
		UInt32 uint32 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"1"}, nil
			},
		},
	}

	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, uint32(1), s.UInt32)
}

func TestFillUInt32WithInvalidValue(t *testing.T) {

	var s struct {
		UInt32 uint32 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"invalid"}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parsedErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
}

func TestFillUInt64(t *testing.T) {

	var s struct {
		UInt64 uint64 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"1"}, nil
			},
		},
	}

	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, uint64(1), s.UInt64)
}

func TestFillUInt64WithInvalidValue(t *testing.T) {

	var s struct {
		UInt64 uint64 `foo:"bar"`
	}
	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"invalid"}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parsedErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
}

func TestFillBool(t *testing.T) {

	var s struct {
		Bool bool `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"true"}, nil
			},
		},
	}

	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, true, s.Bool)
}

func TestFillFloat32(t *testing.T) {

	var s struct {
		Float32 float32 `foo:"bar"`
	}
	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"1.5"}, nil
			},
		},
	}

	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, float32(1.5), s.Float32)
}

func TestFillFloat32WithInvalidValue(t *testing.T) {

	var s struct {
		Float32 float32 `foo:"bar"`
	}
	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"invalid"}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parsedErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
}

func TestFillFloat64(t *testing.T) {

	var s struct {
		Float64 float64 `foo:"bar"`
	}
	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"1.5"}, nil
			},
		},
	}

	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, float64(1.5), s.Float64)
}

func TestFillFloat64WithInvalidValue(t *testing.T) {

	var s struct {
		Float64 float64 `foo:"bar"`
	}
	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"invalid"}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parsedErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
}

func TestFillStruct(t *testing.T) {

	var s struct {
		Struct struct {
			Hello string `json:"hello"`
		} `foo:"bar"`
	}
	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{`{ "hello" : "world" }`}, nil
			},
		},
	}

	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, "world", s.Struct.Hello)
}

func TestFillStructWithInvalidJson(t *testing.T) {

	var s struct {
		Struct struct {
			Hello string `json:"hello"`
		} `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{`{ "hello" : invalidjson`}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parsedErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
}

func TestFillUnsupportedType(t *testing.T) {

	var s struct {
		Chan chan string `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{"helloworld"}, nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parsedErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
}

func TestFillIfSourceReturnsAnError(t *testing.T) {

	var s struct {
		Chan chan string `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) ([]string, error) {
				assert.Equal(t, "bar", field)
				return []string{}, errors.New("I am a test error")
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	parsedErr, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
	assert.Equal(t, "I am a test error", parsedErr.InnerError.Error())

}
