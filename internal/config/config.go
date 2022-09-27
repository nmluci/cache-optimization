package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName string      `json:"serviceName"`
	ServicePort string      `json:"servicePort"`
	Environment Environment `json:"environment"`

	MariaDBConfig MariaDBConfig `json:"mariaDBConfig"`
	RedisConfig   RedisConfig   `json:"redisConfig"`

	MasterKey string `json:"-"`
}

const logTagConfig = "[Init Config]"

var config *Config

func Init() {
	godotenv.Load("conf/.env")

	conf := Config{
		ServiceName: os.Getenv("SERVICE_NAME"),
		ServicePort: os.Getenv("SERVICE_PORT"),
		MariaDBConfig: MariaDBConfig{
			Address:  os.Getenv("MARIADB_ADDRESS"),
			Username: os.Getenv("MARIADB_USERNAME"),
			Password: os.Getenv("MARIADB_PASSWORD"),
			DBName:   os.Getenv("MARIADB_DBNAME"),
		},
		RedisConfig: RedisConfig{
			Address:  os.Getenv("REDIS_ADDRESS"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		MasterKey: os.Getenv("MASTER_KEY"),
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
