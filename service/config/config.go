package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	Host string
	Port string
}

var Env *Environment

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	Env = &Environment{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
	}
}
