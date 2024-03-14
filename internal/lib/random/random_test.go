package random

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRandomString(t *testing.T) {

	tests := []struct {
		size      int
		wantLen   int
		wantError bool
		err       string
	}{
		{
			size:      2,
			wantLen:   2,
			wantError: false,
			err:       "",
		},
		{
			size:      1,
			wantLen:   1,
			wantError: false,
			err:       "",
		},
		{
			size:      0,
			wantLen:   0,
			wantError: false,
			err:       "",
		},
		{
			size:      -1,
			wantLen:   0,
			wantError: true,
			err:       ErrNegativeSize.Error(),
		},
	}
	for _, tt := range tests {
		result, err := NewRandomString(tt.size)

		if (err != nil) != tt.wantError {
			t.Errorf("NewRandomString(%d) error = %s, wantError = %t", tt.size, err.Error(), tt.wantError)
			return
		} else if err != nil {
			require.EqualError(t, err, tt.err)
		}

		require.Equal(t, tt.wantLen, len(result))
	}
}
