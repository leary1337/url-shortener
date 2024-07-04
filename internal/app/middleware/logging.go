package middleware

import (
	"net/http"
	"time"

	"github.com/leary1337/url-shortener/pkg/logger"
)

// responseWriterWrapper обертка для http.ResponseWriter, чтобы перехватывать статус ответа
type responseWriterWrapper struct {
	http.ResponseWriter
	status        int
	contentLength int
}

func (ww *responseWriterWrapper) Write(b []byte) (int, error) {
	if ww.status == 0 {
		ww.status = http.StatusOK
	}
	n, err := ww.ResponseWriter.Write(b)
	ww.contentLength += n
	return n, err
}

func (ww *responseWriterWrapper) WriteHeader(status int) {
	ww.status = status
	ww.ResponseWriter.WriteHeader(status)
}

func LoggingMiddleware(l logger.Interface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Создаем обертку для ResponseWriter, чтобы перехватывать статус ответа
			ww := &responseWriterWrapper{ResponseWriter: w}
			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			// Логируем информацию о запросе и ответе
			l.Info("method=%s uri=%s status=%d duration=%s content_length=%d", r.Method, r.RequestURI, ww.status, duration, ww.contentLength)
		})
	}
}
