# Development Setup Guide

## Prerequisites

### Required Software

1. **Go** (1.22 or later)
   ```bash
   # Download from https://go.dev/dl/
   go version  # Verify installation
   ```

2. **Node.js** (20 or later)
   ```bash
   # Download from https://nodejs.org/
   node --version
   npm --version
   ```

3. **Docker & Docker Compose**
   ```bash
   # Download from https://docs.docker.com/get-docker/
   docker --version
   docker-compose --version
   ```

4. **Make** (optional but recommended)
   ```bash
   # macOS: pre-installed
   # Linux: sudo apt-get install build-essential
   # Windows: install via chocolatey or use WSL
   make --version
   ```

### Optional Tools

- **Air** - Hot reload for Go
  ```bash
  go install github.com/cosmtrek/air@latest
  ```

- **golangci-lint** - Go linter
  ```bash
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```

## Project Setup

### 1. Clone Repository

```bash
git clone https://github.com/yourusername/hertzboard.git
cd hertzboard
```

### 2. Environment Configuration

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` file with your local settings:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=hertzboard
DB_USER=hertzboard
DB_PASSWORD=hertzboard_dev_password

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=hertzboard_redis_password

# JWT Secret (generate a secure random string)
JWT_SECRET=your-super-secret-jwt-key-change-this

# OAuth (optional for development)
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
```

### 3. Quick Start (using Make)

```bash
make init
```

This command will:
- Install backend dependencies
- Install frontend dependencies
- Start Docker containers
- Run database migrations

### 4. Manual Setup

If you prefer to set up manually:

#### Start Infrastructure Services

```bash
docker-compose up -d
```

Wait for services to be ready (check with `docker-compose ps`).

#### Install Backend Dependencies

```bash
cd backend
go mod download
go mod tidy
```

#### Install Frontend Dependencies

```bash
cd frontend
npm install
```

#### Run Database Migrations

```bash
# Using Make
make migrate

# Or manually
cd backend
# TODO: Add migration command when implemented
```

## Running the Application

### Option 1: Using Make Commands

```bash
# Start all infrastructure services
make dev-up

# In separate terminals:

# Terminal 1: Run API Gateway
make backend-run

# Terminal 2: Run Frontend
make frontend-run

# Optional - Terminal 3: Run WebSocket Server
cd backend && go run cmd/ws-server/main.go
```

### Option 2: Manual Execution

#### Backend (API Gateway)

```bash
cd backend
go run cmd/api-gateway/main.go
```

#### Backend (WebSocket Server)

```bash
cd backend
go run cmd/ws-server/main.go
```

#### Frontend

```bash
cd frontend
npm run dev
```

## Accessing Services

Once everything is running, you can access:

- **Frontend**: http://localhost:5173
- **API Gateway**: http://localhost:8080
- **API Health Check**: http://localhost:8080/health
- **WebSocket Server**: ws://localhost:8081
- **MinIO Console**: http://localhost:9001
  - Username: `hertzboard`
  - Password: `hertzboard_minio_password`
- **MailHog (Email Testing)**: http://localhost:8025

## Development Workflow

### Backend Development

1. **Make changes** to Go files
2. **Restart server** (or use Air for hot reload)
3. **Run tests**:
   ```bash
   cd backend
   go test ./...
   ```

#### Using Air for Hot Reload

Create `.air.toml` in backend directory:

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/api-gateway/main.go"
  delay = 1000
  exclude_dir = ["tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_error = true

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false
```

Then run:

```bash
cd backend
air
```

### Frontend Development

SvelteKit provides hot module replacement (HMR) by default.

1. **Make changes** to Svelte/TypeScript files
2. **Browser auto-refreshes** with changes
3. **Run tests**:
   ```bash
   cd frontend
   npm run test
   ```

### Code Quality

#### Linting

```bash
# Backend
make backend-lint
# or
cd backend && golangci-lint run

# Frontend
make frontend-lint
# or
cd frontend && npm run lint
```

#### Formatting

```bash
# Backend
cd backend && go fmt ./...

# Frontend
cd frontend && npm run format
```

#### Type Checking (Frontend)

```bash
cd frontend
npm run check
```

## Database Management

### Running Migrations

```bash
make migrate
```

### Creating New Migration

```bash
make migrate-create
# Enter migration name when prompted
```

### Rolling Back Migration

```bash
make migrate-down
```

### Accessing PostgreSQL

```bash
docker exec -it hertz-board-postgres psql -U hertzboard -d hertzboard
```

### Common SQL Commands

```sql
-- List all tables
\dt

-- Describe table
\d users

-- Select from table
SELECT * FROM users;
```

## Testing

### Backend Tests

```bash
# Run all tests
make backend-test

# Run with coverage
make test-coverage

# Run specific package
cd backend
go test ./internal/domain/user/...
```

### Frontend Tests

```bash
# Run unit tests
make frontend-test

# Run tests in watch mode
cd frontend
npm run test
```

## Troubleshooting

### Port Already in Use

If ports are already in use, you can change them in:
- `.env` file for application ports
- `docker-compose.yml` for service ports

### Database Connection Issues

1. Check if PostgreSQL is running:
   ```bash
   docker-compose ps postgres
   ```

2. Check logs:
   ```bash
   docker-compose logs postgres
   ```

3. Verify connection:
   ```bash
   docker exec hertz-board-postgres pg_isready
   ```

### Redis Connection Issues

1. Check if Redis is running:
   ```bash
   docker-compose ps redis
   ```

2. Test connection:
   ```bash
   docker exec hertz-board-redis redis-cli ping
   ```

### Node Modules Issues

```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Go Modules Issues

```bash
cd backend
go clean -modcache
go mod download
```

## Useful Commands

### Docker

```bash
# View logs
docker-compose logs -f [service_name]

# Restart service
docker-compose restart [service_name]

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Rebuild services
docker-compose up -d --build
```

### Database

```bash
# Backup database
docker exec hertz-board-postgres pg_dump -U hertzboard hertzboard > backup.sql

# Restore database
cat backup.sql | docker exec -i hertz-board-postgres psql -U hertzboard hertzboard
```

## Next Steps

- Read [Architecture Overview](../architecture/overview.md)
- Check [API Documentation](../api/rest-api.md)
- Review [Coding Standards](./coding-standards.md)
- Explore the codebase structure
