package types_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/0vkanix/greenlight/internal/assert"
	"github.com/0vkanix/greenlight/internal/movie/types"
)

// ExampleRuntime_MarshalJSON demonstrates how the Runtime type
// marshals into the custom JSON string format.
func ExampleRuntime_MarshalJSON() {
	r := types.Runtime(102)
	js, _ := r.MarshalJSON()
	fmt.Println(string(js))
	// Output: "102 mins"
}

// TestRuntimeMarshalJSON verifies that the Runtime type correctly
// marshals to the expected JSON string format.
func TestRuntimeMarshalJSON(t *testing.T) {
	r := types.Runtime(102)

	got, err := r.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	want := "\"102 mins\""
	assert.Equal(t, string(got), want)
}

// TestRuntimeUnmarshalJSON verifies that the Runtime type correctly
// parses various valid and invalid JSON string formats.
func TestRuntimeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    types.Runtime
		wantErr error
	}{
		{
			name:    "Valid input",
			input:   "\"102 mins\"",
			want:    types.Runtime(102),
			wantErr: nil,
		},
		{
			name:    "Invalid format (no mins)",
			input:   "\"102\"",
			wantErr: types.ErrInvalidRuntimeFormat,
		},
		{
			name:    "Invalid format (wrong suffix)",
			input:   "\"102 minutes\"",
			wantErr: types.ErrInvalidRuntimeFormat,
		},
		{
			name:    "Invalid format (not quoted)",
			input:   "102 mins",
			wantErr: types.ErrInvalidRuntimeFormat,
		},
		{
			name:    "Invalid format (non-numeric)",
			input:   "\"abc mins\"",
			wantErr: types.ErrInvalidRuntimeFormat,
		},
		{
			name:    "Empty string",
			input:   "\"\"",
			wantErr: types.ErrInvalidRuntimeFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r types.Runtime
			err := r.UnmarshalJSON([]byte(tt.input))

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("got error %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, r, tt.want)
		})
	}
}
