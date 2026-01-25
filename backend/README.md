# HertzBoard Backend

Backend для collaborative workspace platform, построенный на CloudWeGo Hertz.

## Технологический стек

- **Go 1.23** - Язык программирования
- **Hertz 0.10.3** - HTTP framework (CloudWeGo)
- **PostgreSQL 18.1** - Основная база данных
- **Redis 8.4** - Кэш и pub/sub
- **NATS 2.10** - Message queue
- **MailHog** - Email тестирование (development)

## Быстрый старт

### Предварительные требования

- Go 1.23+
- Docker и Docker Compose
- Make (опционально)

### Установка зависимостей

```bash
# Из корня проекта
cd backend
go mod download
```

### Запуск инфраструктуры

```bash
# Из корня проекта
make dev-up

# Или используя docker-compose напрямую
docker-compose up -d
```

Это запустит:
- PostgreSQL на порту 5432
- Redis на порту 6379
- NATS на порту 4222
- MailHog SMTP на порту 1025
- MailHog UI на http://localhost:8025
- MinIO на порту 9000
- ClickHouse на порту 8123

### Настройка environment variables

Скопируйте `.env.example` в `.env` и настройте OAuth credentials:

```bash
# Из корня проекта
cp .env.example .env
```

Для OAuth необходимо получить credentials:

**Google OAuth:**
1. Перейдите в [Google Cloud Console](https://console.cloud.google.com/)
2. Создайте проект
3. Включите Google+ API
4. Создайте OAuth 2.0 credentials
5. Добавьте redirect URI: `http://localhost:8080/api/v1/auth/google/callback`

**GitHub OAuth:**
1. Перейдите в [GitHub Settings > Developer settings > OAuth Apps](https://github.com/settings/developers)
2. Создайте New OAuth App
3. Authorization callback URL: `http://localhost:8080/api/v1/auth/github/callback`

Добавьте credentials в `.env`:
```env
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
```

### Запуск API Gateway

```bash
cd backend
go run cmd/api-gateway/main.go
```

API будет доступен на `http://localhost:8080`

## API Endpoints

### Health Checks

- `GET /health` - Basic health check
- `GET /readiness` - Readiness check (проверяет DB, Redis)

### Authentication

- `POST /api/v1/auth/register` - Регистрация пользователя
- `POST /api/v1/auth/login` - Вход
- `POST /api/v1/auth/refresh` - Обновление access token
- `POST /api/v1/auth/logout` - Выход
- `POST /api/v1/auth/forgot-password` - Запрос сброса пароля
- `POST /api/v1/auth/reset-password` - Сброс пароля

### OAuth

- `GET /api/v1/auth/google` - Google OAuth redirect
- `GET /api/v1/auth/google/callback` - Google callback
- `GET /api/v1/auth/github` - GitHub OAuth redirect
- `GET /api/v1/auth/github/callback` - GitHub callback

### User Profile (Protected)

Требуют JWT токен в header: `Authorization: Bearer <token>`

- `GET /api/v1/users/me` - Получить профиль текущего пользователя
- `PUT /api/v1/users/me` - Обновить профиль
- `PUT /api/v1/users/me/password` - Сменить пароль

## Примеры использования

### Регистрация

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123",
    "name": "John Doe"
  }'
```

Ответ:
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "provider": "email",
    "email_verified": false,
    "created_at": "2026-01-25T12:00:00Z",
    "updated_at": "2026-01-25T12:00:00Z"
  },
  "tokens": {
    "access_token": "eyJhbGc...",
    "refresh_token": "uuid",
    "expires_at": "2026-01-25T12:15:00Z"
  }
}
```

### Вход

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123"
  }'
```

### Получение профиля

```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer eyJhbGc..."
```

### Обновление профиля

```bash
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer eyJhbGc..." \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "avatar_url": "https://example.com/avatar.jpg"
  }'
```

### Сброс пароля

1. Запросить токен:
```bash
curl -X POST http://localhost:8080/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com"
  }'
```

2. Проверьте MailHog UI (http://localhost:8025) для получения токена

3. Сбросить пароль:
```bash
curl -X POST http://localhost:8080/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "token": "token-from-email",
    "new_password": "NewSecurePass123"
  }'
```

## Структура проекта

```
backend/
├── cmd/
│   └── api-gateway/
│       └── main.go              # Entry point
├── internal/
│   ├── config/                  # Configuration
│   ├── database/                # Database connections
│   ├── middleware/              # HTTP middleware
│   ├── models/                  # Data models
│   ├── repository/              # Data access layer
│   ├── service/                 # Business logic
│   ├── handler/                 # HTTP handlers
│   └── router/                  # Route definitions
├── migrations/                  # SQL migrations
├── configs/
│   └── config.yaml             # Configuration file
├── go.mod
└── go.sum
```

## Конфигурация

Основная конфигурация находится в `configs/config.yaml`. Для переопределения используйте environment variables:

```yaml
app:
  port: 8080                    # APP_PORT
  debug: true                   # APP_DEBUG

database:
  host: localhost               # DB_HOST
  port: 5432                    # DB_PORT
  name: hertzboard             # DB_NAME
  user: hertzboard             # DB_USER
  password: password           # DB_PASSWORD

jwt:
  secret: your-secret          # JWT_SECRET
  access_token_expiry: 15m     # JWT_ACCESS_EXPIRY
  refresh_token_expiry: 168h   # JWT_REFRESH_EXPIRY
```

## Миграции

Миграции находятся в директории `migrations/` и применяются автоматически при запуске приложения.

Формат имени файла: `{version}_{description}.sql`

Примеры:
- `001_create_users_table.sql`
- `002_create_refresh_tokens_table.sql`
- `003_create_password_reset_tokens_table.sql`

## Email тестирование

В development режиме используется MailHog:

1. SMTP server на порту 1025 (без authentication)
2. Web UI на http://localhost:8025

Все отправленные email можно просмотреть в MailHog UI.

## Разработка

### Добавление нового endpoint

1. Создайте handler в `internal/handler/`
2. Добавьте route в `internal/router/router.go`
3. Если требуется, создайте service в `internal/service/`
4. Если требуется, создайте repository в `internal/repository/`

### Добавление middleware

1. Создайте middleware в `internal/middleware/`
2. Добавьте в global middleware в `router.Setup()` или в group

### Создание миграции

1. Создайте файл в `migrations/` с номером версии
2. Напишите SQL команды
3. Миграция применится автоматически при следующем запуске

## Troubleshooting

### База данных не подключается

```bash
# Проверьте статус контейнеров
docker-compose ps

# Проверьте логи PostgreSQL
docker-compose logs postgres

# Пересоздайте контейнеры
docker-compose down -v
docker-compose up -d
```

### Миграции не применяются

```bash
# Проверьте логи приложения при запуске
go run cmd/api-gateway/main.go

# Проверьте подключение к БД
docker-compose exec postgres psql -U hertzboard -d hertzboard
```

### OAuth не работает

1. Проверьте что credentials настроены в `.env`
2. Проверьте redirect URIs в OAuth консоли провайдера
3. Убедитесь что используете правильные scopes

## Лицензия

MIT
