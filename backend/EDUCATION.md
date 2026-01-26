# HertzBoard Backend: Полное руководство по Go

Это обучающее руководство создано на основе реального кода твоего проекта HertzBoard. Здесь ты найдёшь объяснение всех концепций Go с примерами прямо из твоего кода и сравнениями с PHP.

---

## Содержание

1. [Основы Go: Структура проекта](#1-основы-go-структура-проекта)
2. [Пакеты и импорты](#2-пакеты-и-импорты)
3. [Типы данных и переменные](#3-типы-данных-и-переменные)
4. [Структуры (Structs) — аналог классов](#4-структуры-structs--аналог-классов)
5. [Методы и получатели (Receivers)](#5-методы-и-получатели-receivers)
6. [Указатели](#6-указатели)
7. [Интерфейсы](#7-интерфейсы)
8. [Обработка ошибок](#8-обработка-ошибок)
9. [Горутины и конкурентность](#9-горутины-и-конкурентность)
10. [Каналы (Channels)](#10-каналы-channels)
11. [Context — управление жизненным циклом](#11-context--управление-жизненным-циклом)
12. [HTTP-фреймворк Hertz](#12-http-фреймворк-hertz)
13. [Middleware — промежуточное ПО](#13-middleware--промежуточное-по)
14. [Работа с базой данных (pgx)](#14-работа-с-базой-данных-pgx)
15. [JSON и сериализация](#15-json-и-сериализация)
16. [JWT-аутентификация](#16-jwt-аутентификация)
17. [WebSocket и реальное время](#17-websocket-и-реальное-время)
18. [CRDT — синхронизация данных](#18-crdt--синхронизация-данных)
19. [Архитектура проекта](#19-архитектура-проекта)
20. [Тест на проверку знаний](#20-тест-на-проверку-знаний)

---

## 1. Основы Go: Структура проекта

### Структура директорий твоего проекта

```
backend/
├── cmd/                    # Точки входа (main.go файлы)
│   ├── api-gateway/        # Основной REST API сервер
│   └── ws-server/          # WebSocket сервер (заглушка)
├── internal/               # Приватный код (не импортируется извне)
│   ├── config/             # Конфигурация
│   ├── database/           # Подключения к БД
│   ├── handler/            # HTTP обработчики (Controllers в PHP)
│   ├── middleware/         # Middleware
│   ├── models/             # Модели данных (Entities)
│   ├── repository/         # Работа с БД (Repository pattern)
│   ├── router/             # Маршрутизация
│   └── service/            # Бизнес-логика (Services)
├── migrations/             # SQL миграции
├── configs/                # Конфигурационные файлы
└── go.mod                  # Зависимости (аналог composer.json)
```

### Сравнение с PHP (Laravel/Symfony)

| Go (твой проект)       | PHP (Laravel)        | PHP (Symfony)           |
|------------------------|----------------------|-------------------------|
| `cmd/api-gateway/`     | `public/index.php`   | `public/index.php`      |
| `internal/handler/`    | `app/Http/Controllers` | `src/Controller`      |
| `internal/service/`    | `app/Services`       | `src/Service`           |
| `internal/repository/` | `app/Repositories`   | `src/Repository`        |
| `internal/models/`     | `app/Models`         | `src/Entity`            |
| `internal/middleware/` | `app/Http/Middleware` | `src/EventSubscriber`  |
| `go.mod`               | `composer.json`      | `composer.json`         |

### Особенность: Папка `internal/`

В Go папка `internal/` имеет **специальное значение**: код внутри неё не может быть импортирован другими проектами. Это встроенная в язык защита приватного кода.

```go
// Это работает (внутри проекта):
import "github.com/bifshteksex/hertz-board/internal/service"

// Это НЕ работает (из другого проекта):
// import "github.com/bifshteksex/hertz-board/internal/service" // Ошибка компиляции!
```

---

## 2. Пакеты и импорты

### Объявление пакета

Каждый файл Go начинается с объявления пакета:

```go
// internal/service/auth_service.go
package service  // Имя пакета — последняя часть пути директории
```

**Сравнение с PHP:**
```php
// PHP
namespace App\Service;
```

### Импорты

Из твоего файла `cmd/api-gateway/main.go`:

```go
import (
    "context"           // Стандартная библиотека
    "fmt"               // Форматирование строк
    "log"               // Логирование
    "os"                // Работа с ОС
    "os/signal"         // Обработка сигналов
    "syscall"           // Системные вызовы
    "time"              // Работа со временем

    // Внешние зависимости (через go.mod)
    "github.com/cloudwego/hertz/pkg/app/server"

    // Локальные пакеты проекта
    "github.com/bifshteksex/hertz-board/internal/config"
    "github.com/bifshteksex/hertz-board/internal/database"
    "github.com/bifshteksex/hertz-board/internal/handler"
    "github.com/bifshteksex/hertz-board/internal/repository"
    "github.com/bifshteksex/hertz-board/internal/router"
    "github.com/bifshteksex/hertz-board/internal/service"
)
```

### Правила именования импортов

```go
import (
    // Стандартная библиотека — одно слово
    "fmt"
    "context"
    
    // Внешние пакеты — URL репозитория
    "github.com/google/uuid"
    
    // Пакет с alias (когда имена конфликтуют или длинные)
    pgx "github.com/jackc/pgx/v5"
    
    // Версионирование в пути
    "github.com/jackc/pgx/v5"  // v5 — мажорная версия
)
```

**Сравнение с PHP:**
```php
// PHP использует use вместо import
use App\Service\AuthService;
use Symfony\Component\HttpFoundation\Request;
```

---

## 3. Типы данных и переменные

### Базовые типы

Из твоего файла `internal/config/config.go`:

```go
type AppConfig struct {
    Name  string `yaml:"name"`   // Строка
    Env   string `yaml:"env"`    // Строка
    Port  int    `yaml:"port"`   // Целое число (int32 или int64 в зависимости от платформы)
    Debug bool   `yaml:"debug"`  // Булево значение
}
```

### Таблица типов: Go vs PHP

| Go            | PHP           | Описание                        |
|---------------|---------------|----------------------------------|
| `string`      | `string`      | Строка (UTF-8)                   |
| `int`         | `int`         | Целое число                      |
| `int64`       | —             | 64-битное целое                  |
| `float64`     | `float`       | Число с плавающей точкой         |
| `bool`        | `bool`        | Булево                           |
| `[]byte`      | `string`      | Массив байтов (бинарные данные)  |
| `[]string`    | `array`       | Срез (динамический массив)       |
| `map[K]V`     | `array`       | Ассоциативный массив             |
| `*T`          | —             | Указатель на тип T               |
| `interface{}` | `mixed`       | Любой тип                        |

### Объявление переменных

```go
// Полное объявление
var name string = "HertzBoard"

// Краткое объявление (внутри функций)
name := "HertzBoard"  // Тип выводится автоматически

// Константы
const (
    maxRequestBodySizeMB   = 10
    shutdownTimeoutSeconds = 5
    bytesInMB              = 1024 * 1024
    defaultConfigPath      = "configs/config.yaml"
)
```

**Сравнение с PHP:**
```php
// PHP
$name = "HertzBoard";  // Без указания типа
const MAX_SIZE = 10;   // Константа
```

### Нулевые значения (Zero Values)

В Go переменные автоматически инициализируются "нулевыми значениями":

```go
var s string   // "" (пустая строка)
var i int      // 0
var b bool     // false
var p *int     // nil (указатель никуда не указывает)
var m map[string]int  // nil (карта не инициализирована!)
```

**Важно!** Nil-карту нельзя использовать для записи:
```go
var m map[string]int
m["key"] = 1  // PANIC! Runtime error

// Правильно:
m := make(map[string]int)
m["key"] = 1  // OK
```

---

## 4. Структуры (Structs) — аналог классов

### Определение структуры

Из твоего файла `internal/models/user.go`:

```go
type User struct {
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
    PasswordHash  *string   `json:"-" db:"password_hash"`        // Указатель — может быть nil
    AvatarURL     *string   `json:"avatar_url,omitempty" db:"avatar_url"`
    ProviderID    *string   `json:"-" db:"provider_id"`
    Email         string    `json:"email" db:"email"`
    Name          string    `json:"name" db:"name"`
    Provider      string    `json:"provider" db:"provider"`
    ID            uuid.UUID `json:"id" db:"id"`
    EmailVerified bool      `json:"email_verified" db:"email_verified"`
}
```

### Теги структур (Struct Tags)

Теги — это метаданные полей, используемые библиотеками:

```go
type User struct {
    Email string `json:"email" db:"email" binding:"required,email"`
    //           ↑ JSON-сериализация  ↑ БД-маппинг  ↑ Валидация
}
```

| Тег | Библиотека | Назначение |
|-----|------------|------------|
| `json:"field_name"` | encoding/json | Имя поля в JSON |
| `json:"-"` | encoding/json | Исключить из JSON |
| `json:",omitempty"` | encoding/json | Пропустить если пусто |
| `db:"column_name"` | sqlx/pgx | Имя колонки в БД |
| `yaml:"field_name"` | yaml.v3 | Имя поля в YAML |
| `binding:"required"` | hertz/gin | Валидация — обязательное |
| `binding:"email"` | hertz/gin | Валидация — email формат |

**Сравнение с PHP (Doctrine):**
```php
// PHP использует аннотации или атрибуты
#[ORM\Column(type: 'string')]
private string $email;
```

### Создание экземпляра структуры

```go
// Способ 1: Литерал структуры
user := User{
    Email: "user@example.com",
    Name:  "John",
}

// Способ 2: Указатель на структуру (используется чаще)
user := &User{
    Email: "user@example.com",
    Name:  "John",
}

// Способ 3: new() — создаёт указатель с нулевыми значениями
user := new(User)  // *User с пустыми полями
```

### Встраивание структур (Embedding)

Из `internal/models/workspace.go`:

```go
// Базовая структура
type Workspace struct {
    ID       uuid.UUID `json:"id"`
    Name     string    `json:"name"`
    OwnerID  uuid.UUID `json:"owner_id"`
    IsPublic bool      `json:"is_public"`
}

// Расширенная структура — встраивает Workspace
type WorkspaceWithRole struct {
    Owner    *User         `json:"owner,omitempty"`
    UserRole WorkspaceRole `json:"user_role"`
    Workspace  // Анонимное встраивание — все поля Workspace доступны напрямую
}
```

Использование:
```go
ws := WorkspaceWithRole{
    Workspace: Workspace{
        ID:   uuid.New(),
        Name: "My Board",
    },
    UserRole: WorkspaceRoleEditor,
}

// Поля Workspace доступны напрямую:
fmt.Println(ws.Name)  // "My Board" — не ws.Workspace.Name
```

**Сравнение с PHP (наследование):**
```php
// В PHP использовалось бы наследование
class WorkspaceWithRole extends Workspace {
    public User $owner;
    public string $userRole;
}
```

---

## 5. Методы и получатели (Receivers)

### Методы с получателем

В Go методы привязываются к типам через "получатель" (receiver):

```go
// internal/models/canvas.go

// Метод для типа ElementType
func (t ElementType) Valid() bool {
    switch t {
    case ElementTypeText, ElementTypeShape, ElementTypeImage:
        return true
    }
    return false
}

// Метод для типа CanvasElement
func (e *CanvasElement) ToResponse() ElementResponse {
    return ElementResponse{
        ID:          e.ID,
        WorkspaceID: e.WorkspaceID,
        ElementType: e.ElementType,
        // ...
    }
}
```

### Значение vs Указатель в получателе

```go
// Получатель-значение: копирует структуру (безопасно, но медленнее для больших структур)
func (t ElementType) Valid() bool {
    // t — это копия
}

// Получатель-указатель: работает с оригиналом (быстрее, может модифицировать)
func (e *CanvasElement) ToResponse() ElementResponse {
    // e — это указатель на оригинал
}
```

**Правило:** Используй указатель-получатель если:
- Метод модифицирует структуру
- Структура большая (экономия памяти)
- Нужна консистентность (все методы одного типа должны использовать одинаковый тип получателя)

**Сравнение с PHP:**
```php
// PHP методы всегда работают с $this
class ElementType {
    public function isValid(): bool {
        return in_array($this->value, ['text', 'shape', 'image']);
    }
}
```

### Конструкторы (Factory Functions)

В Go нет конструкторов как в PHP. Вместо них используются функции-фабрики:

```go
// internal/service/auth_service.go

// NewAuthService — "конструктор" для AuthService
func NewAuthService(userRepo *repository.UserRepository, jwtService *JWTService) *AuthService {
    return &AuthService{
        userRepo:   userRepo,
        jwtService: jwtService,
    }
}
```

**Сравнение с PHP:**
```php
// PHP конструктор
class AuthService {
    public function __construct(
        private UserRepository $userRepo,
        private JWTService $jwtService
    ) {}
}
```

---

## 6. Указатели

### Что такое указатель?

Указатель хранит **адрес памяти**, где находится значение:

```go
var x int = 10
var p *int = &x  // p указывает на x

fmt.Println(x)   // 10 — значение
fmt.Println(&x)  // 0xc000012080 — адрес в памяти
fmt.Println(p)   // 0xc000012080 — то же самое
fmt.Println(*p)  // 10 — разыменование (получить значение по адресу)

*p = 20          // Изменяем значение по адресу
fmt.Println(x)   // 20 — x изменился!
```

### Зачем нужны указатели?

1. **Модификация аргументов функции:**

```go
// Без указателя — работаем с копией
func double(x int) {
    x = x * 2  // Изменяем копию, оригинал не меняется
}

// С указателем — работаем с оригиналом
func double(x *int) {
    *x = *x * 2  // Изменяем оригинал
}

n := 5
double(&n)  // Передаём адрес
fmt.Println(n)  // 10
```

2. **Опциональные поля (nil = отсутствует):**

Из `internal/models/user.go`:
```go
type User struct {
    PasswordHash *string  // nil если OAuth-пользователь
    AvatarURL    *string  // nil если не установлен
}

// Проверка
if user.PasswordHash == nil {
    return errors.New("user registered with OAuth, no password")
}
```

3. **Экономия памяти (большие структуры):**

```go
// Плохо: копируется вся структура при каждом вызове
func processUser(u User) { ... }

// Хорошо: передаётся только адрес (8 байт)
func processUser(u *User) { ... }
```

### Указатели в твоём проекте

Из `internal/models/workspace.go`:
```go
type UpdateWorkspaceRequest struct {
    Name         *string                `json:"name,omitempty"`
    Description  *string                `json:"description,omitempty"`
    IsPublic     *bool                  `json:"is_public,omitempty"`
    ThumbnailURL *string                `json:"thumbnail_url,omitempty"`
    Settings     map[string]interface{} `json:"settings,omitempty"`
}
```

Использование в `internal/service/workspace_service.go`:
```go
func (s *WorkspaceService) UpdateWorkspace(ctx context.Context, workspaceID uuid.UUID, req *UpdateWorkspaceRequest) (*Workspace, error) {
    // Проверяем, указано ли поле (не nil)
    if req.Name != nil {
        workspace.Name = *req.Name  // Разыменование
    }
    if req.IsPublic != nil {
        workspace.IsPublic = *req.IsPublic
    }
    // ...
}
```

**Паттерн:** Указатели в DTO позволяют различать "поле не передано" (nil) и "поле передано пустым" ("").

**Сравнение с PHP:**
В PHP нет указателей, но есть ссылки:
```php
function double(&$x) {  // & — передача по ссылке
    $x = $x * 2;
}
```

---

## 7. Интерфейсы

### Неявная реализация интерфейсов

В Go интерфейсы реализуются **неявно** — не нужно писать `implements`:

```go
// Определение интерфейса
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Любой тип с методом Read автоматически реализует Reader
type MyFile struct {}

func (f *MyFile) Read(p []byte) (n int, err error) {
    // Реализация
    return 0, nil
}

// MyFile теперь реализует Reader без явного объявления
var r Reader = &MyFile{}  // OK!
```

### Интерфейсы в твоём проекте

Из `internal/models/canvas.go` — реализация интерфейсов для работы с БД:

```go
// ElementData реализует sql.Scanner и driver.Valuer для JSONB

// Scan — для чтения из БД (реализует sql.Scanner)
func (e *ElementData) Scan(value interface{}) error {
    if value == nil {
        *e = make(ElementData)
        return nil
    }
    
    bytes, ok := value.([]byte)
    if !ok {
        return fmt.Errorf("failed to scan ElementData: unexpected type %T", value)
    }
    return json.Unmarshal(bytes, e)
}

// Value — для записи в БД (реализует driver.Valuer)
func (e ElementData) Value() (driver.Value, error) {
    if e == nil {
        return "{}", nil
    }
    return json.Marshal(e)
}
```

### Пустой интерфейс `interface{}`

`interface{}` может содержать значение любого типа (аналог `mixed` в PHP):

```go
// Из internal/models/workspace.go
type Workspace struct {
    Settings map[string]interface{} `json:"settings"`  // Любые настройки
}

// Использование
ws.Settings["theme"] = "dark"      // string
ws.Settings["fontSize"] = 14       // int
ws.Settings["features"] = []string{"a", "b"}  // []string
```

### Type Assertion (Приведение типа)

```go
// Из internal/service/crdt_service.go
func (s *CRDTService) applyCreate(op *models.OperationPayload) error {
    var elementData map[string]interface{}
    // ...
    
    // Извлечение значения с проверкой типа
    elementType, ok := elementData["type"].(string)
    if !ok {
        // Не строка или отсутствует
    }
    
    posX, _ := elementData["pos_x"].(float64)  // JSON числа всегда float64
    zIndex, _ := elementData["z_index"].(float64)
}
```

**Сравнение с PHP:**
```php
// PHP использует mixed и instanceof
function process(mixed $value): void {
    if ($value instanceof User) {
        // ...
    }
}
```

---

## 8. Обработка ошибок

### Ошибки как значения

В Go ошибки — это обычные значения, а не исключения:

```go
// Функция возвращает результат И ошибку
func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    // ...
    return &cfg, nil
}
```

### Паттерн обработки ошибок

Из `cmd/api-gateway/main.go`:

```go
// Типичный паттерн — проверка сразу после вызова
cfg, err := config.Load(configPath)
if err != nil {
    log.Fatalf("Failed to load config: %v", err)  // Критическая ошибка
}

dbPool, err := database.NewPostgresPool(&cfg.Database)
if err != nil {
    log.Fatalf("Failed to connect to PostgreSQL: %v", err)
}
defer database.ClosePostgresPool(dbPool)  // Отложенная очистка
```

### Оборачивание ошибок (Error Wrapping)

Из `internal/service/auth_service.go`:

```go
func (s *AuthService) Register(ctx context.Context, req *models.CreateUserRequest) (*models.AuthResponse, error) {
    existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err != nil {
        // %w — оборачивает оригинальную ошибку (сохраняет цепочку)
        return nil, fmt.Errorf("failed to check existing user: %w", err)
    }
    
    if existingUser != nil {
        // Создание новой ошибки
        return nil, fmt.Errorf("user with email %s already exists", req.Email)
    }
    // ...
}
```

### Проверка типа ошибки

```go
import "errors"

// errors.Is — проверяет, содержит ли цепочка конкретную ошибку
if errors.Is(err, sql.ErrNoRows) {
    return nil, nil  // Не найдено — это не ошибка
}

// errors.As — извлекает ошибку определённого типа
var pgErr *pgconn.PgError
if errors.As(err, &pgErr) {
    if pgErr.Code == "23505" {  // Unique violation
        return nil, fmt.Errorf("email already exists")
    }
}
```

**Сравнение с PHP (исключения):**
```php
// PHP использует try/catch
try {
    $user = $this->userRepo->findByEmail($email);
} catch (NotFoundException $e) {
    return null;
} catch (\Exception $e) {
    throw new \RuntimeException("Failed to find user: " . $e->getMessage(), 0, $e);
}
```

### Интересный факт из твоего проекта

В `internal/handler/auth_handler.go` используется паттерн безопасного ответа:

```go
func (h *AuthHandler) ForgotPassword(c context.Context, ctx *app.RequestContext) {
    // ...
    _, err := h.authService.ForgotPassword(c, req.Email)
    if err != nil {
        // Не раскрываем, существует ли email (защита от перебора)
        ctx.JSON(consts.StatusOK, map[string]interface{}{
            "message": "If the email exists, a password reset link has been sent",
        })
        return
    }
    // Тот же ответ при успехе
    ctx.JSON(consts.StatusOK, map[string]interface{}{
        "message": "If the email exists, a password reset link has been sent",
    })
}
```

---

## 9. Горутины и конкурентность

### Что такое горутина?

Горутина — это легковесный поток выполнения. Запускается ключевым словом `go`:

```go
// Обычный вызов — выполняется синхронно
doSomething()

// Горутина — выполняется асинхронно (параллельно)
go doSomething()
```

### Горутины в твоём проекте

Из `cmd/api-gateway/main.go`:

```go
// Запуск сервера в горутине
go func() {
    if err := h.Run(); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}()

log.Printf("API Gateway is running on %s", addr)

// Ожидание сигнала завершения (основной поток продолжает работать)
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit  // Блокируется до получения сигнала
```

Из `internal/service/hub.go`:

```go
func NewHub(redisClient *redis.Client) *Hub {
    hub := &Hub{
        rooms: make(map[uuid.UUID]*models.Room),
        redis: redisClient,
        ctx:   context.Background(),
    }

    // Запуск фоновых задач в горутинах
    go hub.cleanupEmptyRooms()    // Периодическая очистка
    go hub.subscribeToRedis()     // Подписка на Redis pub/sub

    return hub
}
```

### Синхронизация: Mutex

Из `internal/service/crdt_service.go`:

```go
type LamportClock struct {
    counter int64
    mu      sync.Mutex  // Мьютекс для защиты counter
}

func (lc *LamportClock) Tick() int64 {
    lc.mu.Lock()         // Захватываем мьютекс
    defer lc.mu.Unlock() // Освобождаем при выходе из функции
    lc.counter++
    return lc.counter
}
```

### RWMutex — для оптимизации чтения

Из `internal/service/hub.go`:

```go
type Hub struct {
    rooms map[uuid.UUID]*models.Room
    mu    sync.RWMutex  // RWMutex позволяет множественное чтение
}

func (h *Hub) Register(client *models.Client) {
    h.mu.Lock()  // Эксклюзивная блокировка для записи
    defer h.mu.Unlock()
    // Модификация rooms...
}

func (h *Hub) GetRoomStats(workspaceID uuid.UUID) (int, bool) {
    h.mu.RLock()  // Блокировка только для чтения (можно параллельно)
    defer h.mu.RUnlock()
    // Только чтение rooms...
}
```

**Сравнение с PHP:**
```php
// PHP — однопоточный, мьютексы не нужны
// Для параллелизма используются:
// - Symfony Messenger (очереди)
// - ReactPHP/Swoole (event loop)
// - Процессы (pcntl_fork)
```

---

## 10. Каналы (Channels)

### Что такое канал?

Канал — это типизированный "трубопровод" для передачи данных между горутинами:

```go
// Создание канала
ch := make(chan int)      // Небуферизованный
ch := make(chan int, 10)  // Буферизованный (вместимость 10)

// Отправка в канал
ch <- 42

// Получение из канала
value := <-ch
```

### Каналы в твоём проекте

Из `internal/models/websocket.go`:

```go
type Client struct {
    ID          uuid.UUID
    UserID      uuid.UUID
    WorkspaceID uuid.UUID
    Send        chan *WSMessage  // Канал для отправки сообщений клиенту
}

type Room struct {
    WorkspaceID uuid.UUID
    Clients     map[uuid.UUID]*Client
    Broadcast   chan *WSMessage  // Канал для broadcast-сообщений
    Register    chan *Client     // Канал для регистрации клиентов
    Unregister  chan *Client     // Канал для отключения клиентов
}
```

Из `internal/service/hub.go`:

```go
func (h *Hub) runRoom(room *models.Room) {
    for {
        select {
        case client := <-room.Register:
            // Новый клиент подключился
            room.Clients[client.ID] = client
            
        case client := <-room.Unregister:
            // Клиент отключился
            if _, ok := room.Clients[client.ID]; ok {
                delete(room.Clients, client.ID)
                close(client.Send)  // Закрываем канал клиента
            }
            
        case message := <-room.Broadcast:
            // Рассылаем сообщение всем клиентам
            h.broadcastToRoomClients(room, message, uuid.Nil)
        }
    }
}
```

### Паттерн: Select для мультиплексирования

`select` позволяет ждать на нескольких каналах одновременно:

```go
select {
case msg := <-ch1:
    // Получили из ch1
case msg := <-ch2:
    // Получили из ch2
case ch3 <- value:
    // Отправили в ch3
default:
    // Если все каналы заблокированы (необязательно)
}
```

### Паттерн: Graceful Shutdown

Из `cmd/api-gateway/main.go`:

```go
// Создаём канал для сигналов ОС
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

// Блокируемся до получения сигнала
<-quit  // Ctrl+C или kill

log.Println("Shutting down server...")

// Graceful shutdown с таймаутом
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := h.Shutdown(ctx); err != nil {
    log.Fatalf("Server forced to shutdown: %v", err)
}
```

---

## 11. Context — управление жизненным циклом

### Зачем нужен Context?

Context используется для:
1. **Отмены операций** (cancellation)
2. **Таймаутов** (deadlines)
3. **Передачи request-scoped данных**

### Context в твоём проекте

Из `internal/handler/auth_handler.go`:

```go
// c context.Context — первый параметр, передаётся от Hertz
func (h *AuthHandler) Register(c context.Context, ctx *app.RequestContext) {
    // Передаём context в сервис
    resp, err := h.authService.Register(c, &req)
}
```

Из `internal/repository/user_repository.go`:

```go
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
    // ctx передаётся в запрос к БД
    err := r.db.QueryRow(ctx, query, ...).Scan(...)
    // Если ctx отменён — запрос прервётся
}
```

### Создание Context с таймаутом

Из `internal/database/postgres.go`:

```go
func NewPostgresPool(cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
    // ...
    
    // Context с таймаутом 5 секунд
    ctx, cancel := context.WithTimeout(context.Background(), defaultPingTimeout)
    defer cancel()  // Освобождаем ресурсы

    // Ping прервётся, если превысит 5 секунд
    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
}
```

### Передача данных через Context

Из `internal/middleware/auth.go`:

```go
func Auth(jwtService *service.JWTService) app.HandlerFunc {
    return func(c context.Context, ctx *app.RequestContext) {
        // ...валидация токена...
        
        // Сохраняем данные пользователя в context
        ctx.Set("user_id", claims.UserID)
        ctx.Set("user_email", claims.Email)
        
        ctx.Next(c)  // Передаём управление следующему обработчику
    }
}
```

Использование в handler:
```go
func (h *UserHandler) GetProfile(c context.Context, ctx *app.RequestContext) {
    // Извлекаем user_id из context
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(401, map[string]interface{}{"error": "Unauthorized"})
        return
    }
}
```

**Сравнение с PHP:**
```php
// В PHP request-scoped данные хранятся в объекте Request
$request->attributes->get('user_id');

// Или через middleware с атрибутами
$request->getAttribute('user');
```

---

## 12. HTTP-фреймворк Hertz

### Что такое Hertz?

[Hertz](https://www.cloudwego.io/docs/hertz/) — это высокопроизводительный HTTP-фреймворк от CloudWeGo (ByteDance). Аналог Gin/Echo, но быстрее благодаря собственной сетевой библиотеке Netpoll.

**Документация:** https://www.cloudwego.io/docs/hertz/

### Создание сервера

Из `cmd/api-gateway/main.go`:

```go
import "github.com/cloudwego/hertz/pkg/app/server"

// Создание сервера с настройками
h := server.Default(
    server.WithHostPorts(":8080"),
    server.WithMaxRequestBodySize(10 * 1024 * 1024),  // 10 MB
)
```

### Обработчик (Handler)

Из `internal/handler/auth_handler.go`:

```go
// Сигнатура обработчика Hertz
func (h *AuthHandler) Register(c context.Context, ctx *app.RequestContext) {
    // c — стандартный context.Context
    // ctx — контекст запроса Hertz (тело, заголовки, ответ)
    
    // Парсинг JSON тела запроса
    var req models.CreateUserRequest
    if err := ctx.BindAndValidate(&req); err != nil {
        ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
            "error": err.Error(),
        })
        return
    }
    
    // Бизнес-логика
    resp, err := h.authService.Register(c, &req)
    if err != nil {
        ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
            "error": err.Error(),
        })
        return
    }
    
    // Успешный ответ
    ctx.JSON(consts.StatusCreated, resp)
}
```

### Роутинг

Из `internal/router/router.go`:

```go
func Setup(h *server.Hertz, cfg *config.Config, deps *Dependencies) {
    // Глобальный middleware
    h.Use(middleware.Recovery())
    h.Use(middleware.RequestID())
    h.Use(middleware.CORS(&cfg.CORS))
    
    // Health check
    h.GET("/health", healthCheck)
    
    // Группа маршрутов v1
    v1 := h.Group("/api/v1")
    
    // Auth routes (публичные)
    auth := v1.Group("/auth")
    auth.POST("/register", deps.AuthHandler.Register)
    auth.POST("/login", deps.AuthHandler.Login)
    
    // User routes (защищённые)
    users := v1.Group("/users")
    users.Use(middleware.Auth(deps.JWTService))  // Middleware для группы
    users.GET("/me", deps.UserHandler.GetProfile)
    users.PUT("/me", deps.UserHandler.UpdateProfile)
    
    // Workspace routes с проверкой прав
    workspaces := v1.Group("/workspaces")
    workspaces.Use(middleware.Auth(deps.JWTService))
    
    workspaces.GET("/:workspace_id",
        workspaceMiddleware.OptionalWorkspaceAccess(),  // Может быть публичным
        deps.WorkspaceHandler.GetWorkspace,
    )
    
    workspaces.PUT("/:workspace_id",
        workspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleEditor),
        deps.WorkspaceHandler.UpdateWorkspace,
    )
}
```

### Параметры маршрута

```go
// :workspace_id — параметр пути
workspaces.GET("/:workspace_id", handler)

// Получение параметра
func (h *Handler) GetWorkspace(c context.Context, ctx *app.RequestContext) {
    workspaceIDStr := ctx.Param("workspace_id")
    workspaceID, err := uuid.Parse(workspaceIDStr)
}
```

**Сравнение с PHP (Symfony):**
```php
// Symfony routing
#[Route('/api/v1/workspaces/{workspace_id}', methods: ['GET'])]
public function getWorkspace(string $workspace_id): Response
{
    $workspaceId = Uuid::fromString($workspace_id);
}
```

---

## 13. Middleware — промежуточное ПО

### Что такое Middleware?

Middleware — это функции, которые выполняются **до** и/или **после** основного обработчика:

```
Request → Middleware1 → Middleware2 → Handler → Middleware2 → Middleware1 → Response
```

### Middleware в твоём проекте

**1. Recovery — перехват паники:**

Из `internal/middleware/recovery.go`:
```go
func Recovery() app.HandlerFunc {
    return func(c context.Context, ctx *app.RequestContext) {
        defer func() {
            if r := recover(); r != nil {
                // Логируем stack trace
                log.Printf("Panic recovered: %v\n%s", r, debug.Stack())
                // Возвращаем 500
                ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
                    "error": "Internal server error",
                })
            }
        }()
        ctx.Next(c)  // Вызываем следующий обработчик
    }
}
```

**2. Request ID — уникальный ID запроса:**

Из `internal/middleware/request_id.go`:
```go
func RequestID() app.HandlerFunc {
    return func(c context.Context, ctx *app.RequestContext) {
        // Генерируем или берём из заголовка
        requestID := string(ctx.Request.Header.Peek("X-Request-ID"))
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        // Устанавливаем в заголовок ответа
        ctx.Response.Header.Set("X-Request-ID", requestID)
        
        // Сохраняем в context
        ctx.Set("request_id", requestID)
        
        ctx.Next(c)
    }
}
```

**3. Auth — JWT аутентификация:**

Из `internal/middleware/auth.go`:
```go
func Auth(jwtService *service.JWTService) app.HandlerFunc {
    return func(c context.Context, ctx *app.RequestContext) {
        // Извлекаем токен из заголовка
        authHeader := string(ctx.Request.Header.Peek("Authorization"))
        if authHeader == "" {
            ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
                "error": "Authorization header required",
            })
            ctx.Abort()  // Прерываем цепочку!
            return
        }
        
        // Парсим "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
                "error": "Invalid authorization header format",
            })
            ctx.Abort()
            return
        }
        
        // Валидируем JWT
        claims, err := jwtService.ValidateAccessToken(parts[1])
        if err != nil {
            ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
                "error": "Invalid or expired token",
            })
            ctx.Abort()
            return
        }
        
        // Сохраняем данные пользователя
        ctx.Set("user_id", claims.UserID)
        ctx.Set("user_email", claims.Email)
        
        ctx.Next(c)  // Продолжаем цепочку
    }
}
```

**4. Workspace Access — RBAC:**

Из `internal/middleware/workspace.go`:
```go
func (m *WorkspaceMiddleware) RequireWorkspaceAccess(requiredRole models.WorkspaceRole) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // Получаем user_id (установлен Auth middleware)
        userID, exists := c.Get("user_id")
        if !exists {
            c.JSON(http.StatusUnauthorized, map[string]interface{}{
                "error": "Unauthorized",
            })
            c.Abort()
            return
        }
        
        // Получаем workspace_id из URL
        workspaceIDStr := c.Param("workspace_id")
        workspaceID, _ := uuid.Parse(workspaceIDStr)
        
        // Проверяем права
        uid := userID.(uuid.UUID)
        if err := m.workspaceService.CheckPermission(ctx, workspaceID, uid, requiredRole); err != nil {
            c.JSON(http.StatusForbidden, map[string]interface{}{
                "error": "Access denied",
            })
            c.Abort()
            return
        }
        
        // Сохраняем workspace_id для handler
        c.Set("workspace_id", workspaceID)
        c.Next(ctx)
    }
}
```

### ctx.Next() vs ctx.Abort()

```go
ctx.Next(c)   // Вызывает следующий обработчик в цепочке
ctx.Abort()   // Прерывает цепочку, следующие handlers не вызываются
```

**Сравнение с PHP (Symfony):**
```php
// Symfony Event Subscriber
class AuthSubscriber implements EventSubscriberInterface {
    public function onKernelRequest(RequestEvent $event): void {
        $request = $event->getRequest();
        $token = $request->headers->get('Authorization');
        
        if (!$this->isValidToken($token)) {
            $event->setResponse(new JsonResponse(['error' => 'Unauthorized'], 401));
            // Это прерывает дальнейшую обработку
        }
    }
}
```

---

## 14. Работа с базой данных (pgx)

### Что такое pgx?

[pgx](https://github.com/jackc/pgx) — это PostgreSQL драйвер для Go. Быстрее стандартного database/sql и поддерживает PostgreSQL-специфичные фичи.

**Документация:** https://pkg.go.dev/github.com/jackc/pgx/v5

### Connection Pool

Из `internal/database/postgres.go`:

```go
func NewPostgresPool(cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
    // Парсим DSN
    poolConfig, err := pgxpool.ParseConfig(cfg.GetDSN())
    if err != nil {
        return nil, fmt.Errorf("failed to parse database config: %w", err)
    }
    
    // Настраиваем пул соединений
    poolConfig.MaxConns = 100              // Максимум соединений
    poolConfig.MinConns = 10               // Минимум (всегда открыты)
    poolConfig.MaxConnLifetime = 1 * time.Hour
    poolConfig.MaxConnIdleTime = 30 * time.Minute
    poolConfig.HealthCheckPeriod = 1 * time.Minute
    
    // Создаём пул
    pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
    if err != nil {
        return nil, err
    }
    
    // Проверяем подключение
    if err := pool.Ping(ctx); err != nil {
        return nil, err
    }
    
    return pool, nil
}
```

### Выполнение запросов

Из `internal/repository/user_repository.go`:

**1. QueryRow — получение одной записи:**
```go
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
    query := `
        SELECT id, email, password_hash, name, avatar_url, provider, 
               provider_id, email_verified, created_at, updated_at
        FROM users
        WHERE id = $1
    `
    
    var user models.User
    err := r.db.QueryRow(ctx, query, id).Scan(
        &user.ID,
        &user.Email,
        &user.PasswordHash,  // *string — может быть nil
        &user.Name,
        &user.AvatarURL,
        &user.Provider,
        &user.ProviderID,
        &user.EmailVerified,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    // Обработка "не найдено"
    if err == pgx.ErrNoRows {
        return nil, nil  // Не ошибка, просто нет записи
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get user by id: %w", err)
    }
    
    return &user, nil
}
```

**2. Query — получение нескольких записей:**
```go
func (r *WorkspaceRepository) ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]models.WorkspaceMemberWithUser, error) {
    query := `
        SELECT m.id, m.workspace_id, m.user_id, m.role, m.joined_at,
               u.id, u.email, u.name, u.avatar_url
        FROM workspace_members m
        JOIN users u ON m.user_id = u.id
        WHERE m.workspace_id = $1
    `
    
    rows, err := r.db.Query(ctx, query, workspaceID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()  // Важно закрыть!
    
    var members []models.WorkspaceMemberWithUser
    for rows.Next() {
        var m models.WorkspaceMemberWithUser
        err := rows.Scan(
            &m.ID, &m.WorkspaceID, &m.UserID, &m.Role, &m.JoinedAt,
            &m.User.ID, &m.User.Email, &m.User.Name, &m.User.AvatarURL,
        )
        if err != nil {
            return nil, err
        }
        members = append(members, m)
    }
    
    return members, rows.Err()  // Проверяем ошибку итерации
}
```

**3. Exec — INSERT/UPDATE/DELETE:**
```go
func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
    query := `
        UPDATE users
        SET password_hash = $1, updated_at = NOW()
        WHERE id = $2
    `
    
    _, err := r.db.Exec(ctx, query, passwordHash, userID)
    if err != nil {
        return fmt.Errorf("failed to update password: %w", err)
    }
    
    return nil
}
```

**4. INSERT с RETURNING:**
```go
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
    query := `
        INSERT INTO users (email, password_hash, name, provider, provider_id, email_verified)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `
    
    // QueryRow + Scan для получения сгенерированных значений
    err := r.db.QueryRow(ctx, query,
        user.Email,
        user.PasswordHash,
        user.Name,
        user.Provider,
        user.ProviderID,
        user.EmailVerified,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
    
    return err
}
```

### Транзакции

Из `internal/database/migrate.go`:

```go
func Migrate(pool *pgxpool.Pool, migrationsPath string) error {
    for _, migration := range migrations {
        // Начинаем транзакцию
        tx, err := pool.Begin(ctx)
        if err != nil {
            return err
        }
        
        // Выполняем миграцию
        if _, err := tx.Exec(ctx, migration.SQL); err != nil {
            _ = tx.Rollback(ctx)  // Откат при ошибке
            return err
        }
        
        // Записываем в таблицу миграций
        if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations ..."); err != nil {
            _ = tx.Rollback(ctx)
            return err
        }
        
        // Фиксируем транзакцию
        if err := tx.Commit(ctx); err != nil {
            return err
        }
    }
}
```

**Сравнение с PHP (Doctrine):**
```php
// PHP Doctrine
$user = $entityManager->find(User::class, $id);

$entityManager->beginTransaction();
try {
    $entityManager->persist($user);
    $entityManager->flush();
    $entityManager->commit();
} catch (\Exception $e) {
    $entityManager->rollback();
    throw $e;
}
```

---

## 15. JSON и сериализация

### Struct Tags для JSON

```go
type User struct {
    ID           uuid.UUID `json:"id"`                        // Обычное поле
    Email        string    `json:"email"`                     // Обычное поле
    PasswordHash *string   `json:"-"`                         // Исключено из JSON
    AvatarURL    *string   `json:"avatar_url,omitempty"`     // Пропускается если nil
    CreatedAt    time.Time `json:"created_at"`                // time.Time → ISO 8601
}
```

### Результат сериализации

```go
user := &User{
    ID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
    Email:     "test@example.com",
    AvatarURL: nil,  // omitempty — не будет в JSON
}

jsonBytes, _ := json.Marshal(user)
// {"id":"123e4567-e89b-12d3-a456-426614174000","email":"test@example.com","created_at":"0001-01-01T00:00:00Z"}
```

### Кастомная сериализация JSONB

Из `internal/models/canvas.go`:

```go
type ElementData map[string]interface{}

// Scan — чтение из PostgreSQL JSONB
func (e *ElementData) Scan(value interface{}) error {
    if value == nil {
        *e = make(ElementData)
        return nil
    }
    
    bytes, ok := value.([]byte)
    if !ok {
        str, ok := value.(string)
        if !ok {
            return fmt.Errorf("failed to scan ElementData: unexpected type %T", value)
        }
        return json.Unmarshal([]byte(str), e)
    }
    return json.Unmarshal(bytes, e)
}

// Value — запись в PostgreSQL JSONB
func (e ElementData) Value() (driver.Value, error) {
    if e == nil {
        return "{}", nil
    }
    return json.Marshal(e)
}
```

Это позволяет использовать `map[string]interface{}` напрямую с PostgreSQL JSONB:

```go
element := &CanvasElement{
    ElementData: ElementData{
        "position": map[string]interface{}{"x": 100, "y": 200},
        "content":  "Hello World",
        "style":    map[string]interface{}{"fill": "#ffffff"},
    },
}
// INSERT ... element_data = '{"position":{"x":100,"y":200},...}'
```

---

## 16. JWT-аутентификация

### Структура JWT токена

JWT состоит из трёх частей: `header.payload.signature`

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzLi4uIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxNzA5MjM0NTY3fQ.signature
```

### Реализация в твоём проекте

Из `internal/service/jwt_service.go`:

```go
// Claims — данные, хранящиеся в токене
type Claims struct {
    UserID   uuid.UUID `json:"user_id"`
    Email    string    `json:"email"`
    Username string    `json:"username"`
    jwt.RegisteredClaims  // Встроенные поля (exp, iat, iss, ...)
}

type JWTService struct {
    secret               string
    accessTokenDuration  time.Duration   // 15 минут
    refreshTokenDuration time.Duration   // 7 дней
}

// Генерация Access Token
func (s *JWTService) GenerateAccessToken(userID uuid.UUID, email string, username ...string) (string, time.Time, error) {
    expiresAt := time.Now().Add(s.accessTokenDuration)
    
    claims := &Claims{
        UserID:   userID,
        Email:    email,
        Username: username[0],
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "hertzboard",
        },
    }
    
    // Создаём токен с алгоритмом HS256
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    // Подписываем секретным ключом
    tokenString, err := token.SignedString([]byte(s.secret))
    
    return tokenString, expiresAt, err
}

// Валидация Access Token
func (s *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        // Проверяем алгоритм подписи
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(s.secret), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, fmt.Errorf("invalid token")
}
```

### Refresh Token — хранение в БД

```go
// Генерация Refresh Token (не JWT!)
func (s *JWTService) GenerateRefreshToken() (token, tokenHash string, expiresAt time.Time, err error) {
    token = uuid.New().String()           // Просто UUID
    tokenHash = hashToken(token)          // SHA-256 хеш для хранения в БД
    expiresAt = time.Now().Add(s.refreshTokenDuration)
    return
}

// Хеширование токена
func hashToken(token string) string {
    hash := sha256.Sum256([]byte(token))
    return hex.EncodeToString(hash[:])
}
```

### Почему Refresh Token хранится в БД?

1. **Безопасность:** Можно отозвать токен (logout, смена пароля)
2. **Single-use:** После использования старый токен удаляется
3. **Хеширование:** В БД хранится хеш, не сам токен

**Сравнение с PHP:**
```php
// PHP (firebase/php-jwt)
use Firebase\JWT\JWT;

$payload = [
    'user_id' => $user->getId(),
    'email' => $user->getEmail(),
    'exp' => time() + 900,  // 15 минут
];

$token = JWT::encode($payload, $secret, 'HS256');
$decoded = JWT::decode($token, new Key($secret, 'HS256'));
```

---

## 17. WebSocket и реальное время

### Архитектура WebSocket в твоём проекте

```
Клиент ←→ WebSocketHandler ←→ Hub ←→ Room ←→ Клиенты
                                ↓
                             Redis Pub/Sub (для масштабирования)
```

### Модели WebSocket

Из `internal/models/websocket.go`:

```go
// Типы сообщений
type MessageType string

const (
    MessageTypeJoinRoom       MessageType = "join_room"
    MessageTypeUserJoined     MessageType = "user_joined"
    MessageTypeCursorMove     MessageType = "cursor_move"
    MessageTypeOperation      MessageType = "operation"
    MessageTypeSyncRequest    MessageType = "sync_request"
    // ...
)

// Сообщение WebSocket
type WSMessage struct {
    Type      MessageType `json:"type"`
    UserID    uuid.UUID   `json:"user_id,omitempty"`
    Payload   interface{} `json:"payload,omitempty"`
    Timestamp time.Time   `json:"timestamp"`
    RequestID string      `json:"request_id,omitempty"`
}

// Клиент
type Client struct {
    ID          uuid.UUID
    UserID      uuid.UUID
    WorkspaceID uuid.UUID
    Send        chan *WSMessage  // Канал для отправки сообщений
    Presence    *UserPresence
}

// Комната (workspace)
type Room struct {
    WorkspaceID uuid.UUID
    Clients     map[uuid.UUID]*Client
    Broadcast   chan *WSMessage
    Register    chan *Client
    Unregister  chan *Client
}
```

### Hub — центральный координатор

Из `internal/service/hub.go`:

```go
type Hub struct {
    rooms map[uuid.UUID]*models.Room  // workspaceID → Room
    redis *redis.Client               // Для масштабирования
    mu    sync.RWMutex                // Защита карты rooms
}

// Регистрация клиента
func (h *Hub) Register(client *models.Client) {
    h.mu.Lock()
    defer h.mu.Unlock()
    
    room, exists := h.rooms[client.WorkspaceID]
    if !exists {
        // Создаём новую комнату
        room = &models.Room{
            WorkspaceID: client.WorkspaceID,
            Clients:     make(map[uuid.UUID]*models.Client),
            Broadcast:   make(chan *models.WSMessage, 256),
            Register:    make(chan *models.Client),
            Unregister:  make(chan *models.Client),
        }
        h.rooms[client.WorkspaceID] = room
        
        // Запускаем горутину для обработки сообщений комнаты
        go h.runRoom(room)
    }
    
    // Проверяем лимит клиентов
    if len(room.Clients) >= maxClientsPerRoom {
        h.sendErrorToClient(client, "room_full", "Room has reached maximum capacity")
        return
    }
    
    room.Register <- client  // Отправляем клиента в канал регистрации
}
```

### Обработка комнаты

```go
func (h *Hub) runRoom(room *models.Room) {
    for {
        select {
        case client := <-room.Register:
            // Добавляем клиента
            room.Clients[client.ID] = client
            
            // Отправляем существующих пользователей новому клиенту
            h.sendExistingPresences(client, room)
            
            // Уведомляем остальных о новом пользователе
            joinMsg := &models.WSMessage{
                Type:   models.MessageTypeUserJoined,
                UserID: client.UserID,
                Payload: models.UserJoinedPayload{
                    UserID:    client.UserID,
                    UserName:  client.UserName,
                    UserColor: client.UserColor,
                },
            }
            h.broadcastToRoomClients(room, joinMsg, client.ID)
            
        case client := <-room.Unregister:
            // Удаляем клиента
            if _, ok := room.Clients[client.ID]; ok {
                delete(room.Clients, client.ID)
                close(client.Send)
                
                // Уведомляем остальных
                leaveMsg := &models.WSMessage{
                    Type:   models.MessageTypeUserLeft,
                    UserID: client.UserID,
                }
                h.broadcastToRoomClients(room, leaveMsg, uuid.Nil)
            }
            
        case message := <-room.Broadcast:
            // Рассылаем всем клиентам
            h.broadcastToRoomClients(room, message, uuid.Nil)
        }
    }
}
```

### Масштабирование через Redis Pub/Sub

```go
// Публикация сообщения в Redis (для других инстансов)
func (h *Hub) publishToRedis(workspaceID uuid.UUID, msg *models.WSMessage, excludeClientID uuid.UUID) {
    redisMsg := RedisMessage{
        WorkspaceID:     workspaceID,
        Message:         msg,
        ExcludeClientID: excludeClientID,
    }
    
    data, _ := json.Marshal(redisMsg)
    channel := fmt.Sprintf("workspace:%s", workspaceID)
    h.redis.Publish(h.ctx, channel, data)
}

// Подписка на Redis (в отдельной горутине)
func (h *Hub) subscribeToRedis() {
    pubsub := h.redis.PSubscribe(h.ctx, "workspace:*")
    ch := pubsub.Channel()
    
    for msg := range ch {
        var redisMsg RedisMessage
        json.Unmarshal([]byte(msg.Payload), &redisMsg)
        
        // Отправляем локальным клиентам
        h.mu.RLock()
        room, exists := h.rooms[redisMsg.WorkspaceID]
        h.mu.RUnlock()
        
        if exists {
            h.broadcastToRoomClients(room, redisMsg.Message, redisMsg.ExcludeClientID)
        }
    }
}
```

---

## 18. CRDT — синхронизация данных

### Что такое CRDT?

**CRDT** (Conflict-free Replicated Data Type) — структуры данных, которые автоматически разрешают конфликты при одновременном редактировании.

В твоём проекте используется **Last-Write-Wins (LWW)** стратегия с **Lamport Timestamps**.

### Lamport Clock

Из `internal/service/crdt_service.go`:

```go
type LamportClock struct {
    counter int64
    mu      sync.Mutex
}

// Tick — локальное событие (создание/изменение)
func (lc *LamportClock) Tick() int64 {
    lc.mu.Lock()
    defer lc.mu.Unlock()
    lc.counter++
    return lc.counter
}

// Update — получили событие от другого клиента
func (lc *LamportClock) Update(receivedTimestamp int64) int64 {
    lc.mu.Lock()
    defer lc.mu.Unlock()
    // max(local, received) + 1
    if receivedTimestamp > lc.counter {
        lc.counter = receivedTimestamp
    }
    lc.counter++
    return lc.counter
}
```

### Применение операции

```go
func (s *CRDTService) ApplyOperation(op *models.OperationPayload) error {
    // Обновляем Lamport clock
    s.clock.Update(op.Timestamp)
    
    // Сохраняем операцию в БД (для истории/синхронизации)
    s.operationRepo.Create(s.ctx, &models.Operation{
        ID:          uuid.New(),
        WorkspaceID: op.WorkspaceID,
        ElementID:   op.ElementID,
        UserID:      op.UserID,
        OpType:      string(op.OpType),
        Data:        op.Data,
        Timestamp:   op.Timestamp,
    })
    
    // Применяем по типу операции
    switch op.OpType {
    case models.OperationTypeCreate:
        return s.applyCreate(op)
    case models.OperationTypeUpdate:
        return s.applyUpdate(op)
    case models.OperationTypeDelete:
        return s.applyDelete(op)
    case models.OperationTypeMove:
        return s.applyMove(op)
    }
}
```

### Last-Write-Wins для Update

```go
func (s *CRDTService) applyUpdate(op *models.OperationPayload) error {
    existing, err := s.elementRepo.GetByID(s.ctx, op.ElementID)
    if err != nil {
        return err
    }
    
    // LWW: применяем только если timestamp новее
    if op.Timestamp <= existing.Version {
        return nil  // Игнорируем старую операцию
    }
    
    // Применяем изменения...
    existing.Version = op.Timestamp
    return s.elementRepo.Update(s.ctx, existing)
}
```

### Разрешение конфликтов

```go
// Если timestamps равны — детерминированный tiebreaker по UserID
func (s *CRDTService) ResolveConflict(op1, op2 *models.OperationPayload) *models.OperationPayload {
    if op1.Timestamp != op2.Timestamp {
        if op1.Timestamp > op2.Timestamp {
            return op1
        }
        return op2
    }
    
    // Tiebreaker: больший UserID "выигрывает"
    if op1.UserID.String() > op2.UserID.String() {
        return op1
    }
    return op2
}
```

### State Vector для синхронизации

```go
// Клиент отправляет свой state vector:
// {"user1_id": 15, "user2_id": 23}
// Означает: "Я видел операции user1 до timestamp=15, user2 до timestamp=23"

func (s *CRDTService) GetOperationsSince(
    workspaceID uuid.UUID,
    stateVector map[string]int64,
) ([]*models.Operation, error) {
    operations, _ := s.operationRepo.GetByWorkspaceID(s.ctx, workspaceID, 1000)
    
    result := make([]*models.Operation, 0)
    for _, op := range operations {
        userIDStr := op.UserID.String()
        lastSeen := stateVector[userIDStr]
        
        // Включаем операцию если клиент её не видел
        if op.Timestamp > lastSeen {
            result = append(result, op)
        }
    }
    
    return result, nil
}
```

---

## 19. Архитектура проекта

### Слои архитектуры

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Layer (Hertz)                     │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────────────┐   │
│  │ Router  │→│Middleware│→│ Handler │→│ JSON Response   │   │
│  └─────────┘ └─────────┘ └─────────┘ └─────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     Service Layer                           │
│  ┌──────────────┐ ┌─────────────────┐ ┌─────────────────┐  │
│  │ AuthService  │ │ WorkspaceService│ │  CanvasService  │  │
│  └──────────────┘ └─────────────────┘ └─────────────────┘  │
│  ┌──────────────┐ ┌─────────────────┐ ┌─────────────────┐  │
│  │ JWTService   │ │   EmailService  │ │   CRDTService   │  │
│  └──────────────┘ └─────────────────┘ └─────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Repository Layer                          │
│  ┌──────────────┐ ┌─────────────────┐ ┌─────────────────┐  │
│  │ UserRepo     │ │ WorkspaceRepo   │ │  CanvasRepo     │  │
│  └──────────────┘ └─────────────────┘ └─────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Infrastructure                            │
│  ┌──────────────┐ ┌─────────────────┐ ┌─────────────────┐  │
│  │  PostgreSQL  │ │      Redis      │ │      NATS       │  │
│  └──────────────┘ └─────────────────┘ └─────────────────┘  │
│  ┌──────────────┐                                          │
│  │    MinIO     │                                          │
│  └──────────────┘                                          │
└─────────────────────────────────────────────────────────────┘
```

### Dependency Injection (ручной)

Из `cmd/api-gateway/main.go`:

```go
func main() {
    // 1. Загрузка конфигурации
    cfg, _ := config.Load(configPath)
    
    // 2. Подключение к инфраструктуре
    dbPool, _ := database.NewPostgresPool(&cfg.Database)
    redisClient, _ := database.NewRedisClient(&cfg.Redis)
    natsConn, _ := database.NewNATSConnection(&cfg.NATS)
    
    // 3. Создание репозиториев (зависят от БД)
    userRepo := repository.NewUserRepository(dbPool)
    workspaceRepo := repository.NewWorkspaceRepository(dbPool)
    
    // 4. Создание сервисов (зависят от репозиториев)
    jwtService, _ := service.NewJWTService(&cfg.JWT)
    emailService := service.NewEmailService(&cfg.Email, natsConn)
    authService := service.NewAuthService(userRepo, jwtService)
    workspaceService := service.NewWorkspaceService(workspaceRepo, userRepo, emailService)
    
    // 5. Создание обработчиков (зависят от сервисов)
    authHandler := handler.NewAuthHandler(authService)
    workspaceHandler := handler.NewWorkspaceHandler(workspaceService)
    
    // 6. Сборка зависимостей для роутера
    deps := &router.Dependencies{
        JWTService:       jwtService,
        WorkspaceService: workspaceService,
        AuthHandler:      authHandler,
        WorkspaceHandler: workspaceHandler,
        // ...
    }
    
    // 7. Настройка маршрутов
    router.Setup(h, cfg, deps)
}
```

**Сравнение с PHP (Symfony DI Container):**
```php
// Symfony делает это автоматически через autowiring
// services.yaml
services:
    _defaults:
        autowire: true
        autoconfigure: true
    
    App\Service\AuthService:
        arguments:
            $userRepository: '@App\Repository\UserRepository'
            $jwtService: '@App\Service\JWTService'
```

### Graceful Shutdown

```go
// Ожидание сигнала завершения
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

log.Println("Shutting down server...")

// Контекст с таймаутом для graceful shutdown
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Останавливаем сервер (ждём завершения текущих запросов)
if err := h.Shutdown(ctx); err != nil {
    log.Fatalf("Server forced to shutdown: %v", err)
}

// defer'ы закроют соединения с БД, Redis, NATS
```

---

## 20. Тест на проверку знаний

### Часть 1: Основы Go (10 вопросов)

**1. Что выведет этот код?**
```go
var s []int
fmt.Println(s == nil, len(s))
```
- A) `true 0`
- B) `false 0`
- C) `panic`
- D) `true -1`

**2. Какой будет результат?**
```go
func modify(m map[string]int) {
    m["key"] = 100
}

func main() {
    data := map[string]int{"key": 1}
    modify(data)
    fmt.Println(data["key"])
}
```
- A) `1`
- B) `100`
- C) `0`
- D) `panic`

**3. Что не так с этим кодом?**
```go
func getUser() *User {
    user := User{Name: "John"}
    return &user
}
```
- A) Возврат указателя на локальную переменную — undefined behavior
- B) Код правильный, Go автоматически перемещает на heap
- C) Нужно использовать new(User)
- D) Ошибка компиляции

**4. Какой тип имеет переменная `x`?**
```go
x := make(chan *User, 10)
```
- A) `chan User`
- B) `chan *User`
- C) `*chan User`
- D) `[]chan *User`

**5. Что означает `defer`?**
```go
func process() {
    defer fmt.Println("A")
    fmt.Println("B")
    defer fmt.Println("C")
}
```
- A) Выведет: A, B, C
- B) Выведет: B, A, C
- C) Выведет: B, C, A
- D) Выведет: C, B, A

**6. Какой результат?**
```go
type Counter struct {
    value int
}

func (c Counter) Increment() {
    c.value++
}

func main() {
    c := Counter{value: 0}
    c.Increment()
    c.Increment()
    fmt.Println(c.value)
}
```
- A) `0`
- B) `1`
- C) `2`
- D) `panic`

