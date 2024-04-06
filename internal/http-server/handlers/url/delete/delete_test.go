package delete

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrei1998Front/url-shortener/internal/http-server/handlers/url/delete/mocks"
	"github.com/andrei1998Front/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/andrei1998Front/url-shortener/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
		},
		{
			name:      "Empty alias",
			alias:     "",
			respError: "invalid alias",
		},
		{
			name:      "DeleteURL Error",
			alias:     "test_alias",
			respError: "there are no url with this alias",
			mockError: storage.ErrUrlNotFound,
		},
		{
			name:      "DeleteURL Error",
			alias:     "test_alias",
			respError: "failed to delete alias",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", tc.alias).
					Return(tc.mockError).
					Once()
			}

			handler := New(slogdiscard.NewDiscardLogger(), urlDeleterMock)

			input := fmt.Sprintf(`{"alias": "%s"}`, tc.alias)

			req, err := http.NewRequest(http.MethodDelete, "/url", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
