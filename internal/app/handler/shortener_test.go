package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/leary1337/url-shortener/internal/app/entity"
	"github.com/leary1337/url-shortener/internal/app/middleware"
	"github.com/leary1337/url-shortener/internal/app/service"
	"github.com/leary1337/url-shortener/pkg/logger"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) ShortenURL(ctx context.Context, originalURL string) (*entity.ShortURL, error) {
	args := m.Called(originalURL)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.ShortURL), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockService) ShortenBatch(ctx context.Context, sb []entity.ShortenBatchRequestBody) ([]entity.ShortenBatchResponseBody, error) {
	args := m.Called(ctx, sb)
	return args.Get(0).([]entity.ShortenBatchResponseBody), args.Error(1)
}

func (m *MockService) ResolveURL(ctx context.Context, shortURL string) (*entity.ShortURL, error) {
	args := m.Called(shortURL)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.ShortURL), args.Error(1)
	}
	return nil, args.Error(1)
}

var l = logger.New("info", os.Stdout)

func TestHandler_ShortenURL(t *testing.T) {
	mockService := new(MockService)
	handler := NewShortenerHandler(l, mockService)

	router := chi.NewRouter()
	router.Use(middleware.CompressMiddleware)
	router.Post("/", handler.ShortenURL)

	tests := []struct {
		name           string
		body           string
		mockShortURL   *entity.ShortURL
		mockErr        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful shorten URL",
			body:           "https://example.com",
			mockShortURL:   &entity.ShortURL{ShortURL: "http://localhost:8080/abcd1234", OriginalURL: "https://example.com"},
			expectedStatus: http.StatusCreated,
			expectedBody:   "http://localhost:8080/abcd1234",
		},
		{
			name:           "empty body",
			body:           "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.body != "" {
				mockService.On("ShortenURL", tt.body).Return(tt.mockShortURL, tt.mockErr)
			}

			req := httptest.NewRequest("POST", "/", strings.NewReader(tt.body))
			req.Header.Set("Accept-Encoding", "gzip")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, tt.expectedBody, rec.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_ResolveURL(t *testing.T) {
	mockService := new(MockService)
	handler := NewShortenerHandler(l, mockService)

	router := chi.NewRouter()
	router.Use(middleware.CompressMiddleware)
	router.Get("/{shortURL}", handler.ResolveURL)

	tests := []struct {
		name             string
		shortURL         string
		mockShortURL     *entity.ShortURL
		mockErr          error
		expectedStatus   int
		expectedLocation string
	}{
		{
			name:             "successful resolve URL",
			shortURL:         "abcd1234",
			mockShortURL:     &entity.ShortURL{ShortURL: "http://localhost:8080/abcd1234", OriginalURL: "https://example.com"},
			expectedStatus:   http.StatusTemporaryRedirect,
			expectedLocation: "https://example.com",
		},
		{
			name:             "URL not found",
			shortURL:         "nonexistent",
			mockErr:          service.ErrURLNotFound,
			expectedStatus:   http.StatusNotFound,
			expectedLocation: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.On("ResolveURL", tt.shortURL).Return(tt.mockShortURL, tt.mockErr)

			req := httptest.NewRequest("GET", fmt.Sprintf("/%s", tt.shortURL), nil)
			req.Header.Set("Accept-Encoding", "gzip")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedLocation != "" {
				assert.Equal(t, tt.expectedLocation, rec.Header().Get("Location"))
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestShortenerHandler_ShortenURLJSON(t *testing.T) {
	mockService := new(MockService)
	handler := NewShortenerHandler(l, mockService)

	router := chi.NewRouter()
	router.Use(middleware.CompressMiddleware)
	router.Post("/api/shorten", handler.ShortenURLJSON)

	tests := []struct {
		name           string
		body           string
		mockShortURL   *entity.ShortURL
		mockErr        error
		expectedStatus int
		expectedBody   string
		isJSONBody     bool
	}{
		{
			name:           "successful shorten URL",
			body:           `{"url":"https://example.com"}`,
			mockShortURL:   &entity.ShortURL{ShortURL: "http://localhost:8080/abcd1234", OriginalURL: "https://example.com"},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"result":"http://localhost:8080/abcd1234"}`,
			isJSONBody:     true,
		},
		{
			name:           "empty body",
			body:           ``,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ``,
			isJSONBody:     false,
		},
		{
			name:           "invalid JSON body",
			body:           `{invalid}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ``,
			isJSONBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.body != "" && tt.expectedStatus == http.StatusCreated {
				mockService.On("ShortenURL", "https://example.com").Return(tt.mockShortURL, tt.mockErr)
			}

			req := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.isJSONBody {
				assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			} else {
				assert.Equal(t, tt.expectedBody, rec.Body.String())
			}

			if tt.expectedStatus == http.StatusCreated {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestShortenerHandler_ShortenURLBatch(t *testing.T) {
	mockService := new(MockService)
	l := logger.New("info", os.Stdout)
	handler := NewShortenerHandler(l, mockService)

	router := chi.NewRouter()
	router.Use(middleware.CompressMiddleware)
	router.Post("/api/shorten/batch", handler.ShortenURLBatch)

	tests := []struct {
		name           string
		body           string
		mockResp       []entity.ShortenBatchResponseBody
		mockErr        error
		expectedStatus int
		expectedBody   string
		isJSONBody     bool
	}{
		{
			name: "successful shorten URLs",
			body: `[
				{
					"correlation_id": "1",
					"original_url": "https://example1.com"
				},
				{
					"correlation_id": "2",
					"original_url": "https://example2.com"
				},
				{
					"correlation_id": "3",
					"original_url": "https://example3.com"
				}
			]`,
			mockResp: []entity.ShortenBatchResponseBody{
				{Id: "1", ShortURL: "http://localhost:8080/abcd1"},
				{Id: "2", ShortURL: "http://localhost:8080/abcd2"},
				{Id: "3", ShortURL: "http://localhost:8080/abcd3"},
			},
			expectedStatus: http.StatusCreated,
			expectedBody: `[
				{
					"correlation_id": "1",
					"short_url": "http://localhost:8080/abcd1"
				},
				{
					"correlation_id": "2",
					"short_url": "http://localhost:8080/abcd2"
				},
				{
					"correlation_id": "3",
					"short_url": "http://localhost:8080/abcd3"
				}
			]`,
			isJSONBody: true,
		},
		{
			name:           "empty body",
			body:           ``,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ``,
			isJSONBody:     false,
		},
		{
			name:           "invalid JSON body",
			body:           `{invalid}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ``,
			isJSONBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.body != "" && tt.expectedStatus == http.StatusCreated {
				var reqBodies []entity.ShortenBatchRequestBody
				err := json.Unmarshal([]byte(tt.body), &reqBodies)
				assert.NoError(t, err)
				mockService.On("ShortenBatch", mock.Anything, reqBodies).Return(tt.mockResp, tt.mockErr)
			}

			req := httptest.NewRequest("POST", "/api/shorten/batch", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.isJSONBody {
				assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			} else {
				assert.Equal(t, tt.expectedBody, rec.Body.String())
			}

			if tt.expectedStatus == http.StatusCreated {
				mockService.AssertExpectations(t)
			}
		})
	}
}
