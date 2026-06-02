package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config contient toutes les variables de configuration de l'application.
// Toutes les valeurs sont chargées depuis les variables d'environnement au démarrage.
type Config struct {
	Port                   string
	Env                    string
	DatabasePath           string
	SupabaseURL            string
	SupabaseServiceRoleKey string
	ResendAPIKey           string
	JWTSecret              string
}

// Load charge la configuration depuis le fichier .env (si présent) et les variables d'environnement.
// Panique si une variable obligatoire est absente.
func Load() *Config {

	_ = godotenv.Load()

	cfg := &Config{
		Port:                   getEnvOrDefault("PORT", "8080"),
		Env:                    getEnvOrDefault("ENV", "development"),
		DatabasePath:           getEnvOrDefault("DATABASE_PATH", "pfe.db"),
		SupabaseURL:            mustGetEnv("SUPABASE_URL"),
		SupabaseServiceRoleKey: mustGetEnv("SUPABASE_SERVICE_ROLE_KEY"),
		ResendAPIKey:           mustGetEnv("RESEND_API_KEY"),
		JWTSecret:              getEnvOrDefault("JWT_SECRET", "dev-secret"),
	}

	cfg.validate()
	return cfg
}

func (c *Config) validate() {
	if c.Env != "development" && c.Env != "production" && c.Env != "test" {
		panic(fmt.Sprintf("ENV invalide : %s (attendu : development | production | test)", c.Env))
	}
}

// IsDevelopment retourne true si l'environnement est development.
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction retourne true si l'environnement est production.
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("Variable d'environnement obligatoire manquante : %s", key))
	}
	return val
}

func getEnvOrDefault(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
