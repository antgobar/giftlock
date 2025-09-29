package config

import (
	"encoding/json"
	"log"
	"os"
)

type configFile struct {
	ServerAddrEnv  string   `json:"serverAddressEnv"`
	DatabaseUrlEnv string   `json:"databaseUrlEnv"`
	PublicPrefixes []string `json:"publicPrefixes"`
}

type Config struct {
	ServerAddr     string
	DatabaseUrl    string
	PublicPrefixes []string
}

func Load() *Config {
	file, err := os.Open("config/app.json")
	if err != nil {
		log.Fatalf("failed to open config.json: %v", err)
	}
	defer file.Close()

	var cfgFile configFile
	if err := json.NewDecoder(file).Decode(&cfgFile); err != nil {
		log.Fatalf("failed to decode config.json: %v", err)
	}
	return &Config{
		ServerAddr:     mustLoadEnv(cfgFile.ServerAddrEnv),
		DatabaseUrl:    mustLoadEnv(cfgFile.DatabaseUrlEnv),
		PublicPrefixes: cfgFile.PublicPrefixes,
	}
}

func mustLoadEnv(envName string) string {
	env := os.Getenv(envName)
	if env == "" {
		log.Fatalf("missing environment variable: %s", envName)
	}
	return env
}
