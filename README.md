# UTMN MAP Go Backend

## С чего начать

### Требования:

- Golang (main language)
- golangcilint (linter)
- golines (formatter)
- sqlc (sql codegen)
- tern (миграции)

Установить все можно написав `make init`

### Конфигурация

- `config/<service-name>.yaml`

Пример:

```yaml
server:
  host: "0.0.0.0"
  port: 8000
database:
  dsn: "postgres://utmn_user:utmn_password@localhost:5432/utmn_map?sslmode=disable"
```

### Запуск

```
go run cmd/publicapi/main.go --config=config/publicapi.yaml

go run cmd/adminapi/main.go --config=config/adminapi.yaml
```

### Об архитектуре

Проект построен по многослойной архитектуре, разделяя ответственность между различными слоями:

- **`cmd/`** – Точка входа в приложение.

  - `publicapi/main.go` – Запуск публичного API.
  - `publicapi/app.go` – Инициализация приложения.
  - `publicapi/config.go` – Загрузка конфигурации.

- **`config/`** – Конфигурационные файлы.

  - `publicapi.yaml` – Настройки для публичного API.

- **`internal/`** – Основная бизнес-логика и реализация доменной модели.

  - `domain/` – Доменные сущности и их бизнес-логика.
    - `map/` - Поддомен карты
      - `entities/` – Описание объектов карты (здания, двери, этажи и т. д.).
      - `repository/` – Работа с базой данных (конвертеры, SQL-запросы, миграции).
      - `service/` – Логика работы с картой.
    - `route` – Граф маршрутов.
      - `entities/`
      - `repository/`
      - `service/`
    - `search` – Поиск.
      - `entities/`
      - `repository/`
      - `service/`
  - `entrypoints/`
    - `publicapi/http/handler/` – HTTP-обработчики запросов.

- **Файлы и инструменты сборки:**
  - `Makefile` – Скрипты сборки и управления проектом.
  - `docker-compose.yaml` – Конфигурация Docker-окружения.
  - `go.mod`, `go.sum` – Зависимости Go.

Итого:

- **`cmd/`** – Запуск сервисов.
- **`internal/domain/`** – Бизнес-логика и работа с данными.
- **`internal/entrypoints/`** – Обработчики запросов.
- **`config/`** – Конфигурации.
