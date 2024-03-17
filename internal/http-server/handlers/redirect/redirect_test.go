package redirect

import (
	"log/slog"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		log       *slog.Logger
		urlGetter URLGetter
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

		})
	}
}
