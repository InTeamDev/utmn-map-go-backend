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

## Telegram Bot Authentication System

This project includes a Telegram bot for user authentication and role management.

### Features

1. **User Registration via Deep Link**
   - The web server generates a registration hash (TTL 5 min)
   - Deep link format: `https://t.me/utmn_map_bot?start=reg_hash`
   - The bot validates the hash and registers the user with User role

2. **Authentication by Code**
   - User enters their Telegram username in the web interface
   - The server generates a one-time code (TTL 5 min)
   - The bot sends the code to the user
   - User enters the code to receive a JWT token

3. **Role Management**
   - Bot commands for curators:
     - `/promote @username` - Assign Admin role
     - `/demote @username` - Remove Admin role
   - Primary Curator is assigned manually via database seed

4. **Event Notifications**
   - The bot publishes events to a specified developers chat:
     - User registrations
     - Role changes
     - Authentication code requests

### Setup

1. Create a new Telegram bot using BotFather and get the token
2. Create a group chat for developers and get its chat ID
3. Update `config/publicapi.yaml` with the token and chat ID:

```yaml
tgbot:
  token: "YOUR_TELEGRAM_BOT_TOKEN"
  developers_chat_id: -123456789 # Replace with actual chat ID
```

4. Update the JWT secret in the config:

```yaml
jwt:
  secret: "your-secure-secret-key"
  expiration_hours: 24
```

### Testing the Bot

1. Start the server: `make run-publicapi`
2. Generate a registration hash: `curl -X POST http://localhost:8090/api/auth/register/generate-hash`
3. Use the deep link from the response to register via Telegram
4. Request an auth code: `curl -X POST -H "Content-Type: application/json" -d '{"username":"@yourusername"}' http://localhost:8090/api/auth/login/generate-code`
5. Use the code sent to your Telegram to get a JWT token: `curl -X POST -H "Content-Type: application/json" -d '{"username":"@yourusername", "code":"123456"}' http://localhost:8090/api/auth/login`

### API Endpoints

- `POST /api/auth/register/generate-hash` - Generate registration hash
- `POST /api/auth/register` - Register user (called by the bot)
- `POST /api/auth/login/generate-code` - Generate one-time code
- `POST /api/auth/login` - Login with username and code
