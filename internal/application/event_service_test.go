package application

import (
	"calendar/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventRepository - мок для EventRepository
type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) Create(event *domain.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventRepository) Update(event *domain.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventRepository) Delete(id int, userID int) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockEventRepository) GetByID(id int) (*domain.Event, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Event), args.Error(1)
}

func (m *MockEventRepository) GetByUserAndDate(userID int, date time.Time) ([]*domain.Event, error) {
	args := m.Called(userID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Event), args.Error(1)
}

func (m *MockEventRepository) GetByUserAndDateRange(userID int, startDate, endDate time.Time) ([]*domain.Event, error) {
	args := m.Called(userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Event), args.Error(1)
}

func TestCreateEvent(t *testing.T) {
	tests := []struct {
		name        string
		userID      int
		date        time.Time
		text        string
		expectError bool
		errorType   string
	}{
		{
			name:        "Успешное создание события",
			userID:      1,
			date:        time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			text:        "Тестовое событие",
			expectError: false,
		},
		{
			name:        "Пустой текст события",
			userID:      1,
			date:        time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			text:        "",
			expectError: true,
			errorType:   "validation",
		},
		{
			name:        "Некорректный user_id",
			userID:      0,
			date:        time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			text:        "Тестовое событие",
			expectError: true,
			errorType:   "validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockEventRepository)
			service := NewEventService(mockRepo)

			if !tt.expectError {
				mockRepo.On("Create", mock.AnythingOfType("*domain.Event")).Return(nil)
			}

			event, err := service.CreateEvent(tt.userID, tt.date, tt.text)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType == "validation" {
					appErr, ok := err.(*domain.AppError)
					assert.True(t, ok)
					assert.Equal(t, domain.StatusBadRequest, appErr.GetStatusCode())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, event)
				assert.Equal(t, tt.userID, event.UserID)
				assert.Equal(t, tt.date, event.Date)
				assert.Equal(t, tt.text, event.Text)
				assert.NotZero(t, event.CreatedAt)
				assert.NotZero(t, event.UpdatedAt)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	tests := []struct {
		name          string
		id            int
		userID        int
		date          time.Time
		text          string
		existingEvent *domain.Event
		expectError   bool
		errorType     string
	}{
		{
			name:   "Успешное обновление события",
			id:     1,
			userID: 1,
			date:   time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			text:   "Обновленный текст",
			existingEvent: &domain.Event{
				ID:     1,
				UserID: 1,
				Date:   time.Date(2023, 12, 30, 0, 0, 0, 0, time.UTC),
				Text:   "Старый текст",
			},
			expectError: false,
		},
		{
			name:          "Событие не найдено",
			id:            999,
			userID:        1,
			date:          time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			text:          "Новый текст",
			existingEvent: nil,
			expectError:   true,
			errorType:     "not_found",
		},
		{
			name:   "Нет прав доступа",
			id:     1,
			userID: 2,
			date:   time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			text:   "Новый текст",
			existingEvent: &domain.Event{
				ID:     1,
				UserID: 1,
				Date:   time.Date(2023, 12, 30, 0, 0, 0, 0, time.UTC),
				Text:   "Старый текст",
			},
			expectError: true,
			errorType:   "access_denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockEventRepository)
			service := NewEventService(mockRepo)

			if tt.existingEvent != nil {
				mockRepo.On("GetByID", tt.id).Return(tt.existingEvent, nil)
				// Update вызывается только при успешном обновлении
				if !tt.expectError {
					mockRepo.On("Update", mock.AnythingOfType("*domain.Event")).Return(nil)
				}
			} else {
				mockRepo.On("GetByID", tt.id).Return(nil, domain.NewNotFoundError("событие не найдено"))
			}

			event, err := service.UpdateEvent(tt.id, tt.userID, tt.date, tt.text)

			if tt.expectError {
				assert.Error(t, err)
				switch tt.errorType {
				case "not_found":
					appErr, ok := err.(*domain.AppError)
					assert.True(t, ok)
					assert.Equal(t, domain.StatusNotFound, appErr.GetStatusCode())
				case "access_denied":
					appErr, ok := err.(*domain.AppError)
					assert.True(t, ok)
					assert.Equal(t, domain.StatusForbidden, appErr.GetStatusCode())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, event)
				assert.Equal(t, tt.date, event.Date)
				assert.Equal(t, tt.text, event.Text)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	tests := []struct {
		name          string
		id            int
		userID        int
		existingEvent *domain.Event
		expectError   bool
		errorType     string
	}{
		{
			name:   "Успешное удаление события",
			id:     1,
			userID: 1,
			existingEvent: &domain.Event{
				ID:     1,
				UserID: 1,
				Date:   time.Date(2023, 12, 30, 0, 0, 0, 0, time.UTC),
				Text:   "Событие для удаления",
			},
			expectError: false,
		},
		{
			name:          "Событие не найдено",
			id:            999,
			userID:        1,
			existingEvent: nil,
			expectError:   true,
			errorType:     "not_found",
		},
		{
			name:   "Нет прав доступа",
			id:     1,
			userID: 2,
			existingEvent: &domain.Event{
				ID:     1,
				UserID: 1,
				Date:   time.Date(2023, 12, 30, 0, 0, 0, 0, time.UTC),
				Text:   "Событие для удаления",
			},
			expectError: true,
			errorType:   "access_denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockEventRepository)
			service := NewEventService(mockRepo)

			if tt.existingEvent != nil {
				mockRepo.On("GetByID", tt.id).Return(tt.existingEvent, nil)
				// Delete вызывается только при успешном удалении
				if !tt.expectError {
					mockRepo.On("Delete", tt.id, tt.userID).Return(nil)
				}
			} else {
				mockRepo.On("GetByID", tt.id).Return(nil, domain.NewNotFoundError("событие не найдено"))
			}

			err := service.DeleteEvent(tt.id, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				switch tt.errorType {
				case "not_found":
					appErr, ok := err.(*domain.AppError)
					assert.True(t, ok)
					assert.Equal(t, domain.StatusNotFound, appErr.GetStatusCode())
				case "access_denied":
					appErr, ok := err.(*domain.AppError)
					assert.True(t, ok)
					assert.Equal(t, domain.StatusForbidden, appErr.GetStatusCode())
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
