# HertzBoard

[![Release](https://img.shields.io/github/v/release/bifshteksex/hertz-board)](https://github.com/bifshteksex/hertz-board/releases)
[![License](https://img.shields.io/badge/license-GPL--3.0-blue.svg)](LICENSE)

> **🇬🇧 English** | [🇷🇺 Русский](#ru) | [🇨🇳 中文](#zh)

## 🇬🇧 English

> Collaborative Workspace Platform - An infinite canvas for creating notes, drawings, diagrams, and organizing information in free form.

### Features

- **Infinite Canvas** - Free placement of elements in 2D space
- **Real-time Collaboration** - Work together with your team in real-time
- **CRDT-based Synchronization** - Conflict-free synchronization of changes
- **Rich Toolset** - Text, shapes, images, freehand drawing
- **Cross-platform** - Web, Desktop (Tauri), Mobile (Capacitor)

### Technology Stack

#### Backend
- **Hertz v0.10.3** - HTTP/WebSocket server
- **Kitex** - gRPC microservices
- **PostgreSQL** - Primary database
- **Redis** - Cache and pub/sub
- **MinIO** - File storage
- **ClickHouse** - Analytics
- **NATS** - Message queues

#### Frontend
- **SvelteKit 2.50+** - Framework
- **TypeScript 5.9+** - Type safety
- **Tailwind CSS 4.1+** - Styling
- **Yjs 13.6+** - CRDT synchronization
- **TipTap 3.18+** - Rich text editor
- **Vite 7.3+** - Build tool

### Quick Start

#### Requirements

- **Go 1.25.6+**
- **Node.js 20+** (pnpm or npm)
- **Docker & Docker Compose**
- **Make** (optional)

#### Installation

1. Clone the repository:
```bash
git clone https://github.com/bifshteksex/hertz-board.git
cd hertz-board
```

2. Initialize the project:
```bash
make init
```

This command will:
- Create `.env` file from `.env.example`
- Install backend and frontend dependencies
- Start Docker Compose services
- Run database migrations

#### Running for Development

**Option 1: Using Make (Recommended)**

```bash
# Start all services
make dev-up

# In separate terminals:
make backend-run    # API Gateway
make frontend-run   # Frontend dev server
```

**Option 2: Manual**

1. Start infrastructure:
```bash
docker-compose up -d
```

2. Start backend:
```bash
cd backend
go run cmd/api-gateway/main.go
```

3. Start frontend:
```bash
cd frontend
npm run dev
# or with pnpm
pnpm run dev
```

#### Available Services

After startup, the following services will be available:

- **Frontend**: http://localhost:5173
- **API Gateway**: http://localhost:8080
- **MinIO Console**: http://localhost:9001 (hertzboard / hertzboard_minio_password)
- **MailHog UI**: http://localhost:8025
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **ClickHouse**: localhost:8123
- **NATS**: localhost:4222

### Project Structure

```
hertzboard/
├── backend/              # Go backend services
│   ├── cmd/             # Service entry points
│   ├── internal/        # Internal logic
│   ├── configs/         # Configuration files
│   └── idl/             # IDL files (Protobuf, Thrift)
├── frontend/            # SvelteKit frontend
│   └── src/            
│       ├── lib/         # Components, stores, utils
│       └── routes/      # Application pages
├── nginx/               # Nginx configuration
├── docs/                # Documentation
└── deploy/              # Deployment configs
```

### Available Commands

#### Make Commands

```bash
make help              # Show all available commands
make dev-up            # Start development environment
make dev-down          # Stop development environment
make dev-logs          # Show logs from all services
make install           # Install all dependencies
make migrate           # Run database migrations
make test              # Run all tests
make lint              # Run linters
make build             # Build all services
make clean             # Clean build artifacts
```

#### Backend Commands

```bash
cd backend
go run cmd/api-gateway/main.go    # Run API Gateway (includes WebSocket)
go test -v ./...                  # Run tests
```

#### Frontend Commands

```bash
cd frontend
npm run dev           # or: pnpm run dev
npm run build         # Production build
npm run preview       # Preview production build
npm test              # Run tests
npm run lint          # Run linter
npm run format        # Format code
```

### Testing

```bash
make test              # All tests
make backend-test      # Backend tests only
make frontend-test     # Frontend tests only
make test-coverage     # Tests with coverage report
```

### Deployment

See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed production deployment instructions.

### Development

For contribution guidelines, code style, and commit conventions, see [CONTRIBUTING.md](CONTRIBUTING.md).

### Documentation

Detailed documentation is available in the `docs/` folder:

- [Architecture](docs/architecture/overview.md)
- [API Documentation](docs/api/rest-api.md)
- [WebSocket Protocol](docs/api/websocket-api.md)
- [Development Guide](docs/development/setup.md)

### License

GPL-3.0 License - see [LICENSE](LICENSE) file for details.

---

<a name="ru"></a>

## 🇷🇺 Русский

> Платформа для совместной работы - Бесконечный canvas для создания заметок, рисунков, диаграмм и организации информации в свободной форме.

### Особенности

- **Бесконечный Canvas** - Свободное размещение элементов в 2D пространстве
- **Real-time Коллаборация** - Работайте вместе с командой в реальном времени
- **CRDT-based Синхронизация** - Бесконфликтная синхронизация изменений
- **Богатый набор инструментов** - Текст, фигуры, изображения, рисование от руки
- **Кросс-платформенность** - Web, Desktop (Tauri), Mobile (Capacitor)

### Технологический Стек

#### Backend
- **Hertz v0.10.3** - HTTP/WebSocket сервер
- **Kitex** - gRPC микросервисы
- **PostgreSQL** - Основная база данных
- **Redis** - Кэш и pub/sub
- **MinIO** - Хранилище файлов
- **ClickHouse** - Аналитика
- **NATS** - Очереди сообщений

#### Frontend
- **SvelteKit 2.50+** - Фреймворк
- **TypeScript 5.9+** - Типизация
- **Tailwind CSS 4.1+** - Стилизация
- **Yjs 13.6+** - CRDT синхронизация
- **TipTap 3.18+** - Rich text редактор
- **Vite 7.3+** - Инструмент сборки

### Быстрый старт

#### Требования

- **Go 1.25.6+**
- **Node.js 20+** (pnpm или npm)
- **Docker & Docker Compose**
- **Make** (опционально)

#### Установка

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

#### Запуск для разработки

**Вариант 1: Используя Make (рекомендуется)**

```bash
# Запустить все сервисы
make dev-up

# В отдельных терминалах:
make backend-run    # API Gateway
make frontend-run   # Frontend dev server
```

**Вариант 2: Вручную**

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
# или с pnpm
pnpm run dev
```

#### Доступные сервисы

После запуска будут доступны:

- **Frontend**: http://localhost:5173
- **API Gateway**: http://localhost:8080
- **MinIO Console**: http://localhost:9001 (hertzboard / hertzboard_minio_password)
- **MailHog UI**: http://localhost:8025
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **ClickHouse**: localhost:8123
- **NATS**: localhost:4222

### Структура проекта

```
hertzboard/
├── backend/              # Go backend сервисы
│   ├── cmd/             # Точки входа для сервисов
│   ├── internal/        # Внутренняя логика
│   ├── configs/         # Конфигурационные файлы
│   └── idl/             # IDL файлы (Protobuf, Thrift)
├── frontend/            # SvelteKit frontend
│   └── src/            
│       ├── lib/         # Компоненты, stores, utils
│       └── routes/      # Страницы приложения
├── nginx/               # Конфигурация Nginx
├── docs/                # Документация
└── deploy/              # Конфигурации для развертывания
```

### Доступные команды

#### Make команды

```bash
make help              # Показать все доступные команды
make dev-up            # Запустить dev окружение
make dev-down          # Остановить dev окружение
make dev-logs          # Показать логи всех сервисов
make install           # Установить зависимости
make migrate           # Запустить миграции базы данных
make test              # Запустить все тесты
make lint              # Запустить линтеры
make build             # Собрать все сервисы
make clean             # Очистить build артефакты
```

#### Backend команды

```bash
cd backend
go run cmd/api-gateway/main.go    # Запустить API Gateway (включает WebSocket)
go test -v ./...                  # Запустить тесты
```

#### Frontend команды

```bash
cd frontend
npm run dev           # или: pnpm run dev
npm run build         # Production сборка
npm run preview       # Превью production сборки
npm test              # Запустить тесты
npm run lint          # Запустить линтер
npm run format        # Форматировать код
```

### Тестирование

```bash
make test              # Все тесты
make backend-test      # Только backend тесты
make frontend-test     # Только frontend тесты
make test-coverage     # Тесты с отчетом покрытия
```

### Деплой

Подробные инструкции по развертыванию в production см. в [DEPLOYMENT.md](DEPLOYMENT.md).

### Разработка

Правила внесения вклада, стиль кода и соглашения о коммитах см. в [CONTRIBUTING.md](CONTRIBUTING.md).

### Документация

Подробная документация доступна в папке `docs/`:

- [Архитектура](docs/architecture/overview.md)
- [API Документация](docs/api/rest-api.md)
- [WebSocket Protocol](docs/api/websocket-api.md)
- [Руководство по разработке](docs/development/setup.md)

### Лицензия

GPL-3.0 License - см. [LICENSE](LICENSE) файл для деталей.

---

<a name="zh"></a>

## 🇨🇳 中文

> 协作工作空间平台 - 用于创建笔记、绘图、图表和自由组织信息的无限画布。

### 特性

- **无限画布** - 在2D空间中自由放置元素
- **实时协作** - 与团队实时协作
- **基于CRDT的同步** - 无冲突的变更同步
- **丰富的工具集** - 文本、形状、图像、手绘
- **跨平台** - Web、桌面端（Tauri）、移动端（Capacitor）

### 技术栈

#### 后端
- **Hertz v0.10.3** - HTTP/WebSocket 服务器
- **Kitex** - gRPC 微服务
- **PostgreSQL** - 主数据库
- **Redis** - 缓存和发布/订阅
- **MinIO** - 文件存储
- **ClickHouse** - 分析
- **NATS** - 消息队列

#### 前端
- **SvelteKit 2.50+** - 框架
- **TypeScript 5.9+** - 类型安全
- **Tailwind CSS 4.1+** - 样式
- **Yjs 13.6+** - CRDT 同步
- **TipTap 3.18+** - 富文本编辑器
- **Vite 7.3+** - 构建工具

### 快速开始

#### 要求

- **Go 1.25.6+**
- **Node.js 20+**（pnpm 或 npm）
- **Docker & Docker Compose**
- **Make**（可选）

#### 安装

1. 克隆仓库：
```bash
git clone https://github.com/bifshteksex/hertz-board.git
cd hertz-board
```

2. 初始化项目：
```bash
make init
```

此命令将：
- 从 `.env.example` 创建 `.env` 文件
- 安装后端和前端依赖项
- 启动 Docker Compose 服务
- 运行数据库迁移

#### 开发运行

**选项1：使用 Make（推荐）**

```bash
# 启动所有服务
make dev-up

# 在单独的终端中：
make backend-run    # API 网关
make frontend-run   # 前端开发服务器
```

**选项2：手动**

1. 启动基础设施：
```bash
docker-compose up -d
```

2. 启动后端：
```bash
cd backend
go run cmd/api-gateway/main.go
```

3. 启动前端：
```bash
cd frontend
npm run dev
# 或使用 pnpm
pnpm run dev
```

#### 可用服务

启动后，以下服务将可用：

- **前端**: http://localhost:5173
- **API 网关**: http://localhost:8080
- **MinIO 控制台**: http://localhost:9001 (hertzboard / hertzboard_minio_password)
- **MailHog UI**: http://localhost:8025
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **ClickHouse**: localhost:8123
- **NATS**: localhost:4222

### 项目结构

```
hertzboard/
├── backend/              # Go 后端服务
│   ├── cmd/             # 服务入口点
│   ├── internal/        # 内部逻辑
│   ├── configs/         # 配置文件
│   └── idl/             # IDL 文件（Protobuf、Thrift）
├── frontend/            # SvelteKit 前端
│   └── src/            
│       ├── lib/         # 组件、stores、工具
│       └── routes/      # 应用程序页面
├── nginx/               # Nginx 配置
├── docs/                # 文档
└── deploy/              # 部署配置
```

### 可用命令

#### Make 命令

```bash
make help              # 显示所有可用命令
make dev-up            # 启动开发环境
make dev-down          # 停止开发环境
make dev-logs          # 显示所有服务的日志
make install           # 安装所有依赖项
make migrate           # 运行数据库迁移
make test              # 运行所有测试
make lint              # 运行代码检查器
make build             # 构建所有服务
make clean             # 清理构建产物
```

#### 后端命令

```bash
cd backend
go run cmd/api-gateway/main.go    # 运行 API 网关（包括 WebSocket）
go test -v ./...                  # 运行测试
```

#### 前端命令

```bash
cd frontend
npm run dev           # 或：pnpm run dev
npm run build         # 生产构建
npm run preview       # 预览生产构建
npm test              # 运行测试
npm run lint          # 运行代码检查器
npm run format        # 格式化代码
```

### 测试

```bash
make test              # 所有测试
make backend-test      # 仅后端测试
make frontend-test     # 仅前端测试
make test-coverage     # 带覆盖率报告的测试
```

### 部署

有关详细的生产部署说明，请参阅 [DEPLOYMENT.md](DEPLOYMENT.md)。

### 开发

有关贡献指南、代码风格和提交约定，请参阅 [CONTRIBUTING.md](CONTRIBUTING.md)。

### 文档

详细文档可在 `docs/` 文件夹中找到：

- [架构](docs/architecture/overview.md)
- [API 文档](docs/api/rest-api.md)
- [WebSocket 协议](docs/api/websocket-api.md)
- [开发指南](docs/development/setup.md)

### 许可证

GPL-3.0 License - 有关详细信息，请参阅 [LICENSE](LICENSE) 文件。

---

Made with ❤️ by Roman Shangin
