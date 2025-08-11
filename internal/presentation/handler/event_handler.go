package handler

import (
	"calendar/internal/application"
	"calendar/internal/domain"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// EventHandler обрабатывает HTTP-запросы для событий
type EventHandler struct {
	*BaseHandler
	eventService *application.EventService
}

// NewEventHandler создает новый экземпляр обработчика событий
func NewEventHandler(eventService *application.EventService) *EventHandler {
	return &EventHandler{
		BaseHandler:  NewBaseHandler(),
		eventService: eventService,
	}
}

// RegisterRoutes регистрирует маршруты для событий
func (h *EventHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_event", h.CreateEvent).Methods("POST")
	router.HandleFunc("/update_event", h.UpdateEvent).Methods("POST")
	router.HandleFunc("/delete_event", h.DeleteEvent).Methods("POST")
	router.HandleFunc("/events_for_day", h.GetEventsForDay).Methods("GET")
	router.HandleFunc("/events_for_week", h.GetEventsForWeek).Methods("GET")
	router.HandleFunc("/events_for_month", h.GetEventsForMonth).Methods("GET")
}

// CreateEvent создает новое событие
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	// Парсим и валидируем форму
	fields, err := h.GetValidator().ParseFormAndValidate(r, []string{"user_id", "date", "text"})
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Парсим и валидируем параметры
	userID, err := h.GetValidator().ParseAndValidateUserID(fields["user_id"])
	if err != nil {
		h.handleError(w, err)
		return
	}

	date, err := h.GetValidator().ParseAndValidateDate(fields["date"])
	if err != nil {
		h.handleError(w, err)
		return
	}

	text := fields["text"]

	// Создаем событие
	event, err := h.eventService.CreateEvent(userID, date, text)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeSuccess(w, event)
}

// UpdateEvent обновляет существующее событие
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	// Парсим и валидируем форму
	fields, err := h.GetValidator().ParseFormAndValidate(r, []string{"id", "user_id", "date", "text"})
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Парсим и валидируем параметры
	id, err := h.GetValidator().ParseAndValidateID(fields["id"])
	if err != nil {
		h.handleError(w, err)
		return
	}

	userID, err := h.GetValidator().ParseAndValidateUserID(fields["user_id"])
	if err != nil {
		h.handleError(w, err)
		return
	}

	date, err := h.GetValidator().ParseAndValidateDate(fields["date"])
	if err != nil {
		h.handleError(w, err)
		return
	}

	text := fields["text"]

	// Обновляем событие
	event, err := h.eventService.UpdateEvent(id, userID, date, text)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeSuccess(w, event)
}

// DeleteEvent удаляет событие
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	// Парсим и валидируем форму
	fields, err := h.GetValidator().ParseFormAndValidate(r, []string{"id", "user_id"})
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Парсим и валидируем параметры
	id, err := h.GetValidator().ParseAndValidateID(fields["id"])
	if err != nil {
		h.handleError(w, err)
		return
	}

	userID, err := h.GetValidator().ParseAndValidateUserID(fields["user_id"])
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Удаляем событие
	err = h.eventService.DeleteEvent(id, userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeSuccess(w, map[string]string{"message": "Событие успешно удалено"})
}

// getEventsByDateRange общий метод для получения событий по диапазону дат
func (h *EventHandler) getEventsByDateRange(w http.ResponseWriter, r *http.Request, getter func(int, time.Time) ([]*domain.Event, error)) {
	// Извлекаем параметры из query string
	userIDStr := r.URL.Query().Get("user_id")
	dateStr := r.URL.Query().Get("date")

	if userIDStr == "" || dateStr == "" {
		h.writeError(w, http.StatusBadRequest, "Необходимы параметры: user_id, date")
		return
	}

	// Парсим и валидируем параметры
	userID, err := h.GetValidator().ParseAndValidateUserID(userIDStr)
	if err != nil {
		h.handleError(w, err)
		return
	}

	date, err := h.GetValidator().ParseAndValidateDate(dateStr)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Получаем события
	events, err := getter(userID, date)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeSuccess(w, events)
}

// GetEventsForDay возвращает события на день
func (h *EventHandler) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	h.getEventsByDateRange(w, r, h.eventService.GetEventsForDay)
}

// GetEventsForWeek возвращает события на неделю
func (h *EventHandler) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	h.getEventsByDateRange(w, r, h.eventService.GetEventsForWeek)
}

// GetEventsForMonth возвращает события на месяц
func (h *EventHandler) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры из query string
	userIDStr := r.URL.Query().Get("user_id")
	yearMonthStr := r.URL.Query().Get("date")

	if userIDStr == "" || yearMonthStr == "" {
		h.writeError(w, http.StatusBadRequest, "Необходимы параметры: user_id, date")
		return
	}

	// Парсим и валидируем параметры
	userID, err := h.GetValidator().ParseAndValidateUserID(userIDStr)
	if err != nil {
		h.handleError(w, err)
		return
	}

	yearMonth, err := h.GetValidator().ParseAndValidateYearMonth(yearMonthStr)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Получаем события
	events, err := h.eventService.GetEventsForMonth(userID, yearMonth)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeSuccess(w, events)
}