**7. Что делает `ctx.Abort()` в middleware Hertz?**
- A) Вызывает panic
- B) Прерывает цепочку middleware, следующие handlers не вызываются
- C) Возвращает 500 ошибку
- D) Закрывает соединение

**8. Почему в Go нет конструкторов?**
- A) Это ошибка дизайна языка
- B) Go использует функции-фабрики (New...) вместо конструкторов
- C) Можно использовать init() как конструктор
- D) Структуры не нуждаются в инициализации

**9. Что такое `interface{}`?**
- A) Пустой интерфейс, который не реализует ни один тип
- B) Интерфейс, который реализуют все типы
- C) Ошибка синтаксиса
- D) Указатель на любой тип

**10. Как правильно проверить ошибку в Go?**
```go
result, err := someFunction()
```
- A) `if err != nil { return err }`
- B) `try { } catch (err) { }`
- C) `if err.IsError() { }`
- D) `throw err`

### Часть 2: Твой проект (10 вопросов)

**11. В каком файле определена структура `User`?**
- A) `internal/handler/user_handler.go`
- B) `internal/models/user.go`
- C) `internal/service/auth_service.go`
- D) `internal/repository/user_repository.go`

**12. Какой HTTP-фреймворк используется в проекте?**
- A) Gin
- B) Echo
- C) Hertz
- D) Fiber

