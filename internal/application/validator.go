package application

import (
	"calendar/internal/domain"
)

// ServiceValidator содержит методы для валидации в сервисном слое
type ServiceValidator struct{}

// NewServiceValidator создает новый валидатор сервиса
func NewServiceValidator() *ServiceValidator {
	return &ServiceValidator{}
}

// ValidateUserID проверяет корректность ID пользователя
func (v *ServiceValidator) ValidateUserID(userID int) error {
	if userID <= 0 {
		return domain.NewValidationError("некорректный ID пользователя")
	}
	return nil
}

// ValidateEventText проверяет корректность текста события
func (v *ServiceValidator) ValidateEventText(text string) error {
	if text == "" {
		return domain.NewValidationError("текст события не может быть пустым")
	}
	return nil
}

// ValidateEventID проверяет корректность ID события
func (v *ServiceValidator) ValidateEventID(id int) error {
	if id <= 0 {
		return domain.NewValidationError("некорректный ID события")
	}
	return nil
}

// ValidateEventData проверяет все данные события
func (v *ServiceValidator) ValidateEventData(userID int, text string) error {
	if err := v.ValidateUserID(userID); err != nil {
		return err
	}
	if err := v.ValidateEventText(text); err != nil {
		return err
	}
	return nil
}
