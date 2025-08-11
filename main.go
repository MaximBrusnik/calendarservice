package main

import (
	"log"
	"os"

	"calendar/internal/presentation/server"
)

func main() {
	// Получаем порт из переменной окружения или используем значение по умолчанию
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Создаем и запускаем сервер
	srv := server.NewServer(port)

	log.Printf("Сервер календаря запущен на порту %s", port)
	log.Printf("Health check: http://localhost:%s/health", port)

	if err := srv.Start(); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
