package main

import (
	"os"

	"github.com/AndreyShep2012/go-company-handler/internal/app"
	"github.com/AndreyShep2012/go-company-handler/internal/config"
)

func main() {
	path := getEnv("CONFIG_PATH", "config.yml")
	cfg := initConfig(path)
	app.Serve(cfg)
}

func initConfig(path string) config.Config {
	cfg, err := config.Load(path)
	if err != nil {
		panic(err)
	}

	return cfg
}

func getEnv(key string, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}