**13. Где хранится refresh token?**
- A) В JWT payload
- B) В Redis
- C) В PostgreSQL (хеш)
- D) В cookie

**14. Что делает `WorkspaceMiddleware.RequireWorkspaceAccess(models.WorkspaceRoleEditor)`?**
- A) Требует роль Owner
- B) Требует роль Editor или выше (Editor, Owner)
- C) Требует роль Viewer или выше
- D) Разрешает доступ всем

**15. Какой алгоритм используется для разрешения конфликтов в CRDT?**
- A) First-Write-Wins
- B) Last-Write-Wins с Lamport Timestamps
- C) Version Vectors
- D) Operational Transformation

**16. Для чего используется NATS в проекте?**
- A) Хранение сессий
- B) Асинхронная отправка email
- C) WebSocket broadcast
- D) Кеширование

**17. Что означает тег `json:"-"` в структуре?**
```go
PasswordHash *string `json:"-"`
```
- A) Поле обязательное
- B) Поле исключено из JSON сериализации
- C) Поле nullable
- D) Поле private

**18. Какой паттерн используется для работы с БД?**
- A) Active Record
- B) Repository Pattern
- C) Table Gateway
- D) Data Mapper

**19. Где происходит валидация JWT токена?**
- A) В handler
- B) В service
- C) В middleware
- D) В repository

