package repository

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
}

type repository struct {
	mariaDB *sql.DB
	mongoDB *mongo.Database
	redis   *redis.Client
	logger  *logrus.Entry
	conf    *repositoryConfig
}

type repositoryConfig struct {
}

type NewRepositoryParams struct {
	Logger  *logrus.Entry
	MariaDB *sql.DB
	MongoDB *mongo.Database
	Redis   *redis.Client
}

func NewRepository(params *NewRepositoryParams) Repository {
	return &repository{
		logger:  params.Logger,
		conf:    &repositoryConfig{},
		mariaDB: params.MariaDB,
		mongoDB: params.MongoDB,
		redis:   params.Redis,
	}
}
