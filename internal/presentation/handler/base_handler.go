package handler

import (
	"calendar/internal/domain"
	"encoding/json"
	"net/http"
)

// BaseHandler содержит общие методы для всех обработчиков
type BaseHandler struct {
	validator *RequestValidator
}

// NewBaseHandler создает новый базовый обработчик
func NewBaseHandler() *BaseHandler {
	return &BaseHandler{
		validator: NewRequestValidator(),
	}
}

// writeResponse записывает JSON-ответ
func (h *BaseHandler) writeResponse(w http.ResponseWriter, statusCode int, response domain.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// writeError записывает ошибку
func (h *BaseHandler) writeError(w http.ResponseWriter, statusCode int, message string) {
	h.writeResponse(w, statusCode, domain.Response{Error: message})
}

// writeSuccess записывает успешный ответ
func (h *BaseHandler) writeSuccess(w http.ResponseWriter, data interface{}) {
	h.writeResponse(w, http.StatusOK, domain.Response{Result: data})
}

// handleError обрабатывает ошибки и возвращает соответствующий HTTP статус-код
func (h *BaseHandler) handleError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*domain.AppError); ok {
		// Это наша типизированная ошибка
		h.writeError(w, appErr.GetStatusCode(), appErr.Error())
		return
	}

	// Неизвестная ошибка - возвращаем 500
	h.writeError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
}

// GetValidator возвращает валидатор запросов
func (h *BaseHandler) GetValidator() *RequestValidator {
	return h.validator
}
