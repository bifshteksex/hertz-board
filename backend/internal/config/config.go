package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Default configuration values
const (
	DefaultAppPort                 = 8080
	DefaultDBPort                  = 5432
	DefaultDBMaxConnections        = 100
	DefaultDBMaxIdleConnections    = 10
	DefaultDBConnectionMaxLifetime = 3600
	DefaultRedisPort               = 6379
	DefaultRedisMaxRetries         = 3
	DefaultRedisPoolSize           = 10
	DefaultClickHousePort          = 8123
	DefaultNATSMaxReconnect        = 10
	DefaultNATSReconnectWait       = 2
	DefaultSMTPPort                = 1025
	DefaultWSPort                  = 8080
	DefaultWSReadBufferSize        = 1024
	DefaultWSWriteBufferSize       = 1024
	DefaultWSMaxMessageSize        = 10485760 // 10MB
	DefaultWSPingPeriod            = 54
	DefaultWSPongWait              = 60
	DefaultWSWriteWait             = 10
	DefaultMaxUploadSize           = 10485760 // 10MB
	DefaultRateLimitRequests       = 100
	DefaultMetricsPort             = 9090
)

type Config struct {
	App        AppConfig        `yaml:"app"`
	Database   DatabaseConfig   `yaml:"database"`
	Redis      RedisConfig      `yaml:"redis"`
	MinIO      MinIOConfig      `yaml:"minio"`
	ClickHouse ClickHouseConfig `yaml:"clickhouse"`
	NATS       NATSConfig       `yaml:"nats"`
	JWT        JWTConfig        `yaml:"jwt"`
	OAuth      OAuthConfig      `yaml:"oauth"`
	Email      EmailConfig      `yaml:"email"`
	CORS       CORSConfig       `yaml:"cors"`
	WebSocket  WebSocketConfig  `yaml:"websocket"`
	Upload     UploadConfig     `yaml:"upload"`
	RateLimit  RateLimitConfig  `yaml:"rate_limit"`
	Logging    LoggingConfig    `yaml:"logging"`
	Metrics    MetricsConfig    `yaml:"metrics"`
	Tracing    TracingConfig    `yaml:"tracing"`
}

type AppConfig struct {
	Name  string `yaml:"name"`
	Env   string `yaml:"env"`
	Port  int    `yaml:"port"`
	Debug bool   `yaml:"debug"`
}

type DatabaseConfig struct {
	Host                  string `yaml:"host"`
	Port                  int    `yaml:"port"`
	Name                  string `yaml:"name"`
	User                  string `yaml:"user"`
	Password              string `yaml:"password"`
	SSLMode               string `yaml:"ssl_mode"`
	MaxConnections        int    `yaml:"max_connections"`
	MaxIdleConnections    int    `yaml:"max_idle_connections"`
	ConnectionMaxLifetime int    `yaml:"connection_max_lifetime"`
}

type RedisConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Password   string `yaml:"password"`
	DB         int    `yaml:"db"`
	MaxRetries int    `yaml:"max_retries"`
	PoolSize   int    `yaml:"pool_size"`
}

type MinIOConfig struct {
	Endpoint       string `yaml:"endpoint"`
	PublicEndpoint string `yaml:"public_endpoint"`
	AccessKey      string `yaml:"access_key"`
	SecretKey      string `yaml:"secret_key"`
	UseSSL         bool   `yaml:"use_ssl"`
	PublicUseSSL   bool   `yaml:"public_use_ssl"`
	BucketAssets   string `yaml:"bucket_assets"`
	BucketExports  string `yaml:"bucket_exports"`
	BucketBackups  string `yaml:"bucket_backups"`
}

type ClickHouseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type NATSConfig struct {
	URL           string `yaml:"url"`
	MaxReconnect  int    `yaml:"max_reconnect"`
	ReconnectWait int    `yaml:"reconnect_wait"`
}

type JWTConfig struct {
	Secret             string `yaml:"secret"`
	AccessTokenExpiry  string `yaml:"access_token_expiry"`
	RefreshTokenExpiry string `yaml:"refresh_token_expiry"`
}

type OAuthProviderConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURL  string `yaml:"redirect_url"`
}

type OAuthConfig struct {
	Google OAuthProviderConfig `yaml:"google"`
	GitHub OAuthProviderConfig `yaml:"github"`
}

type EmailConfig struct {
	SMTPHost     string `yaml:"smtp_host"`
	SMTPPort     int    `yaml:"smtp_port"`
	SMTPUser     string `yaml:"smtp_user"`
	SMTPPassword string `yaml:"smtp_password"`
	From         string `yaml:"from"`
}

type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

