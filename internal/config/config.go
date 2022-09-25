package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName string      `json:"serviceName"`
	ServicePort string      `json:"servicePort"`
	Environment Environment `json:"environment"`

	MariaDBConfig MariaDBConfig `json:"mariaDBConfig"`
	MongoDBConfig MongoDBConfig `json:"mongoDBConfig"`
	RedisConfig   RedisConfig   `json:"redisConfig"`
}

const logTagConfig = "[Init Config]"

var config *Config

func Init() {
	godotenv.Load("conf/.env")

	conf := Config{
		ServiceName: os.Getenv("SERVICE_NAME"),
		ServicePort: os.Getenv("SERVICE_PORT"),
		MariaDBConfig: MariaDBConfig{
			Address:  fmt.Sprintf("%s:%s", os.Getenv("MARIADB_ADDRESS"), os.Getenv("MARIADB_PORT")),
			Username: os.Getenv("MARIADB_USERNAME"),
			Password: os.Getenv("MARIADB_PASSWORD"),
			DBName:   os.Getenv("MARIADB_DBNAME"),
		},
		MongoDBConfig: MongoDBConfig{
			Address:  fmt.Sprintf("%s:%s", os.Getenv("MONGODB_ADDRESS"), os.Getenv("MONGODB_PORT")),
			Username: os.Getenv("MONGODB_USERNAME"),
			Password: os.Getenv("MONGODB_PASSWORD"),
			DBName:   os.Getenv("MONGODB_DBNAME"),
		},
		RedisConfig: RedisConfig{
			Address:  os.Getenv("REDIS_ADDRESS"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
	}

	if conf.ServiceName == "" {
		log.Fatalf("%s service name should not be empty", logTagConfig)
	}

	if conf.ServicePort == "" {
		log.Fatalf("%s service port should not be empty", logTagConfig)
	}

	if conf.MariaDBConfig.Address == "" || conf.MariaDBConfig.DBName == "" {
		log.Fatalf("%s address and db name cannot be empty", logTagConfig)
	}

	envString := os.Getenv("ENVIRONMENT")
	if envString != "dev" && envString != "prod" {
		log.Fatalf("%s environment must be either dev or prod, found: %s", logTagConfig, envString)
	}

	conf.Environment = Environment(envString)

	config = &conf
}

func Get() (conf *Config) {
	conf = config
	return
}
