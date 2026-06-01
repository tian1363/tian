package config

import "os"

type Config struct {
	HTTPAddr string
	MySQLDSN string
}

func Load() Config {
	return Config{
		HTTPAddr: resolveHTTPAddr(),
		MySQLDSN: getEnv("MYSQL_DSN", "root:root@tcp(127.0.0.1:3306)/note?parseTime=true&charset=utf8mb4"),
	}
}

func resolveHTTPAddr() string {
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		return addr
	}

	// Railway and similar platforms provide PORT.
	if port := os.Getenv("PORT"); port != "" {
		if port[0] == ':' {
			return port
		}
		return ":" + port
	}

	return ":8080"
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
