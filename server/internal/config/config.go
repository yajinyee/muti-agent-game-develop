// Package config 伺服器設定
package config

import (
	"os"
	"strconv"
)

// Config 伺服器設定
type Config struct {
	Port            string
	MaxPlayersPerRoom int
	InitialCoins    int
	DebugMode       bool
}

// Load 從環境變數載入設定
func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "7777"),
		MaxPlayersPerRoom: getEnvInt("MAX_PLAYERS", 10),
		InitialCoins:    getEnvInt("INITIAL_COINS", 10000),
		DebugMode:       getEnv("DEBUG", "false") == "true",
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return defaultVal
}
