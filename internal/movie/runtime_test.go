package movie

import (
	"errors"
	"testing"

	"github.com/0vkanix/greenlight/internal/assert"
)

func TestRuntimeMarshalJSON(t *testing.T) {
	r := Runtime(102)

	got, err := r.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	want := "\"102 mins\""
	assert.Equal(t, string(got), want)
}

func TestRuntimeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      Runtime
		wantErr   error
	}{
		{
			name:    "Valid input",
			input:   "\"102 mins\"",
			want:    Runtime(102),
			wantErr: nil,
		},
		{
			name:    "Invalid format (no mins)",
			input:   "\"102\"",
			wantErr: ErrInvalidRuntimeFormat,
		},
		{
			name:    "Invalid format (wrong suffix)",
			input:   "\"102 minutes\"",
			wantErr: ErrInvalidRuntimeFormat,
		},
		{
			name:    "Invalid format (not quoted)",
			input:   "102 mins",
			wantErr: ErrInvalidRuntimeFormat,
		},
		{
			name:    "Invalid format (non-numeric)",
			input:   "\"abc mins\"",
			wantErr: ErrInvalidRuntimeFormat,
		},
		{
			name:    "Empty string",
			input:   "\"\"",
			wantErr: ErrInvalidRuntimeFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r Runtime
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
