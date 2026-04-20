package config

import "os"

const defaultBaseURL = "https://qa-internship.avito.com"

type Config struct {
	BaseURL string
}

func Load() *Config {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Config{BaseURL: baseURL}
}
