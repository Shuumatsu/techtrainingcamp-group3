package config

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	TotalMoney      = 2e3
	MaxMoney        = 500
	MinMoney        = 100
	UserAmount      = 6e3
	SnatchProb      = 0.6
	MaxSnatchAmount = 5
	TotalAmount     = 7
)

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
	// profiler
	Profiler string
	// bloom filter
	Bloomfilter string
	// rpc server
	RpcHost string
	RpcPort string
	// http server
	HttpHost string
	HttpPort string
}

var Env *Environment

func init() {
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

		Profiler:    os.Getenv("PROFILER"),
		Bloomfilter: os.Getenv("BLOOMFILTER"),

		// rpc server
		RpcHost: os.Getenv("RPC_HOST"),
		RpcPort: os.Getenv("RPC_PORT"),
		// http server
		HttpHost: os.Getenv("HTTP_HOST"),
		HttpPort: os.Getenv("HTTP_PORT"),
	}
}
