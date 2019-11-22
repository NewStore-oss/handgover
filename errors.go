package handgover

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func FromError(err error) (Error, bool) {
	tErr, ok := err.(Error)
	return tErr, ok
}

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
