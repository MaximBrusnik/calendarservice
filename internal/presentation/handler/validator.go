package handler

import (
	"calendar/internal/domain"
	"net/http"
	"strconv"
	"time"
)

// RequestValidator содержит методы для валидации HTTP-запросов
type RequestValidator struct{}

// NewRequestValidator создает новый экземпляр валидатора
func NewRequestValidator() *RequestValidator {
	return &RequestValidator{}
}

// ParseAndValidateUserID парсит и валидирует user_id из формы или query
func (v *RequestValidator) ParseAndValidateUserID(value string) (int, error) {
	if value == "" {
		return 0, domain.NewValidationError("параметр user_id обязателен")
	}

	userID, err := strconv.Atoi(value)
	if err != nil {
		return 0, domain.NewValidationError("некорректный user_id")
	}

	if userID <= 0 {
		return 0, domain.NewValidationError("user_id должен быть положительным числом")
	}

	return userID, nil
}

// ParseAndValidateID парсит и валидирует id из формы
func (v *RequestValidator) ParseAndValidateID(value string) (int, error) {
	if value == "" {
		return 0, domain.NewValidationError("параметр id обязателен")
	}

	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, domain.NewValidationError("некорректный id")
	}

	if id <= 0 {
		return 0, domain.NewValidationError("id должен быть положительным числом")
	}

	return id, nil
}

// ParseAndValidateDate парсит и валидирует дату из строки
func (v *RequestValidator) ParseAndValidateDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, domain.NewValidationError("параметр date обязателен")
	}

	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, domain.NewValidationError("некорректный формат даты, используйте YYYY-MM-DD")
	}

	return date, nil
}

// ParseAndValidateYearMonth парсит и валидирует год и месяц из строки
func (v *RequestValidator) ParseAndValidateYearMonth(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, domain.NewValidationError("параметр date обязателен")
	}

	yearMonth, err := time.Parse("2006-01", value)
	if err != nil {
		return time.Time{}, domain.NewValidationError("некорректный формат даты, используйте YYYY-MM")
	}

	return yearMonth, nil
}

// ValidateRequiredFields проверяет, что все обязательные поля присутствуют
func (v *RequestValidator) ValidateRequiredFields(fields map[string]string) error {
	for name, value := range fields {
		if value == "" {
			return domain.NewValidationError("параметр " + name + " обязателен")
		}
	}
	return nil
}

// ParseFormAndValidate парсит форму и валидирует обязательные поля
func (v *RequestValidator) ParseFormAndValidate(r *http.Request, requiredFields []string) (map[string]string, error) {
	if err := r.ParseForm(); err != nil {
		return nil, domain.NewValidationError("ошибка парсинга формы")
	}

	fields := make(map[string]string)
	for _, field := range requiredFields {
		fields[field] = r.FormValue(field)
	}

	if err := v.ValidateRequiredFields(fields); err != nil {
		return nil, err
	}

	return fields, nil
}
