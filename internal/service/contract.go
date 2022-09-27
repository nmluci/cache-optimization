package service

import (
	"context"

	"github.com/nmluci/cache-optimization/internal/model"
	"github.com/nmluci/cache-optimization/internal/repository"
	"github.com/nmluci/cache-optimization/pkg/dto"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Ping() (pingResponse dto.PublicPingResponse)

	Register(ctx context.Context, payload *dto.PublicUserPayload) (err error)
	Login(ctx context.Context, payload *dto.PublicUserLoginPayload) (sessionKey string, err error)
	EditUser(ctx context.Context, id uint64, payload *dto.PublicUserPayload) (err error)
	DeleteUser(ctx context.Context, id uint64, sessionKey string) (err error)

	FindProductByID(ctx context.Context, id uint64) (res *model.Product, err error)
	FindProducts(ctx context.Context) (res []*model.Product, err error)
	InsertProduct(ctx context.Context, payload *model.Product) (err error)
	UpdateProduct(ctx context.Context, id uint64, payload *model.Product) (err error)
	DeleteProduct(ctx context.Context, id uint64) (err error)

	Checkout(ctx context.Context, payload *dto.PublicCheckout) (err error)

	// NO CACHE
	ForceLogin(ctx context.Context, payload *dto.PublicUserLoginPayload) (sessionKey string, err error)
	ForceEditUser(ctx context.Context, id uint64, payload *dto.PublicUserPayload) (err error)
	ForceFindProductByID(ctx context.Context, id uint64) (res *model.Product, err error)
	ForceFindProducts(ctx context.Context) (res []*model.Product, err error)
}

type service struct {
	logger     *logrus.Entry
	conf       *serviceConfig
	repository repository.Repository
}

type serviceConfig struct {
}

type NewServiceParams struct {
	Logger     *logrus.Entry
	Repository repository.Repository
}

func NewService(params *NewServiceParams) Service {
	return &service{
		logger:     params.Logger,
		conf:       &serviceConfig{},
		repository: params.Repository,
	}
}
