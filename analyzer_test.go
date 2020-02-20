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
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	A string
	B int      `b:"b"`
	C []string `c:"c c c"`
	D struct{ DD string }
}

func TestAnalyzeWithStruct(t *testing.T) {
	result := analyze(testStruct{})

	assert.Len(t, result.Properties, 4)

	assert.Equal(t, "A", result.Properties[0].Name)
	assert.Equal(t, 0, result.Properties[0].Index)
	assert.Equal(t, reflect.String, result.Properties[0].Kind)
	assert.EqualValues(t, "", result.Properties[0].Tag)

	assert.Equal(t, "B", result.Properties[1].Name)
	assert.Equal(t, 1, result.Properties[1].Index)
	assert.Equal(t, reflect.Int, result.Properties[1].Kind)
	assert.EqualValues(t, "b:\"b\"", result.Properties[1].Tag)

	assert.Equal(t, "C", result.Properties[2].Name)
	assert.Equal(t, 2, result.Properties[2].Index)
	assert.Equal(t, reflect.Slice, result.Properties[2].Kind)
	assert.EqualValues(t, "c:\"c c c\"", result.Properties[2].Tag)

	assert.Equal(t, "D", result.Properties[3].Name)
	assert.Equal(t, 3, result.Properties[3].Index)
	assert.Equal(t, reflect.Struct, result.Properties[3].Kind)
	assert.EqualValues(t, "", result.Properties[3].Tag)
}

func TestAnalyzeWithStructAsPointer(t *testing.T) {
	result := analyze(&testStruct{})
	assert.Len(t, result.Properties, 4)

	assert.Equal(t, "A", result.Properties[0].Name)
	assert.Equal(t, 0, result.Properties[0].Index)
	assert.Equal(t, reflect.String, result.Properties[0].Kind)
	assert.EqualValues(t, "", result.Properties[0].Tag)

	assert.Equal(t, "B", result.Properties[1].Name)
	assert.Equal(t, 1, result.Properties[1].Index)
	assert.Equal(t, reflect.Int, result.Properties[1].Kind)
	assert.EqualValues(t, "b:\"b\"", result.Properties[1].Tag)

	assert.Equal(t, "C", result.Properties[2].Name)
	assert.Equal(t, 2, result.Properties[2].Index)
	assert.Equal(t, reflect.Slice, result.Properties[2].Kind)
	assert.EqualValues(t, "c:\"c c c\"", result.Properties[2].Tag)

	assert.Equal(t, "D", result.Properties[3].Name)
	assert.Equal(t, 3, result.Properties[3].Index)
	assert.Equal(t, reflect.Struct, result.Properties[3].Kind)
	assert.EqualValues(t, "", result.Properties[3].Tag)
}
