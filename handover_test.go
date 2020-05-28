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

func TestReturnValuerAsNilWithoutError(t *testing.T) {
	var s struct {
		String string `foo:"bar"`
	}
	s.String = "hello world"

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return nil, nil
			},
		},
	}
	err := From(sources).To(&s)
	assert.NoError(t, err)

	assert.Equal(t, "hello world", s.String)
}

func TestReturnValuerAsNilWithError(t *testing.T) {
	var s struct {
		String string `foo:"bar"`
	}
	s.String = "hello world"

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return nil, errors.New("test error")
			},
		},
	}
	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error
	assert.True(t, errors.As(err, &parsedErr))

	assert.Equal(t, "", parsedErr.Value)
	assert.Equal(t, "hello world", s.String)
}

func TestFillWithNilStruct(t *testing.T) {

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("helloworld"), nil
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
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("helloworld"), nil
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
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Values([]string{"hello", "world"}), nil
			},
		},
		{
			Tag: "john",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "doe", field)
				return Value(`{ "some": "json" }`), nil
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
	s.Slice = []int{1}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Values([]string{"invalid", "value"}), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error
	assert.True(t, errors.As(err, &parsedErr))

	assert.Equal(t, "bar", parsedErr.Field)
	assert.Equal(t, "invalid", parsedErr.Value)
	assert.Error(t, parsedErr.InnerError)

	assert.Equal(t, 1, s.Slice[0])
}

