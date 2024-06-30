package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/leary1337/url-shortener/internal/app/config"
)

func TestServerHandler_GenerateShortURL(t *testing.T) {
	type want struct {
		code         int
		response     string
		redirectAddr string
	}
	tests := []struct {
		name   string
		method string
		body   string
		want   want
	}{
		{
			"valid POST test #1",
			http.MethodPost,
			"https://google.com",
			want{
				code:         http.StatusCreated,
				redirectAddr: "http://testURL",
			},
		},
		{
			"empty body POST test #2",
			http.MethodPost,
			"",
			want{
				code: http.StatusBadRequest,
			},
		},
	}
	serverHandler := &ServerHandler{
		cfg: &config.Config{
			Addr:         "localhost:8080",
			RedirectAddr: "http://testURL",
		},
		urlMap: map[string]string{"shortTestUrl1": "https://google.com"},
	}
	ts := httptest.NewServer(ShortenerRouter(serverHandler))
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, body := testRequest(t, ts, tt.method, `/`, strings.NewReader(tt.body))

			require.Equal(t, tt.want.code, response.StatusCode)
			if tt.want.code == http.StatusCreated {
				assert.NotEmpty(t, body)
			}
			if tt.want.redirectAddr != "" {
				assert.Contains(t, body, tt.want.redirectAddr)
			}
		})
	}
}

func TestServerHandler_GetOriginalURL(t *testing.T) {
	type want struct {
		code     int
		location string
	}
	tests := []struct {
		name       string
		method     string
		requestURI string
		want       want
	}{
		{
			"valid GET test #1",
			http.MethodGet,
			"/shortTestUrl1",
			want{
				code:     http.StatusTemporaryRedirect,
				location: "https://google.com",
			},
		},
		{
			"not found GET test #2",
			http.MethodGet,
			"/shortTestUrl2",
			want{
				code: http.StatusNotFound,
			},
		},
	}
	serverHandler := &ServerHandler{
		cfg: &config.Config{
			Addr:         "localhost:8080",
			RedirectAddr: "localhost:8080",
		},
		urlMap: map[string]string{"shortTestUrl1": "https://google.com"},
	}

	ts := httptest.NewServer(ShortenerRouter(serverHandler))
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, _ := testRequest(t, ts, tt.method, tt.requestURI, nil)

			assert.Equal(t, tt.want.location, response.Header.Get("Location"))
			assert.Equal(t, tt.want.code, response.StatusCode)
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
