package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

//nolint:govet // fieldalignment: struct field order optimized for readability over memory
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

//nolint:govet // fieldalignment: struct field order optimized for readability
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

//nolint:govet // fieldalignment: struct field order optimized for readability
type RedisConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Password   string `yaml:"password"`
	DB         int    `yaml:"db"`
	MaxRetries int    `yaml:"max_retries"`
	PoolSize   int    `yaml:"pool_size"`
}

//nolint:govet // fieldalignment: struct field order optimized for readability
type MinIOConfig struct {
	Endpoint      string `yaml:"endpoint"`
	AccessKey     string `yaml:"access_key"`
	SecretKey     string `yaml:"secret_key"`
	UseSSL        bool   `yaml:"use_ssl"`
	BucketAssets  string `yaml:"bucket_assets"`
	BucketExports string `yaml:"bucket_exports"`
	BucketBackups string `yaml:"bucket_backups"`
}

//nolint:govet // fieldalignment: struct field order optimized for readability
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

//nolint:govet // fieldalignment: struct field order optimized for readability
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

//nolint:govet // fieldalignment: struct field order optimized for readability
type UploadConfig struct {
	MaxSize      int64    `yaml:"max_size"`
	AllowedTypes []string `yaml:"allowed_types"`
}

//nolint:govet // fieldalignment: struct field order optimized for readability
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

//nolint:govet // fieldalignment: struct field order optimized for readability
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
