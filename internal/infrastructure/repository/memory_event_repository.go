package repository

import (
	"calendar/internal/domain"
	"sync"
	"time"
)

// MemoryEventRepository реализует in-memory репозиторий для событий
type MemoryEventRepository struct {
	events map[int]*domain.Event
	users  map[int][]int // map[userID][]eventIDs
	mu     sync.RWMutex
	nextID int
}

// NewMemoryEventRepository создает новый экземпляр in-memory репозитория
func NewMemoryEventRepository() *MemoryEventRepository {
	return &MemoryEventRepository{
		events: make(map[int]*domain.Event),
		users:  make(map[int][]int),
		nextID: 1,
	}
}

// Create создает новое событие
func (r *MemoryEventRepository) Create(event *domain.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	event.ID = r.nextID
	r.events[event.ID] = event
	r.users[event.UserID] = append(r.users[event.UserID], event.ID)
	r.nextID++

	return nil
}

// Update обновляет существующее событие
func (r *MemoryEventRepository) Update(event *domain.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[event.ID]; !exists {
		return domain.NewNotFoundError("событие не найдено")
	}

	r.events[event.ID] = event
	return nil
}

// Delete удаляет событие
func (r *MemoryEventRepository) Delete(id int, userID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	event, exists := r.events[id]
	if !exists {
		return domain.NewNotFoundError("событие не найдено")
	}

	if event.UserID != userID {
		return domain.NewAccessDeniedError("нет прав для удаления этого события")
	}

	// Удаляем событие
	delete(r.events, id)

	// Удаляем из списка пользователя
	if userEvents, ok := r.users[userID]; ok {
		for i, eventID := range userEvents {
			if eventID == id {
				r.users[userID] = append(userEvents[:i], userEvents[i+1:]...)
				break
			}
		}
	}

	return nil
}

// GetByID возвращает событие по ID
func (r *MemoryEventRepository) GetByID(id int) (*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	event, exists := r.events[id]
	if !exists {
		return nil, domain.NewNotFoundError("событие не найдено")
	}

	return event, nil
}

// GetByUserAndDate возвращает события пользователя на конкретную дату
func (r *MemoryEventRepository) GetByUserAndDate(userID int, date time.Time) ([]*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var events []*domain.Event
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	for _, eventID := range r.users[userID] {
		if event, exists := r.events[eventID]; exists {
			if event.Date.After(startOfDay) && event.Date.Before(endOfDay) {
				events = append(events, event)
			}
		}
	}

	return events, nil
}

// GetByUserAndDateRange возвращает события пользователя в указанном диапазоне дат
func (r *MemoryEventRepository) GetByUserAndDateRange(userID int, startDate, endDate time.Time) ([]*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var events []*domain.Event

	for _, eventID := range r.users[userID] {
		if event, exists := r.events[eventID]; exists {
			if (event.Date.After(startDate) || event.Date.Equal(startDate)) &&
				(event.Date.Before(endDate) || event.Date.Equal(endDate)) {
				events = append(events, event)
			}
		}
	}

	return events, nil
}
