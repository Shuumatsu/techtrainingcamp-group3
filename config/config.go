package config

import (
	"github.com/joho/godotenv"
	"github.com/philandstuff/dhall-golang/v6"
	"os"
)

const (
	TotalMoney      = 20000
	MaxMoney        = 500
	MinMoney        = 100
	TotalAmount     = 100
	SnatchProb      = 0.6
	PoolCapacity    = 128
	PoolWorkerNUM   = 10
	MaxSnatchAmount = 5
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

	RedisPasswd string
	RedisHost   string
	RedisPort   string

	GinMode string

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

		RedisPasswd: os.Getenv("REDIS_PASSWD"),
		RedisHost:   os.Getenv("REDIS_HOST"),
		RedisPort:   os.Getenv("REDIS_PORT"),

		GinMode: os.Getenv("GIN_MODE"),

		LogLevel: os.Getenv("LOG_LEVEL"),
	}
}
