package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/leary1337/url-shortener/internal/app/service"
	"github.com/leary1337/url-shortener/pkg/logger"
)

// MockPingService is a mock implementation of the Ping service interface.
type MockPingService struct {
	mock.Mock
}

func (m *MockPingService) PingDB(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestPingHandler_PingDB(t *testing.T) {
	mockService := new(MockPingService)
	l := logger.New("info", os.Stdout)
	handler := NewPingHandler(l, mockService)

	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	tests := []struct {
		name           string
		mockErr        error
		expectedStatus int
	}{
		{
			name:           "ping successful",
			mockErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ping failure",
			mockErr:        service.ErrPostgresPingFailed,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock before each test
			mockService.ExpectedCalls = nil
			mockService.Calls = nil

			mockService.On("PingDB", mock.Anything).Return(tt.mockErr)

			req := httptest.NewRequest("GET", "/ping", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockService.AssertCalled(t, "PingDB", mock.Anything)
			mockService.AssertExpectations(t)
		})
	}
}
