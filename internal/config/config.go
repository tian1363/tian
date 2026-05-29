package config

import "os"

type Config struct {
	HTTPAddr string
	MySQLDSN string
}

func Load() Config {
	return Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		MySQLDSN: getEnv("MYSQL_DSN", "root:root@tcp(127.0.0.1:3306)/note?parseTime=true&charset=utf8mb4"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
