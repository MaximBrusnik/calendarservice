package server

import (
	"fmt"
	"net/http"

	"calendar/internal/application"
	"calendar/internal/infrastructure/repository"
	"calendar/internal/presentation/handler"
	"calendar/internal/presentation/middleware"

	"github.com/gorilla/mux"
)

// Server представляет HTTP-сервер
type Server struct {
	router       *mux.Router
	port         string
	eventHandler *handler.EventHandler
}

// NewServer создает новый экземпляр HTTP-сервера
func NewServer(port string) *Server {
	// Создаем репозиторий
	eventRepo := repository.NewMemoryEventRepository()

	// Создаем сервис приложения
	eventService := application.NewEventService(eventRepo)

	// Создаем обработчик
	eventHandler := handler.NewEventHandler(eventService)

	// Создаем роутер
	router := mux.NewRouter()

	// middleware
	router.Use(middleware.LoggingMiddleware)

	// Создаем сервер
	server := &Server{
		router:       router,
		port:         port,
		eventHandler: eventHandler,
	}

	// Настраиваем маршруты
	server.setupRoutes()

	return server
}

// setupRoutes настраивает маршруты сервера
func (s *Server) setupRoutes() {
	// Регистрируем маршруты для событий
	s.eventHandler.RegisterRoutes(s.router)

	// Добавляем health check endpoint
	s.router.HandleFunc("/health", s.healthCheck).Methods("GET")
}

// healthCheck обрабатывает запрос на проверку состояния сервера
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "ok", "service": "calendar"}`)
}

// Start запускает HTTP-сервер
func (s *Server) Start() error {
	addr := ":" + s.port
	fmt.Printf("Сервер запущен на порту %s\n", s.port)
	return http.ListenAndServe(addr, s.router)
}
