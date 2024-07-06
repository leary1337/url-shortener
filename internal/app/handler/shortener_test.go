package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/leary1337/url-shortener/internal/app/service"
	"github.com/leary1337/url-shortener/pkg/logger"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) ShortenURL(originalURL string) string {
	args := m.Called(originalURL)
	return args.String(0)
}

func (m *MockService) ResolveURL(shortURL string) (string, error) {
	args := m.Called(shortURL)
	return args.String(0), args.Error(1)
}

var l = logger.New("info")

func TestHandler_ShortenURL(t *testing.T) {
	mockService := new(MockService)
	redirectAddr := "http://localhost:8080"
	handler := NewShortenerHandler(l, mockService, redirectAddr)

	router := chi.NewRouter()
	router.Post("/", handler.ShortenURL)

	tests := []struct {
		name           string
		body           string
		mockShortURL   string
		mockErr        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful shorten URL",
			body:           "https://example.com",
			mockShortURL:   "abcd1234",
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
				mockService.On("ShortenURL", tt.body).Return(tt.mockShortURL)
			}

			req := httptest.NewRequest("POST", "/", strings.NewReader(tt.body))
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
	redirectAddr := "http://localhost:8080"
	handler := NewShortenerHandler(l, mockService, redirectAddr)

	router := chi.NewRouter()
	router.Get("/{shortURL}", handler.ResolveURL)

	tests := []struct {
		name             string
		shortURL         string
		mockOriginalURL  string
		mockErr          error
		expectedStatus   int
		expectedLocation string
	}{
		{
			name:             "successful resolve URL",
			shortURL:         "abcd1234",
			mockOriginalURL:  "https://example.com",
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
			mockService.On("ResolveURL", tt.shortURL).Return(tt.mockOriginalURL, tt.mockErr)

			req := httptest.NewRequest("GET", fmt.Sprintf("/%s", tt.shortURL), nil)
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
	redirectAddr := "http://localhost:8080"
	handler := NewShortenerHandler(l, mockService, redirectAddr)

	router := chi.NewRouter()
	router.Post("/api/shorten", handler.ShortenURLJSON)

	tests := []struct {
		name           string
		body           string
		mockShortURL   string
		mockErr        error
		expectedStatus int
		expectedBody   string
		isJSONBody     bool
	}{
		{
			name:           "successful shorten URL",
			body:           `{"url":"https://example.com"}`,
			mockShortURL:   "abcd1234",
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
				mockService.On("ShortenURL", "https://example.com").Return(tt.mockShortURL)
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
