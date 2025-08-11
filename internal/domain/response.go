package domain

// Response представляет стандартный ответ API
type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// CreateEventRequest представляет запрос на создание события
type CreateEventRequest struct {
	UserID int    `json:"user_id" form:"user_id"`
	Date   string `json:"date" form:"date"`
	Text   string `json:"text" form:"text"`
}

// UpdateEventRequest представляет запрос на обновление события
type UpdateEventRequest struct {
	ID     int    `json:"id" form:"id"`
	UserID int    `json:"user_id" form:"user_id"`
	Date   string `json:"date" form:"date"`
	Text   string `json:"text" form:"text"`
}

// DeleteEventRequest представляет запрос на удаление события
type DeleteEventRequest struct {
	ID     int `json:"id" form:"id"`
	UserID int `json:"user_id" form:"user_id"`
}

// GetEventsRequest представляет запрос на получение событий
type GetEventsRequest struct {
	UserID int    `json:"user_id" form:"user_id"`
	Date   string `json:"date" form:"date"`
}
