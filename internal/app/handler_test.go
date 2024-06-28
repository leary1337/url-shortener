package app

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServerHandler_MainPage(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
		location    string
	}
	tests := []struct {
		name       string
		method     string
		requestURI string
		body       string
		want       want
	}{
		{
			"valid GET test #1",
			http.MethodGet,
			"/shortTestUrl1",
			"",
			want{
				code:     http.StatusTemporaryRedirect,
				location: "https://google.com",
			},
		},
		{
			"not found GET test #2",
			http.MethodGet,
			"/shortTestUrl2",
			"",
			want{
				code: http.StatusNotFound,
			},
		},
		{
			"valid POST test #3",
			http.MethodPost,
			"/",
			"https://google.com",
			want{
				code: http.StatusCreated,
			},
		},
		{
			"empty body POST test #4",
			http.MethodPost,
			"/",
			"",
			want{
				code: http.StatusBadRequest,
			},
		},
	}
	serverHandler := &ServerHandler{
		serverAddr: "localhost:8080",
		urlMap:     map[string]string{"shortTestUrl1": "https://google.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.requestURI, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			serverHandler.MainPage(w, r)

			result := w.Result()

			if tt.method == http.MethodGet {
				assert.Equal(t, tt.want.location, result.Header.Get("Location"))
				assert.Equal(t, tt.want.code, result.StatusCode)
			} else {
				assert.Equal(t, tt.want.code, result.StatusCode)
				assert.NotEmpty(t, result.Body)
			}
		})
	}
}
