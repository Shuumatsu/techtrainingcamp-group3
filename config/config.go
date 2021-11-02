package config

import "os"

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

var Env *Environment

func init() {
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