**20. Зачем в Hub используется Redis Pub/Sub?**
- A) Для хранения WebSocket сообщений
- B) Для масштабирования на несколько серверов
- C) Для аутентификации WebSocket
- D) Для логирования

---

### Ответы

<details>
<summary>Нажми, чтобы увидеть ответы</summary>

1. **A** — nil slice имеет длину 0, но сам равен nil
2. **B** — maps передаются по ссылке
3. **B** — Go делает escape analysis и выделяет на heap если нужно
4. **B** — буферизованный канал указателей на User
5. **C** — defer выполняются в обратном порядке (LIFO)
6. **A** — receiver по значению, оригинал не меняется
7. **B** — Abort() прерывает цепочку handlers
8. **B** — Go использует New... функции-фабрики
9. **B** — пустой интерфейс реализуется любым типом
10. **A** — стандартный паттерн проверки ошибок

11. **B** — `internal/models/user.go`
12. **C** — Hertz от CloudWeGo
13. **C** — PostgreSQL (хранится хеш, не сам токен)
14. **B** — Editor или Owner могут выполнять действие
15. **B** — Last-Write-Wins с Lamport Timestamps
16. **B** — Асинхронная очередь для отправки email
17. **B** — Поле не сериализуется в JSON (безопасность)
18. **B** — Repository Pattern
19. **C** — В `middleware/auth.go`
20. **B** — Для масштабирования (сообщения между инстансами)

**Оценка:**
- 18-20: Отлично! Ты хорошо понимаешь Go и свой проект
- 14-17: Хорошо! Есть что повторить
- 10-13: Удовлетворительно. Перечитай соответствующие разделы
- <10: Рекомендуется пройти материал заново

</details>

---

## Полезные ссылки

- [Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Hertz Documentation](https://www.cloudwego.io/docs/hertz/)
- [pgx Documentation](https://pkg.go.dev/github.com/jackc/pgx/v5)
- [JWT Go](https://pkg.go.dev/github.com/golang-jwt/jwt/v5)
- [Go by Example](https://gobyexample.com/)

---

*Документ создан на основе анализа кода HertzBoard Backend*
