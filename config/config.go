package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBName                 string
	DBAddress              string
	JWTExpirationInSeconds int64
	JWTSecret              string
}

var Envs = NewConfig()

func NewConfig() *Config {
	return &Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "4000"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", "password"),
		DBName:                 getEnv("DB_NAME", "ecommerce"),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "3306")),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600),
		JWTSecret:              getEnv("JWT_SECRET", "JDNMRVJIRJVRRVERVMKVFJVJVFJNV"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			return fallback
		}
		return i
	}
	return fallback

}
