package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// LoggingMiddleware логирует каждый HTTP-запрос
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Логируем начало запроса
		fmt.Printf("[%s] %s %s - Начало запроса\n",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path)

		// Создаем wrapper для ResponseWriter для перехвата статус-кода
		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Выполняем следующий обработчик
		next.ServeHTTP(wrappedWriter, r)

		// Логируем завершение запроса
		duration := time.Since(start)
		fmt.Printf("[%s] %s %s - Завершение запроса (статус: %d, время: %v)\n",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			wrappedWriter.statusCode,
			duration)
	})
}

// responseWriter оборачивает http.ResponseWriter для перехвата статус-кода
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader перехватывает статус-код
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
