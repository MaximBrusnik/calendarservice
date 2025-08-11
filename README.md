# HTTP-сервер "Календарь" - 

HTTP-сервер для работы с небольшим календарем событий, реализованный на Go с соблюдением принципов Clean Architecture.

## Архитектура

Проект следует принципам Clean Architecture и разделен на следующие слои:

```
calendar/
├── cmd/                    # Точка входа в приложение
├── internal/              # Внутренний код приложения
│   ├── domain/           # Доменный слой (бизнес-логика, интерфейсы)
│   ├── application/      # Слой приложения (сервисы, use cases)
│   ├── infrastructure/   # Слой инфраструктуры (репозитории, внешние сервисы)
│   └── presentation/     # Слой представления (HTTP handlers, middleware)
├── go.mod                # Зависимости Go
└── README.md             # Документация
```

### Слои архитектуры:

1. **Domain Layer** (`internal/domain/`)
   - Бизнес-модели (`Event`)
   - Интерфейсы репозиториев (`EventRepository`)
   - Интерфейсы сервисов (`EventService`)
   - Доменные ошибки

2. **Application Layer** (`internal/application/`)
   - Бизнес-логика (`EventService`)
   - Валидация данных
   - Координация между доменными объектами

3. **Infrastructure Layer** (`internal/infrastructure/`)
   - Реализация репозиториев (`MemoryEventRepository`)
   - Внешние сервисы
   - База данных

4. **Presentation Layer** (`internal/presentation/`)
   - HTTP обработчики (`EventHandler`)
   - Middleware (`LoggingMiddleware`)
   - HTTP сервер (`Server`)

## Возможности

- **CRUD операции для событий:**
  - Создание нового события
  - Обновление существующего события
  - Удаление события
  - Получение событий на день, неделю, месяц

- **Безопасность:** Проверка прав доступа пользователей к событиям
- **Валидация:** Проверка корректности входных данных
- **Логирование:** Middleware для логирования всех запросов
- **Тестирование:** Покрытие unit-тестами с использованием моков
- **Чистая архитектура:** Четкое разделение ответственности между слоями

## API Endpoints

### Создание события
```
POST /create_event
Content-Type: application/x-www-form-urlencoded

user_id=1&date=2025-12-18&text=Текст1
```

### Обновление события
```
POST /update_event
Content-Type: application/x-www-form-urlencoded

id=1&user_id=1&date=2025-12-18&text=Текст2 
```

### Удаление события
```
POST /delete_event
Content-Type: application/x-www-form-urlencoded

id=1&user_id=1
```

### Получение событий на день
```
GET /events_for_day?user_id=1&date=2025-12-18
```

### Получение событий на неделю
```
GET /events_for_week?user_id=1&date=2025-12-18
```

### Получение событий на месяц
```
GET /events_for_month?user_id=1&date=2025-12
```

### Health Check
```
GET /health
```

## Форматы данных

- **Дата:** YYYY-MM-DD (например, 2025-12-18)
- **Месяц:** YYYY-MM (например, 2025-12)
- **user_id:** Целое число, идентификатор пользователя
- **id:** Целое число, идентификатор события

## Ответы API

### Успешный ответ
```json
{
  "result": "данные"
}
```

### Ответ с ошибкой
```json
{
  "error": "описание ошибки"
}
```

## HTTP статус-коды

- **200 OK** - успешное выполнение запроса
- **400 Bad Request** - ошибки ввода (некорректные параметры)
- **503 Service Unavailable** - ошибки бизнес-логики (событие не найдено, нет прав)
- **500 Internal Server Error** - прочие ошибки

## Установка и запуск

### Требования
- Go 1.21 или выше

### Установка зависимостей
```bash
go mod tidy
```

### Запуск сервера
```bash
go run main.go
```

### Запуск с указанием порта
```bash
PORT=8080 go run main.go
```

### Проверка качества кода
```bash
# Проверка с помощью go vet
go vet ./...

# Проверка с помощью golangci-lint
golangci-lint run
```


### Структура тестов:
- `internal/application/event_service_test.go` - тесты бизнес-логики
- `internal/infrastructure/repository/memory_event_repository_test.go` - тесты репозитория


