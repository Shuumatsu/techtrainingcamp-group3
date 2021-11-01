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
	// database config
	DBUser   string
	DBPasswd string
	DBHost   string
	DBPort   string
	DBName   string

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
		DBUser:   os.Getenv("DB_USER"),
		DBPasswd: os.Getenv("DB_PASSWD"),
		DBHost:   os.Getenv("DB_HOST"),
		DBPort:   os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),

		GinMode: os.Getenv("GIN_MODE"),

		LogLevel: os.Getenv("LOG_LEVEL"),
	}
}
