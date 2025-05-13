package config

import (
	"action-detector-backend/pkg/logger"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

const (
	EnvironmentProd = "production"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	HTTPHost string
	HTTPPort int

	Environment string
	Debug       bool

	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string

	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int

	JWTSecret                string
	JWTAccessExpirationHours int
	JWTRefreshExpirationDays int

	HashKey string

	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioUseSSL     bool
	MinioBucketName string
	MinioFileUrl    string

	PredictApiURL string
}

func GetConfig() *Config {
	once.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			logger.GetLogger().Error(".env file not found")
		}

		instance = &Config{
			HTTPHost:    cast.ToString(getOrReturnDefault("HTTP_HOST", "84.247.173.245")),
			HTTPPort:    cast.ToInt(getOrReturnDefault("HTTP_PORT", 4050)),
			Environment: cast.ToString(getOrReturnDefault("ENVIRONMENT", "production")),
			Debug:       cast.ToBool(getOrReturnDefault("DEBUG", false)),

			PostgresHost:     cast.ToString(getOrReturnDefault("POSTGRE_HOST", "84.247.173.245")),
			PostgresPort:     cast.ToInt(getOrReturnDefault("POSTGRE_PORT", 5432)),
			PostgresDatabase: cast.ToString(getOrReturnDefault("POSTGRE_DB", "usta_db")),
			PostgresUser:     cast.ToString(getOrReturnDefault("POSTGRE_USER", "postgres")),
			PostgresPassword: cast.ToString(getOrReturnDefault("POSTGRE_PASSWORD", "Hasanov@1209DD")),

			RedisHost:     cast.ToString(getOrReturnDefault("REDIS_HOST", "localhost")),
			RedisPort:     cast.ToInt(getOrReturnDefault("REDIS_PORT", 6379)),
			RedisPassword: cast.ToString(getOrReturnDefault("REDIS_PASSWORD", "")),
			RedisDB:       cast.ToInt(getOrReturnDefault("REDIS_DB", 0)),

			JWTSecret:                cast.ToString(getOrReturnDefault("JWT_SECRET", "usta365-forever-2025")),
			JWTAccessExpirationHours: cast.ToInt(getOrReturnDefault("JWT_ACCESS_EXPIRATION_HOURS", 12)),
			JWTRefreshExpirationDays: cast.ToInt(getOrReturnDefault("JWT_REFRESH_EXPIRATION_DAYS", 3)),

			HashKey: cast.ToString(getOrReturnDefault("HASH_KEY", "skd32r8w3hoqN2HSdvw")),

			MinioEndpoint:   cast.ToString(getOrReturnDefault("MINIO_ENDPOINT", "api.cdn.demo.tn.uz")),
			MinioAccessKey:  cast.ToString(getOrReturnDefault("MINIO_ACCESS_KEY", "Mrb0224")),
			MinioSecretKey:  cast.ToString(getOrReturnDefault("MINIO_SECRET_KEY", "Hero571Hero")),
			MinioUseSSL:     cast.ToBool(getOrReturnDefault("MINIO_USE_SSL", true)),
			MinioBucketName: cast.ToString(getOrReturnDefault("MINIO_BUCKET_NAME", "turonlife")),
			MinioFileUrl:    cast.ToString(getOrReturnDefault("MINIO_FILE_URL", "")),

			PredictApiURL: cast.ToString(getOrReturnDefault("PREDICT_API_URL", "")),
		}
	})

	return instance
}

func getOrReturnDefault(key string, defaultValue any) any {
	val, exists := os.LookupEnv(key)
	if exists {
		return val
	}
	return defaultValue
}
