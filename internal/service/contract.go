package service

import (
	"github.com/nmluci/go-backend/internal/repository"
	"github.com/nmluci/go-backend/pkg/dto"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Ping() (pingResponse dto.PublicPingResponse)
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
