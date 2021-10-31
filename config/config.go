package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/philandstuff/dhall-golang/v6"
)

type Config struct {
	ZapLogFile string `dhall:"zap_log_file"`
	GinLogFile string `dhall:"gin_log_file"`
}

type Environment struct {
	DBHost string
	DBPort string

	GinMode  string
	LogLevel string
}

var Conf *Config
var Env *Environment

func init() {
	if err := dhall.UnmarshalFile("config.dhall", &Conf); err != nil {
		panic(err)
	}

	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	Env = &Environment{
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),

		GinMode: os.Getenv("GIN_MODE"),

		LogLevel: os.Getenv("DEBUG_LEVEL"),
	}
}
