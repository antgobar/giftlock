package config

import (
	"log"
	"os"
)

type Config struct {
	ServerAddr  string
	DatabaseUrl string
}

func Load() *Config {
	return &Config{
		ServerAddr:  mustLoadEnv("SERVER_ADDR"),
		DatabaseUrl: mustLoadEnv("DATABASE_URL"),
	}
}

func mustLoadEnv(envName string) string {
	env := os.Getenv(envName)
	if env == "" {
		log.Fatalf("missing environment variable: %s", envName)
	}
	return env
}
