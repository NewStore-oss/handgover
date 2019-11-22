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
