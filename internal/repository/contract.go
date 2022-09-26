package repository

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/nmluci/cache-optimization/internal/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Repository interface {
	NewUserSession(ctx context.Context, data *model.Users) (sessionKey string, err error)
	InvalidateUserSession(ctx context.Context, sessionKey string) (err error)
	FindUserByID(ctx context.Context, id uint64) (res *model.Users, err error)
	FindUserByEmail(ctx context.Context, email string) (res *model.Users, err error)
	InsertNewUser(ctx context.Context, data *model.Users) (err error)
	UpdateUserByID(ctx context.Context, id uint64, data *model.Users) (err error)
	DeleteUserByID(ctx context.Context, id uint64) (err error)

	FindProducts(ctx context.Context) (res []*model.Product, err error)
	FindProductByID(ctx context.Context, id uint64) (res *model.Product, err error)
	InsertNewProduct(ctx context.Context, data *model.Product) (err error)
	UpdateProduct(ctx context.Context, id uint64, data *model.Product) (err error)
	DeleteProductByID(ctx context.Context, id uint64) (err error)
}

type repository struct {
	mariaDB *sql.DB
	redis   *redis.Client
	logger  *logrus.Entry
	conf    *repositoryConfig
}

type repositoryConfig struct {
}

type NewRepositoryParams struct {
	Logger  *logrus.Entry
	MariaDB *sql.DB
	Redis   *redis.Client
}

func NewRepository(params *NewRepositoryParams) Repository {
	return &repository{
		logger:  params.Logger,
		conf:    &repositoryConfig{},
		mariaDB: params.MariaDB,
		redis:   params.Redis,
	}
}
