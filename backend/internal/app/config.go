package app

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Address            string `yaml:"address"`
	DatabaseURL        string `yaml:"database_url"`
	FrontendOrigin     string `yaml:"frontend_origin"`
	FrontendBaseURL    string `yaml:"frontend_base_url"`
	FrontendDist       string `yaml:"frontend_dist"`
	CookieSecure       bool   `yaml:"cookie_secure"`
	GoogleClientID     string `yaml:"google_client_id"`
	GoogleClientSecret string `yaml:"google_client_secret"`
	GoogleRedirectURL  string `yaml:"google_redirect_url"`
}

func DefaultConfig() Config {
	return Config{
		Address:           "localhost:8080",
		DatabaseURL:       "file:./data/formbuilder.db?_pragma=foreign_keys(1)",
		FrontendOrigin:    "http://localhost:5173",
		FrontendBaseURL:   "http://localhost:5173",
		FrontendDist:      "../frontend/dist",
		GoogleRedirectURL: "http://localhost:8080/api/auth/google/callback",
	}
}

func LoadConfig(args []string) (Config, error) {
	cfg := DefaultConfig()

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	configPath := fs.String("config", "", "path to YAML config file")
	address := fs.String("address", "", "HTTP listen address")
	databaseURL := fs.String("database-url", "", "database URL, SQLite or Postgres")
	frontendOrigin := fs.String("frontend-origin", "", "allowed frontend origin")
	frontendBaseURL := fs.String("frontend-base-url", "", "public frontend base URL")
	frontendDist := fs.String("frontend-dist", "", "frontend build directory to serve")
	cookieSecure := fs.Bool("cookie-secure", false, "set Secure on auth cookies")
	googleClientID := fs.String("google-client-id", "", "Google OAuth client ID")
	googleClientSecret := fs.String("google-client-secret", "", "Google OAuth client secret")
	googleRedirectURL := fs.String("google-redirect-url", "", "Google OAuth callback URL")
	if err := fs.Parse(args); err != nil {
		return cfg, err
	}

	if *configPath != "" {
		data, err := os.ReadFile(*configPath)
		if err != nil {
			return cfg, fmt.Errorf("read config: %w", err)
		}
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return cfg, fmt.Errorf("parse config: %w", err)
		}
	}

	applyEnv(&cfg)
	if *address != "" {
		cfg.Address = *address
	}
	if *databaseURL != "" {
		cfg.DatabaseURL = *databaseURL
	}
	if *frontendOrigin != "" {
		cfg.FrontendOrigin = *frontendOrigin
	}
	if *frontendBaseURL != "" {
		cfg.FrontendBaseURL = *frontendBaseURL
	}
	if *frontendDist != "" {
		cfg.FrontendDist = *frontendDist
	}
	if fs.Lookup("cookie-secure").Value.String() == "true" {
		cfg.CookieSecure = *cookieSecure
	}
	if *googleClientID != "" {
		cfg.GoogleClientID = *googleClientID
	}
	if *googleClientSecret != "" {
		cfg.GoogleClientSecret = *googleClientSecret
	}
	if *googleRedirectURL != "" {
		cfg.GoogleRedirectURL = *googleRedirectURL
	}

	if cfg.Address == "" {
		return cfg, fmt.Errorf("address is required")
	}
	if cfg.DatabaseURL == "" {
		return cfg, fmt.Errorf("database_url is required")
	}
	if cfg.FrontendOrigin == "" {
		cfg.FrontendOrigin = cfg.FrontendBaseURL
	}
	return cfg, nil
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("ADDRESS"); v != "" {
		cfg.Address = v
	} else if v := os.Getenv("PORT"); v != "" {
		cfg.Address = "0.0.0.0:" + v
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		cfg.DatabaseURL = v
	}
	if v := os.Getenv("FRONTEND_ORIGIN"); v != "" {
		cfg.FrontendOrigin = v
	}
	if v := os.Getenv("FRONTEND_BASE_URL"); v != "" {
		cfg.FrontendBaseURL = v
	}
	if v := os.Getenv("FRONTEND_DIST"); v != "" {
		cfg.FrontendDist = v
	}
	if v := os.Getenv("COOKIE_SECURE"); v == "true" || v == "1" {
		cfg.CookieSecure = true
	}
	if v := os.Getenv("GOOGLE_CLIENT_ID"); v != "" {
		cfg.GoogleClientID = v
	}
	if v := os.Getenv("GOOGLE_CLIENT_SECRET"); v != "" {
		cfg.GoogleClientSecret = v
	}
	if v := os.Getenv("GOOGLE_REDIRECT_URL"); v != "" {
		cfg.GoogleRedirectURL = v
	}
}
