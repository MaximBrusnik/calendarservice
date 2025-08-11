package repository

import (
	"calendar/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryEventRepository_Create(t *testing.T) {
	repo := NewMemoryEventRepository()
	event := &domain.Event{
		UserID: 1,
		Date:   time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
		Text:   "Тестовое событие",
	}

	err := repo.Create(event)
	assert.NoError(t, err)
	assert.Equal(t, 1, event.ID)
	// CreatedAt и UpdatedAt устанавливаются в сервисе, а не в репозитории
}

func TestMemoryEventRepository_GetByID(t *testing.T) {
	repo := NewMemoryEventRepository()
	event := &domain.Event{
		UserID: 1,
		Date:   time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
		Text:   "Тестовое событие",
	}

	// Создаем событие
	err := repo.Create(event)
	assert.NoError(t, err)

	// Получаем событие по ID
	retrievedEvent, err := repo.GetByID(event.ID)
	assert.NoError(t, err)
	assert.Equal(t, event, retrievedEvent)

	// Пытаемся получить несуществующее событие
	_, err = repo.GetByID(999)
	assert.Error(t, err)
	appErr, ok := err.(*domain.AppError)
	assert.True(t, ok)
	assert.Equal(t, domain.StatusNotFound, appErr.GetStatusCode())
}

func TestMemoryEventRepository_Update(t *testing.T) {
	repo := NewMemoryEventRepository()
	event := &domain.Event{
		UserID: 1,
		Date:   time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
		Text:   "Тестовое событие",
	}

	// Создаем событие
	err := repo.Create(event)
	assert.NoError(t, err)

	// Обновляем событие
	event.Text = "Обновленное событие"
	err = repo.Update(event)
	assert.NoError(t, err)

	// Проверяем, что событие обновлено
	retrievedEvent, err := repo.GetByID(event.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Обновленное событие", retrievedEvent.Text)

	// Пытаемся обновить несуществующее событие
	nonExistentEvent := &domain.Event{ID: 999, UserID: 1, Date: time.Now(), Text: "Несуществующее"}
	err = repo.Update(nonExistentEvent)
	assert.Error(t, err)
	appErr, ok := err.(*domain.AppError)
	assert.True(t, ok)
	assert.Equal(t, domain.StatusNotFound, appErr.GetStatusCode())
}

func TestMemoryEventRepository_Delete(t *testing.T) {
	repo := NewMemoryEventRepository()
	event := &domain.Event{
		UserID: 1,
		Date:   time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
		Text:   "Тестовое событие",
	}

	// Создаем событие
	err := repo.Create(event)
	assert.NoError(t, err)

	// Удаляем событие
	err = repo.Delete(event.ID, event.UserID)
	assert.NoError(t, err)

	// Проверяем, что событие удалено
	_, err = repo.GetByID(event.ID)
	assert.Error(t, err)
	appErr, ok := err.(*domain.AppError)
	assert.True(t, ok)
	assert.Equal(t, domain.StatusNotFound, appErr.GetStatusCode())
}

func TestMemoryEventRepository_Delete_AccessDenied(t *testing.T) {
	repo := NewMemoryEventRepository()
	event := &domain.Event{
		UserID: 1,
		Date:   time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
		Text:   "Тестовое событие",
	}

	// Создаем событие
	err := repo.Create(event)
	assert.NoError(t, err)

	// Пытаемся удалить событие от имени другого пользователя
	err = repo.Delete(event.ID, 2)
	assert.Error(t, err)
	appErr, ok := err.(*domain.AppError)
	assert.True(t, ok)
	assert.Equal(t, domain.StatusForbidden, appErr.GetStatusCode())
}

func TestMemoryEventRepository_GetByUserAndDate(t *testing.T) {
	repo := NewMemoryEventRepository()

	// Создаем несколько событий для одного пользователя
	event1 := &domain.Event{
		UserID: 1,
		Date:   time.Date(2023, 12, 31, 10, 0, 0, 0, time.UTC),
		Text:   "Событие 1",
	}
	event2 := &domain.Event{
		UserID: 1,
		Date:   time.Date(2023, 12, 31, 15, 0, 0, 0, time.UTC),
		Text:   "Событие 2",
	}
	event3 := &domain.Event{
		UserID: 1,
		Date:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Text:   "Событие 3",
	}

	repo.Create(event1)
	repo.Create(event2)
	repo.Create(event3)

	// Получаем события на 31 декабря
	events, err := repo.GetByUserAndDate(1, time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC))
	assert.NoError(t, err)
	assert.Len(t, events, 2)
}

func TestMemoryEventRepository_GetByUserAndDateRange(t *testing.T) {
	repo := NewMemoryEventRepository()

	// Создаем несколько событий
	event1 := &domain.Event{
		UserID: 1,
		Date:   time.Date(2023, 12, 30, 0, 0, 0, 0, time.UTC),
		Text:   "Событие 1",
	}
	event2 := &domain.Event{
		UserID: 1,
		Date:   time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
		Text:   "Событие 2",
	}
	event3 := &domain.Event{
		UserID: 1,
		Date:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Text:   "Событие 3",
	}

	repo.Create(event1)
	repo.Create(event2)
	repo.Create(event3)

	// Получаем события за период с 30 декабря по 1 января
	startDate := time.Date(2023, 12, 30, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 1, 23, 59, 59, 0, time.UTC)

	events, err := repo.GetByUserAndDateRange(1, startDate, endDate)
	assert.NoError(t, err)
	assert.Len(t, events, 3)
}

func TestMemoryEventRepository_ConcurrentAccess(t *testing.T) {
	repo := NewMemoryEventRepository()

	// Создаем несколько горутин для одновременного доступа
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			event := &domain.Event{
				UserID: id,
				Date:   time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
				Text:   "Событие",
			}

			repo.Create(event)
			repo.GetByID(event.ID)
			repo.Update(event)
			repo.Delete(event.ID, event.UserID)

			done <- true
		}(i)
	}

	// Ждем завершения всех горутин
	for i := 0; i < 10; i++ {
		<-done
	}

	// Проверяем, что репозиторий не поврежден
	assert.Equal(t, 11, repo.nextID) // 10 событий + 1 для следующего ID
}
