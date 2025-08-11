package domain

import (
	"time"
)

// Event представляет событие в календаре
type Event struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Date      time.Time `json:"date"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// EventRepository определяет интерфейс для работы с событиями
type EventRepository interface {
	Create(event *Event) error
	Update(event *Event) error
	Delete(id int, userID int) error
	GetByID(id int) (*Event, error)
	GetByUserAndDate(userID int, date time.Time) ([]*Event, error)
	GetByUserAndDateRange(userID int, startDate, endDate time.Time) ([]*Event, error)
}

// EventService определяет бизнес-логику для работы с событиями
type EventService interface {
	CreateEvent(userID int, date time.Time, text string) (*Event, error)
	UpdateEvent(id int, userID int, date time.Time, text string) (*Event, error)
	DeleteEvent(id int, userID int) error
	GetEventsForDay(userID int, date time.Time) ([]*Event, error)
	GetEventsForWeek(userID int, startDate time.Time) ([]*Event, error)
	GetEventsForMonth(userID int, yearMonth time.Time) ([]*Event, error)
}
