package config

import "os"

type Config struct {
	Port              string
	DBUsername        string
	DBPassword        string
	DBHost            string
	DBPort            string
	BootstrapServers  string
	KafkaAuthPassword string
	JaegerHost        string
	LokiHost          string
}

func NewConfig() *Config {
	return &Config{
		Port:              os.Getenv("SERVICE_PORT"),
		DBUsername:        os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
		DBPassword:        os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
		DBHost:            os.Getenv("DB_HOST"),
		DBPort:            os.Getenv("DB_PORT"),
		BootstrapServers:  os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		KafkaAuthPassword: os.Getenv("KAFKA_AUTH_PASSWORD"),
		JaegerHost:        os.Getenv("JAEGER_ENDPOINT"),
		LokiHost:          os.Getenv("LOKI_ENDPOINT"),
	}
}
