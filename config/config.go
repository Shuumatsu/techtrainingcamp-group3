package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/philandstuff/dhall-golang/v6"
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
	// mysql config
	DBUser         string
	DBPasswd       string
	DBHost         string
	DBPort         string
	DBName         string
	DBMaxIdleConns string
	DBMaxOpenConns string
	// redis config
	RedisPasswd   string
	RedisHost     string
	RedisPort     string
	RedisPoolSize string
	// kafka config
	KafkaHost   string
	KafkaPort   string
	KafkaTopics []string
	// tokenBucket config
	TokenInterval string
	TokenMaxCount string
	// gin config
	GinMode string
	// log config
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
		DBUser:         os.Getenv("DB_USER"),
		DBPasswd:       os.Getenv("DB_PASSWD"),
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBName:         os.Getenv("DB_NAME"),
		DBMaxIdleConns: os.Getenv("DB_MAX_IDLECONNS"),
		DBMaxOpenConns: os.Getenv("DB_MAX_OPENCONNS"),

		RedisPasswd:   os.Getenv("REDIS_PASSWD"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPoolSize: os.Getenv("REDIS_POOLSIZE"),

		KafkaHost:   os.Getenv("KAFKA_HOST"),
		KafkaPort:   os.Getenv("KAFKA_PORT"),
		KafkaTopics: []string{"OpenEnvelope", "AddUser", "AddEnvelopeToUser"},

		TokenInterval: os.Getenv("TOKEN_INTERVAL"),
		TokenMaxCount: os.Getenv("TOKEN_MAXCOUNT"),

		GinMode: os.Getenv("GIN_MODE"),

		LogLevel: os.Getenv("LOG_LEVEL"),
	}
}
