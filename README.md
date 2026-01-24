# HertzBoard

> Collaborative Workspace Platform - Бесконечный canvas для создания заметок, рисунков, диаграмм и организации информации в свободной форме.

[![License](https://img.shields.io/badge/license-GPL--3.0-blue.svg)](LICENSE)

## Особенности

- **Бесконечный Canvas** - Свободное размещение элементов в 2D пространстве
- **Real-time Коллаборация** - Работайте вместе с командой в реальном времени
- **CRDT-based Синхронизация** - Бесконфликтная синхронизация изменений
- **Богатый набор инструментов** - Текст, фигуры, изображения, рисование от руки
- **Кросс-платформенность** - Web, Desktop (Tauri), Mobile (Capacitor)

## Технологический Стек

### Backend
- **Hertz** - HTTP/WebSocket сервер
- **Kitex** - gRPC микросервисы
- **PostgreSQL** - Основная база данных
- **Redis** - Кэш и pub/sub
- **MinIO** - Хранилище файлов
- **ClickHouse** - Аналитика
- **NATS** - Очереди сообщений

### Frontend
- **SvelteKit** - Фреймворк
- **TypeScript** - Типизация
- **Tailwind CSS** - Стилизация
- **Yjs** - CRDT синхронизация
- **TipTap** - Rich text редактор

## Быстрый старт

### Требования

- Go 1.22+
- Node.js 20+
- Docker & Docker Compose
- Make (опционально)

### Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/bifshteksex/hertz-board.git
cd hertz-board
```

2. Инициализируйте проект:
```bash
make init
```

Эта команда:
- Создаст `.env` файл из `.env.example`
- Установит зависимости для backend и frontend
- Запустит Docker Compose с сервисами
- Выполнит миграции базы данных

### Запуск для разработки

#### Вариант 1: Используя Make (рекомендуется)

```bash
# Запустить все сервисы
make dev-up

# В отдельных терминалах:
make backend-run    # API Gateway
make frontend-run   # Frontend dev server
```

#### Вариант 2: Вручную

1. Запустите инфраструктуру:
```bash
docker-compose up -d
```

2. Запустите backend:
```bash
cd backend
go run cmd/api-gateway/main.go
```

3. Запустите frontend:
```bash
cd frontend
npm run dev
```

### Доступные сервисы

После запуска будут доступны:

- **Frontend**: http://localhost:5173
- **API Gateway**: http://localhost:8080
- **WebSocket Server**: ws://localhost:8081
- **MinIO Console**: http://localhost:9001 (hertzboard / hertzboard_minio_password)
- **MailHog UI**: http://localhost:8025
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

## Структура проекта

```
hertzboard/
├── backend/              # Go backend services
│   ├── cmd/             # Точки входа для сервисов
│   ├── internal/        # Внутренняя логика
│   ├── configs/         # Конфигурационные файлы
│   └── scripts/         # Скрипты миграций
├── frontend/            # SvelteKit frontend
│   └── src/            
│       ├── lib/         # Компоненты, stores, utils
│       └── routes/      # Страницы приложения
├── desktop/             # Tauri desktop app
├── mobile/              # Capacitor mobile app
├── deploy/              # Kubernetes, Docker configs
└── docs/                # Документация
```

## Доступные команды

### Make команды

```bash
make help              # Показать все доступные команды
make dev-up            # Запустить dev окружение
make dev-down          # Остановить dev окружение
make install           # Установить зависимости
make test              # Запустить все тесты
make lint              # Запустить линтеры
make build             # Собрать все сервисы
make clean             # Очистить build артефакты
```

### Backend команды

```bash
cd backend
go run cmd/api-gateway/main.go    # Запустить API Gateway
go run cmd/ws-server/main.go      # Запустить WebSocket Server
go test ./...                     # Запустить тесты
```

### Frontend команды

```bash
cd frontend
npm run dev           # Dev server
npm run build         # Production build
npm run preview       # Preview production build
npm test              # Запустить тесты
npm run lint          # Линтер
```

## Тестирование

### Backend тесты

```bash
make backend-test
# или
cd backend && go test -v ./...
```

### Frontend тесты

```bash
make frontend-test
# или
cd frontend && npm run test
```

### Все тесты

```bash
make test
```

## Деплой

### Staging

Push в ветку `develop` автоматически запустит деплой в staging окружение.

### Production

Создайте тег версии для деплоя в production:

```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

## Разработка

### Структура Backend

Backend следует Clean Architecture:
- `domain/` - Бизнес-логика и entities
- `infrastructure/` - Внешние зависимости (DB, Redis, MinIO)
- `interfaces/` - HTTP, WebSocket, gRPC handlers
- `pkg/` - Переиспользуемые пакеты

### Структура Frontend

Frontend организован по функциональным модулям:
- `components/` - Reusable UI компоненты
- `stores/` - Svelte stores для state management
- `services/` - API клиенты и сервисы
- `utils/` - Утилиты
- `types/` - TypeScript типы

### Git Flow

1. Создайте ветку от `develop`:
```bash
git checkout -b feature/your-feature-name
```

2. Сделайте изменения и коммит:
```bash
git add .
git commit -m "feat: add new feature"
```

3. Отправьте в remote и создайте Pull Request:
```bash
git push origin feature/your-feature-name
```

### Commit Convention

Используем [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - Новая функциональность
- `fix:` - Исправление бага
- `docs:` - Изменения в документации
- `style:` - Форматирование кода
- `refactor:` - Рефакторинг
- `test:` - Добавление тестов
- `chore:` - Рутинные задачи

## Документация

Подробная документация доступна в папке `docs/`:

- [Архитектура](docs/architecture/overview.md)
- [API Документация](docs/api/rest-api.md)
- [WebSocket Protocol](docs/api/websocket-api.md)
- [Руководство по разработке](docs/development/setup.md)

## Лицензия

GPL-3.0 License - см. [LICENSE](LICENSE) файл для деталей.

---

Made with ❤️ by Roman Shangin
