package tests

import (
	"net/http"
	"net/url"
	"path"
	"testing"

	"github.com/andrei1998Front/url-shortener/internal/lib/api"
	"github.com/andrei1998Front/url-shortener/internal/lib/random"
	"github.com/stretchr/testify/require"

	dl "github.com/andrei1998Front/url-shortener/internal/http-server/handlers/url/delete"
	"github.com/andrei1998Front/url-shortener/internal/http-server/handlers/url/save"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:8082"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	alias, _ := random.NewRandomString(10)

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: alias,
		}).
		WithBasicAuth("myuser", "mypass").
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("Alias")
}

//nolint:funlen
func TestURLShortener_SaveRedirectDelete(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: gofakeit.Word(),
			error: "field URL is not a valid URL",
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// Save

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("myuser", "mypass").
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("Alias")

				resp.Value("Error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("Alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("Alias").String().NotEmpty()

				alias = resp.Value("Alias").String().Raw()
			}

			// Redirect

			testRedirect(t, alias, tc.url)

			// Delete

			reqDel := e.DELETE("/"+path.Join("url")).
				WithJSON(dl.Request{
					Alias: tc.alias,
				}).
				WithBasicAuth("myuser", "mypass").
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.alias != "" {
				reqDel.Value("Status").String().IsEqual("OK")
				testRedirectNotFound(t, alias)
			} else {
				reqDel.Value("Status").String().IsEqual("Error")
			}
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}

func testRedirectNotFound(t *testing.T, alias string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	_, err := api.GetRedirect(u.String())

	require.ErrorIs(t, err, api.ErrInvalidStatusCode)
}
