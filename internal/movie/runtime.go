package movie

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrInvalidRuntimeFormat is returned when a JSON runtime string is 
// not in the expected "<runtime> mins" format.
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

// Runtime represents a movie's runtime in minutes. It is a custom type 
// to allow for custom JSON marshaling and unmarshaling.
type Runtime int32

// MarshalJSON returns the runtime as a JSON string in the format "<runtime> mins".
func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

// UnmarshalJSON parses a JSON string in the format "<runtime> mins" and 
// stores the result in the Runtime object.
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(i)

	return nil
}
