package application

import (
	"calendar/internal/domain"
	"time"
)

// EventService реализует бизнес-логику для работы с событиями
type EventService struct {
	repo      domain.EventRepository
	validator *ServiceValidator
}

// NewEventService создает новый экземпляр сервиса событий
func NewEventService(repo domain.EventRepository) *EventService {
	return &EventService{
		repo:      repo,
		validator: NewServiceValidator(),
	}
}

// CreateEvent создает новое событие
func (s *EventService) CreateEvent(userID int, date time.Time, text string) (*domain.Event, error) {
	// Валидация входных данных
	if err := s.validator.ValidateEventData(userID, text); err != nil {
		return nil, err
	}

	// Создаем событие
	event := &domain.Event{
		UserID:    userID,
		Date:      date,
		Text:      text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Сохраняем в репозитории
	if err := s.repo.Create(event); err != nil {
		return nil, domain.NewInternalError("ошибка при создании события", err)
	}

	return event, nil
}

// UpdateEvent обновляет существующее событие
func (s *EventService) UpdateEvent(id int, userID int, date time.Time, text string) (*domain.Event, error) {
	// Валидация входных данных
	if err := s.validator.ValidateEventData(userID, text); err != nil {
		return nil, err
	}

	if err := s.validator.ValidateEventID(id); err != nil {
		return nil, err
	}

	// Получаем существующее событие
	event, err := s.repo.GetByID(id)
	if err != nil {
		return nil, domain.NewNotFoundError("событие не найдено")
	}

	// Проверяем права доступа
	if event.UserID != userID {
		return nil, domain.NewAccessDeniedError("нет прав для изменения этого события")
	}

	// Обновляем поля
	event.Date = date
	event.Text = text
	event.UpdatedAt = time.Now()

	// Сохраняем изменения
	if err := s.repo.Update(event); err != nil {
		return nil, domain.NewInternalError("ошибка при обновлении события", err)
	}

	return event, nil
}

// DeleteEvent удаляет событие
func (s *EventService) DeleteEvent(id int, userID int) error {
	// Валидация входных данных
	if err := s.validator.ValidateEventID(id); err != nil {
		return err
	}

	if err := s.validator.ValidateUserID(userID); err != nil {
		return err
	}

	// Получаем существующее событие
	event, err := s.repo.GetByID(id)
	if err != nil {
		return domain.NewNotFoundError("событие не найдено")
	}

	// Проверяем права доступа
	if event.UserID != userID {
		return domain.NewAccessDeniedError("нет прав для удаления этого события")
	}

	// Удаляем событие
	if err := s.repo.Delete(id, userID); err != nil {
		return domain.NewInternalError("ошибка при удалении события", err)
	}

	return nil
}

// getEventsByUserID общий метод для получения событий пользователя
func (s *EventService) getEventsByUserID(userID int, getter func(int, time.Time) ([]*domain.Event, error), date time.Time) ([]*domain.Event, error) {
	if err := s.validator.ValidateUserID(userID); err != nil {
		return nil, err
	}

	events, err := getter(userID, date)
	if err != nil {
		return nil, domain.NewInternalError("ошибка при получении событий", err)
	}

	return events, nil
}

// GetEventsForDay возвращает события на конкретный день
func (s *EventService) GetEventsForDay(userID int, date time.Time) ([]*domain.Event, error) {
	return s.getEventsByUserID(userID, s.repo.GetByUserAndDate, date)
}

// GetEventsForWeek возвращает события на неделю, начиная с указанной даты
func (s *EventService) GetEventsForWeek(userID int, startDate time.Time) ([]*domain.Event, error) {
	endDate := startDate.AddDate(0, 0, 7)
	return s.getEventsByUserID(userID, func(uid int, date time.Time) ([]*domain.Event, error) {
		return s.repo.GetByUserAndDateRange(uid, startDate, endDate)
	}, startDate)
}

// GetEventsForMonth возвращает события на месяц
func (s *EventService) GetEventsForMonth(userID int, yearMonth time.Time) ([]*domain.Event, error) {
	// Начало месяца
	startDate := time.Date(yearMonth.Year(), yearMonth.Month(), 1, 0, 0, 0, 0, yearMonth.Location())
	// Конец месяца
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return s.getEventsByUserID(userID, func(uid int, date time.Time) ([]*domain.Event, error) {
		return s.repo.GetByUserAndDateRange(uid, startDate, endDate)
	}, yearMonth)
}
