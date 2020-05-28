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
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Error struct {
	Field      string
	Source     string
	Value      string
	InnerError error
}

func newError(field, source string, values []string, err error) Error {

	e := Error{
		Field:      field,
		Source:     source,
		InnerError: err,
	}

	switch ie := e.InnerError.(type) {
	case *strconv.NumError:
		e.Value = ie.Num
	case *time.ParseError:
		e.Value = ie.Value
	case *json.UnsupportedValueError:
		e.Value = ie.Str
	default:
		if len(values) <= 0 {
			return e
		}
		if len(values) == 1 {
			e.Value = values[0]
			return e
		}
		e.Value = "[" + strings.Join(values, " ") + "]"
	}

	return e
}

func (te Error) Error() string {
	return fmt.Sprintf("failed to set field %q from source %q: %s", te.Field, te.Source, te.InnerError)
}
