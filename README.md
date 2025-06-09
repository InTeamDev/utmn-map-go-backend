# UTMN MAP Go Backend

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/InTeamDev/utmn-map-go-backend)

## Описание проекта

Backend-сервис на Go для работы с картой (здания, двери, этажи), маршрутизации и поиска по данным кампуса УТМН.

## Ключевые компоненты

```plain
.
├── api
│   └── openapi           # OpenAPI спецификации и сгенерированный код
├── cmd                   # Точки входа для каждого сервиса
├── config                # Файлы конфигурации и Docker-настройки
├── docs                  # Документация по API и архитектуре
├── dump                  # Экспорт схемы базы данных
├── internal              # Business-логика и доменные реализации
├── pkg                   # Утильные пакеты (валидация, БД)
├── temp-frontend         # Временный фронтенд для тестирования
├── docker-compose.yaml   # Сборка окружения через Docker
├── Dockerfile            # Образ для контейнеризации основного сервиса
├── Makefile              # Сценарии сборки и управления проектом
├── prometheus.yml        # Конфигурация Prometheus
├── dashboard.json        # Конфигурация Grafana дашбордов
└── README.md             # Этот файл
```

## С чего начать

### Требования

- Go 1.20+
- Make
- golangci-lint
- golines
- sqlc
- tern (миграции базы данных)

Установить зависимости и инструменты можно с помощью:

```bash
make init
```

### Конфигурация

Все сервисы читают настройки из `config/<service>.yaml` и при необходимости из `config/<service>.docker.yaml` для Docker.

Пример `config/publicapi.yaml`:

```yaml
server:
  host: "0.0.0.0"
  port: 8000
database:
  dsn: "postgres://utmn_user:utmn_password@localhost:5432/utmn_map?sslmode=disable"
```

### Запуск сервисов

Для локального запуска каждого сервиса выполните:

```bash
go run cmd/publicapi/main.go --config=config/publicapi.yaml
```

Аналогично для остальных сервисов:

```bash
go run cmd/adminapi/main.go --config=config/adminapi.yaml
go run cmd/authapi/main.go --config=config/authapi.yaml
go run cmd/bot/main.go --config=config/bot.yaml
```

## Структура каталогов

### `api/openapi`

- `adminapi/`, `authapi/`, `publicapi/` — OpenAPI спецификации и сгенерированный код (Go).
- Файлы:

  - `openapi.yaml` — схема API.
  - `oapi-config.yaml` — конфигурация генератора.
  - `<api>.gen.go` — сгенерированный код.

### `cmd`

Точки входа для сервисов:

- `adminapi`, `authapi`, `publicapi`, `bot`.
- Каждый каталог содержит `main.go`.

### `config`

- `*.yaml` — конфигурации для каждого сервиса.
- `*.docker.yaml` — настройки для Docker.
- `*.go` — структура конфигурации и загрузка.

### `internal`

Основная бизнес-логика и доменные модели.

- `domain/` — сущности и сервисы домена:

  - `auth`, `map`, `route`, `search`.

- `entrypoints/` — HTTP-обработчики (handler) и инициализация приложений.
- `middleware/` — JWT, basic auth, метрики.
- `migrations/` — SQL-миграции (tern).
- `server/` — общая точка запуска сервера.

### `pkg`

Утильные пакеты:

- `database` — инициализация подключения к PostgreSQL.
- `validate` — функции валидации (например, дверей).

### Дополнительно

- `docs/` — документация по backend и примеру использования API.
- `dump/campus_schema.json` — экспорт JSON-схемы базы данных.
- `temp-frontend/` — тестовый фронтенд (HTML+JS).
- `prometheus.yml` и `dashboard.json` — мониторинг и дашборды Grafana.

## Сборка и CI

- `Makefile` содержит цели для сборки, тестирования и проверки кода:

  - `make build`
  - `make lint`
  - `make test`
  - `make migrate`

## Важно

- Следуйте стилю Go (gofmt, golines).
- Покрывайте код тестами (unit и integration).
- Обновляйте OpenAPI спецификации и регенерируйте клиент после изменений.
