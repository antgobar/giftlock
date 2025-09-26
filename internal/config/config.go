package config

import (
	"log"
	"os"
)

type Config struct {
	ServerAddr     string
	DatabaseUrl    string
	PublicPrefixes []string
}

func Load() *Config {
	return &Config{
		ServerAddr:  mustLoadEnv("SERVER_ADDR"),
		DatabaseUrl: mustLoadEnv("DATABASE_URL"),
		PublicPrefixes: []string{
			"/register",
			"/login",
			"/logout",
			"/static",
			"/favicon.ico",
			"/assets",
			"/api/login",
			"/api/register",
			"/api/logout",
		},
	}
}

func mustLoadEnv(envName string) string {
	env := os.Getenv(envName)
	if env == "" {
		log.Fatalf("missing environment variable: %s", envName)
	}
	return env
}
