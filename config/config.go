package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"

	_ "github.com/lib/pq"
)

type Config struct {
	HTTP_PORT        string
	GRPC_PORT        string
	DB_HOST          string
	DB_PORT          string
	DB_USER          string
	DB_NAME          string
	DB_PASSWORD      string
	DB_CASBIN_DRIVER string
	REFRESH_TOKEN    string
	ACCESS_TOKEN     string
}

func Load() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	config := Config{}

	config.HTTP_PORT = cast.ToString(coalesce("HTTP_PORT", ":8081"))
	config.GRPC_PORT = cast.ToString(coalesce("GRPC_PORT", ":50051"))
	config.DB_HOST = cast.ToString(coalesce("DB_HOST", "localhost"))
	config.DB_PORT = cast.ToString(coalesce("DB_PORT", "5432"))
	config.DB_USER = cast.ToString(coalesce("DB_USER", "postgres"))
	config.DB_NAME = cast.ToString(coalesce("DB_NAME", "ecommerce_auth_service"))
	config.DB_PASSWORD = cast.ToString(coalesce("DB_PASSWORD", "123321"))
	config.DB_CASBIN_DRIVER = cast.ToString(coalesce("DB_CASBIN_DRIVER", "postgres"))
	config.ACCESS_TOKEN = cast.ToString(coalesce("ACCESS_TOKEN", "key_is_really_easy"))
	config.REFRESH_TOKEN = cast.ToString(coalesce("REFRESH_TOKEN", "key_is_not_hand"))

	return config
}

func coalesce(env string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(env)
	if !exists {
		return defaultValue
	}
	return value
}
