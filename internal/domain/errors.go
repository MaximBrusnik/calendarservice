package domain

import (
	"net/http"
)

// AppError представляет ошибку приложения с HTTP статус-кодом
type AppError struct {
	Message    string
	StatusCode int
	Err        error
}

// Error возвращает сообщение об ошибке
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// Unwrap возвращает вложенную ошибку
func (e *AppError) Unwrap() error {
	return e.Err
}

// GetStatusCode возвращает HTTP статус-код для ошибки
func (e *AppError) GetStatusCode() int {
	return e.StatusCode
}

// NewAppError создает новую ошибку приложения
func NewAppError(message string, statusCode int, err error) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// HTTP статус-коды для ошибок
const (
	StatusBadRequest          = http.StatusBadRequest          // 400
	StatusUnauthorized        = http.StatusUnauthorized        // 401
	StatusForbidden           = http.StatusForbidden           // 403
	StatusNotFound            = http.StatusNotFound            // 404
	StatusInternalServerError = http.StatusInternalServerError // 500
	StatusServiceUnavailable  = http.StatusServiceUnavailable  // 503
)

// Функции для создания типизированных ошибок
func NewValidationError(message string) *AppError {
	return NewAppError(message, StatusBadRequest, nil)
}

func NewBusinessLogicError(message string) *AppError {
	return NewAppError(message, StatusServiceUnavailable, nil)
}

func NewNotFoundError(message string) *AppError {
	return NewAppError(message, StatusNotFound, nil)
}

func NewAccessDeniedError(message string) *AppError {
	return NewAppError(message, StatusForbidden, nil)
}

func NewInternalError(message string, err error) *AppError {
	return NewAppError(message, StatusInternalServerError, err)
}