type WebSocketConfig struct {
	Port            int `yaml:"port"`
	ReadBufferSize  int `yaml:"read_buffer_size"`
	WriteBufferSize int `yaml:"write_buffer_size"`
	MaxMessageSize  int `yaml:"max_message_size"`
	PingPeriod      int `yaml:"ping_period"`
	PongWait        int `yaml:"pong_wait"`
	WriteWait       int `yaml:"write_wait"`
}

type UploadConfig struct {
	MaxSize      int64    `yaml:"max_size"`
	AllowedTypes []string `yaml:"allowed_types"`
}

type RateLimitConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Requests int    `yaml:"requests"`
	Duration string `yaml:"duration"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

type MetricsConfig struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type TracingConfig struct {
	Enabled        bool   `yaml:"enabled"`
	JaegerEndpoint string `yaml:"jaeger_endpoint"`
}

// Load reads configuration from a YAML file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Expand environment variables in the config
	expandedData := []byte(os.ExpandEnv(string(data)))

	var cfg Config
	if err := yaml.Unmarshal(expandedData, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// GetDSN returns PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// GetRedisAddr returns Redis address
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetAccessTokenDuration parses access token expiry duration
func (c *JWTConfig) GetAccessTokenDuration() (time.Duration, error) {
	return time.ParseDuration(c.AccessTokenExpiry)
}

// GetRefreshTokenDuration parses refresh token expiry duration
func (c *JWTConfig) GetRefreshTokenDuration() (time.Duration, error) {
	return time.ParseDuration(c.RefreshTokenExpiry)
}

// LoadFromEnv loads configuration from environment variables
// This is used in production when CONFIG_PATH is not set
func LoadFromEnv() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Name:  getEnvOrDefault("APP_NAME", "HertzBoard"),
			Env:   getEnvOrDefault("APP_ENV", "production"),
			Port:  getEnvAsIntOrDefault("APP_PORT", DefaultAppPort),
			Debug: getEnvAsBoolOrDefault("APP_DEBUG", false),
		},
		Database: DatabaseConfig{
			Host:                  getEnvOrDefault("DB_HOST", "postgres"),
			Port:                  getEnvAsIntOrDefault("DB_PORT", DefaultDBPort),
			Name:                  getEnvOrDefault("DB_NAME", "hertzboard"),
			User:                  getEnvOrDefault("DB_USER", "hertzboard"),
			Password:              os.Getenv("DB_PASSWORD"),
			SSLMode:               getEnvOrDefault("DB_SSL_MODE", "require"),
			MaxConnections:        getEnvAsIntOrDefault("DB_MAX_CONNECTIONS", DefaultDBMaxConnections),
			MaxIdleConnections:    getEnvAsIntOrDefault("DB_MAX_IDLE_CONNECTIONS", DefaultDBMaxIdleConnections),
			ConnectionMaxLifetime: getEnvAsIntOrDefault("DB_CONNECTION_MAX_LIFETIME", DefaultDBConnectionMaxLifetime),
		},
		Redis: RedisConfig{
			Host:       getEnvOrDefault("REDIS_HOST", "redis"),
			Port:       getEnvAsIntOrDefault("REDIS_PORT", DefaultRedisPort),
			Password:   os.Getenv("REDIS_PASSWORD"),
			DB:         getEnvAsIntOrDefault("REDIS_DB", 0),
			MaxRetries: getEnvAsIntOrDefault("REDIS_MAX_RETRIES", DefaultRedisMaxRetries),
			PoolSize:   getEnvAsIntOrDefault("REDIS_POOL_SIZE", DefaultRedisPoolSize),
		},
		MinIO: MinIOConfig{
			Endpoint:       getEnvOrDefault("MINIO_ENDPOINT", "minio:9000"),
			PublicEndpoint: getEnvOrDefault("MINIO_PUBLIC_ENDPOINT", ""), // Default to endpoint if not set
			AccessKey:      getEnvOrDefault("MINIO_ACCESS_KEY", "hertzboard"),
			SecretKey:      os.Getenv("MINIO_SECRET_KEY"),
			UseSSL:         getEnvAsBoolOrDefault("MINIO_USE_SSL", false),
			PublicUseSSL:   getEnvAsBoolOrDefault("MINIO_PUBLIC_USE_SSL", false),
			BucketAssets:   getEnvOrDefault("MINIO_BUCKET_ASSETS", "hertzboard-assets"),
			BucketExports:  getEnvOrDefault("MINIO_BUCKET_EXPORTS", "hertzboard-exports"),
			BucketBackups:  getEnvOrDefault("MINIO_BUCKET_BACKUPS", "hertzboard-backups"),
		},
		ClickHouse: ClickHouseConfig{
			Host:     getEnvOrDefault("CLICKHOUSE_HOST", "clickhouse"),
			Port:     getEnvAsIntOrDefault("CLICKHOUSE_PORT", DefaultClickHousePort),
			Database: getEnvOrDefault("CLICKHOUSE_DATABASE", "hertzboard_analytics"),
			User:     getEnvOrDefault("CLICKHOUSE_USER", "hertzboard"),
			Password: os.Getenv("CLICKHOUSE_PASSWORD"),
		},
		NATS: NATSConfig{
			URL:           getEnvOrDefault("NATS_URL", "nats://nats:4222"),
			MaxReconnect:  getEnvAsIntOrDefault("NATS_MAX_RECONNECT", DefaultNATSMaxReconnect),
			ReconnectWait: getEnvAsIntOrDefault("NATS_RECONNECT_WAIT", DefaultNATSReconnectWait),
		},
		JWT: JWTConfig{
			Secret:             os.Getenv("JWT_SECRET"),
			AccessTokenExpiry:  getEnvOrDefault("JWT_ACCESS_TOKEN_EXPIRY", "15m"),
			RefreshTokenExpiry: getEnvOrDefault("JWT_REFRESH_TOKEN_EXPIRY", "7d"),
		},
		OAuth: OAuthConfig{
			Google: OAuthProviderConfig{
				ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
				ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
				RedirectURL:  getEnvOrDefault("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
			},
			GitHub: OAuthProviderConfig{
				ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
				ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
				RedirectURL:  getEnvOrDefault("GITHUB_REDIRECT_URL", "http://localhost:8080/auth/github/callback"),
			},
		},
		Email: EmailConfig{
			SMTPHost:     getEnvOrDefault("SMTP_HOST", "localhost"),
			SMTPPort:     getEnvAsIntOrDefault("SMTP_PORT", DefaultSMTPPort),
			SMTPUser:     os.Getenv("SMTP_USER"),
			SMTPPassword: os.Getenv("SMTP_PASSWORD"),
			From:         getEnvOrDefault("SMTP_FROM", "noreply@hertzboard.dev"),
		},
		CORS: CORSConfig{
			AllowedOrigins:   []string{getEnvOrDefault("CORS_ALLOWED_ORIGINS", "http://localhost:5173")},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowedHeaders:   []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           86400,
		},
		WebSocket: WebSocketConfig{
			Port:            getEnvAsIntOrDefault("WS_PORT", DefaultWSPort),
			ReadBufferSize:  getEnvAsIntOrDefault("WS_READ_BUFFER_SIZE", DefaultWSReadBufferSize),
			WriteBufferSize: getEnvAsIntOrDefault("WS_WRITE_BUFFER_SIZE", DefaultWSWriteBufferSize),
			MaxMessageSize:  getEnvAsIntOrDefault("WS_MAX_MESSAGE_SIZE", DefaultWSMaxMessageSize),
			PingPeriod:      getEnvAsIntOrDefault("WS_PING_PERIOD", DefaultWSPingPeriod),
			PongWait:        getEnvAsIntOrDefault("WS_PONG_WAIT", DefaultWSPongWait),
			WriteWait:       getEnvAsIntOrDefault("WS_WRITE_WAIT", DefaultWSWriteWait),
		},
		Upload: UploadConfig{
			MaxSize:      int64(getEnvAsIntOrDefault("MAX_UPLOAD_SIZE", DefaultMaxUploadSize)),
			AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "image/webp", "image/svg+xml"},
		},
		RateLimit: RateLimitConfig{
			Enabled:  getEnvAsBoolOrDefault("RATE_LIMIT_ENABLED", true),
			Requests: getEnvAsIntOrDefault("RATE_LIMIT_REQUESTS", DefaultRateLimitRequests),
			Duration: getEnvOrDefault("RATE_LIMIT_DURATION", "1m"),
		},
		Logging: LoggingConfig{
			Level:  getEnvOrDefault("LOG_LEVEL", "info"),
			Format: getEnvOrDefault("LOG_FORMAT", "json"),
			Output: getEnvOrDefault("LOG_OUTPUT", "stdout"),
		},
		Metrics: MetricsConfig{
			Enabled: getEnvAsBoolOrDefault("METRICS_ENABLED", true),
			Port:    getEnvAsIntOrDefault("METRICS_PORT", DefaultMetricsPort),
		},
		Tracing: TracingConfig{
			Enabled:        getEnvAsBoolOrDefault("TRACING_ENABLED", false),
			JaegerEndpoint: os.Getenv("JAEGER_ENDPOINT"),
		},
	}

	return cfg, nil
}

// Helper functions for environment variable parsing
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}
