package config

import (
	"log"
	"os"
)

func MustLoadEnv(envName string) string {
	env := os.Getenv(envName)
	if env == "" {
		log.Fatalf("missing environment variable: %s", envName)
	}
	return env
}