func TestFillString(t *testing.T) {

	var s struct {
		String string `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("helloworld"), nil
			},
		},
	}
	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, "helloworld", s.String)
}

func TestFillTimeDuration(t *testing.T) {

	var s struct {
		Duration time.Duration `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("1h"), nil
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
	s.Duration = time.Second

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("1"), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Equal(t, "1", parsedErr.Value)
	assert.Error(t, parsedErr.InnerError)

	assert.Equal(t, time.Second, s.Duration)
}

func TestFillInt(t *testing.T) {

	var s struct {
		Int int `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("1"), nil
			},
		},
	}
	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, int(1), s.Int)
}

func TestFillIntWithInvalidValue(t *testing.T) {

	var s struct {
		Int int `foo:"bar"`
	}
	s.Int = 1

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("invalid"), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Equal(t, "invalid", parsedErr.Value)
	assert.Error(t, parsedErr.InnerError)

	assert.Equal(t, 1, s.Int)
}

func TestFillInt8(t *testing.T) {

	var s struct {
		Int8 int8 `foo:"bar"`
	}
	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("1"), nil
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
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("invalid"), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Equal(t, "invalid", parsedErr.Value)
	assert.Error(t, parsedErr.InnerError)
}

func TestFillInt16(t *testing.T) {

	var s struct {
		Int16 int16 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("1"), nil
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
	s.Int16 = int16(1)

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("invalid"), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Equal(t, "invalid", parsedErr.Value)
	assert.Error(t, parsedErr.InnerError)

	assert.Equal(t, int16(1), s.Int16)
}

func TestFillInt32(t *testing.T) {

	var s struct {
		Int32 int32 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("1"), nil
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
	s.Int32 = int32(1)

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("invalid"), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Equal(t, "invalid", parsedErr.Value)
	assert.Error(t, parsedErr.InnerError)

	assert.Equal(t, int32(1), s.Int32)
}

func TestFillInt64(t *testing.T) {

	var s struct {
		Int64 int64 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("1"), nil
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
	s.Int64 = int64(1)

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("invalid"), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Equal(t, "invalid", parsedErr.Value)
	assert.Error(t, parsedErr.InnerError)

	assert.Equal(t, int64(1), s.Int64)
}

func TestFillUInt(t *testing.T) {

	var s struct {
		UInt uint `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("1"), nil
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
	s.UInt = uint(1)

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("invalid"), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Equal(t, "invalid", parsedErr.Value)
	assert.Error(t, parsedErr.InnerError)

	assert.Equal(t, uint(1), s.UInt)
}

func TestFillUInt8(t *testing.T) {

	var s struct {
		UInt8 uint8 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("1"), nil
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
	s.UInt8 = uint8(1)

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("invalid"), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Equal(t, "invalid", parsedErr.Value)
	assert.Error(t, parsedErr.InnerError)

	assert.Equal(t, uint8(1), s.UInt8)
}

func TestFillUInt16(t *testing.T) {

	var s struct {
		UInt16 uint16 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("1"), nil
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
	s.UInt16 = uint16(1)

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("invalid"), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Equal(t, "invalid", parsedErr.Value)
	assert.Error(t, parsedErr.InnerError)

	assert.Equal(t, uint16(1), s.UInt16)
}

func TestFillUInt32(t *testing.T) {

	var s struct {
		UInt32 uint32 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("1"), nil
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
	s.UInt32 = uint32(1)

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("invalid"), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Equal(t, "invalid", parsedErr.Value)
	assert.Error(t, parsedErr.InnerError)

	assert.Equal(t, uint32(1), s.UInt32)
}

func TestFillUInt64(t *testing.T) {

	var s struct {
		UInt64 uint64 `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("1"), nil
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
	s.UInt64 = uint64(1)

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("invalid"), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Equal(t, "invalid", parsedErr.Value)
	assert.Error(t, parsedErr.InnerError)

	assert.Equal(t, uint64(1), s.UInt64)
}

func TestFillBool(t *testing.T) {

	var s struct {
		Bool bool `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("true"), nil
			},
		},
	}

	assert.NoError(t, From(sources).To(&s))
	assert.Equal(t, true, s.Bool)
}

func TestFillBoolWithInvalidValue(t *testing.T) {

	var s struct {
		Bool bool `foo:"bar"`
	}
	s.Bool = true

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("invalid"), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Equal(t, "invalid", parsedErr.Value)
	assert.Error(t, parsedErr.InnerError)

	assert.True(t, s.Bool)
}

func TestFillFloat32(t *testing.T) {

	var s struct {
		Float32 float32 `foo:"bar"`
	}
	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("1.5"), nil
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
	s.Float32 = float32(1.5)

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("invalid"), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Equal(t, "invalid", parsedErr.Value)
	assert.Error(t, parsedErr.InnerError)

	assert.Equal(t, float32(1.5), s.Float32)
}

func TestFillFloat64(t *testing.T) {

	var s struct {
		Float64 float64 `foo:"bar"`
	}
	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("1.5"), nil
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
	s.Float64 = float64(1.5)

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("invalid"), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Equal(t, "invalid", parsedErr.Value)
	assert.Error(t, parsedErr.InnerError)

	assert.Equal(t, float64(1.5), s.Float64)
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
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value(`{ "hello" : "world" }`), nil
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
	s.Struct.Hello = "world"

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value(`{ "hello" : invalidjson`), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Equal(t, `{ "hello" : invalidjson`, parsedErr.Value)
	assert.Error(t, parsedErr.InnerError)

	assert.Equal(t, "world", s.Struct.Hello)
}

func TestFillUnsupportedType(t *testing.T) {

	var s struct {
		Chan chan string `foo:"bar"`
	}

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value("helloworld"), nil
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
}

func TestFillIfSourceReturnsAnError(t *testing.T) {

	var s struct {
		String string `foo:"bar"`
	}
	s.String = "hello world"

	sources := []Source{
		{
			Tag: "foo",
			Get: func(field string) (Valuer, error) {
				assert.Equal(t, "bar", field)
				return Value(""), errors.New("I am a test error")
			},
		},
	}

	err := From(sources).To(&s)
	assert.Error(t, err)

	var parsedErr Error

	assert.True(t, errors.As(err, &parsedErr))
	assert.Equal(t, "bar", parsedErr.Field)
	assert.Error(t, parsedErr.InnerError)
	assert.Equal(t, "I am a test error", parsedErr.InnerError.Error())

	assert.Equal(t, "hello world", s.String)
}
