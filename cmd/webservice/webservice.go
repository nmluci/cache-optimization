package webservice

import (
	"github.com/labstack/echo/v4"
	"github.com/nmluci/cache-optimization/cmd/webservice/router"
	"github.com/nmluci/cache-optimization/internal/component"
	"github.com/nmluci/cache-optimization/internal/config"
	"github.com/nmluci/cache-optimization/internal/repository"
	"github.com/nmluci/cache-optimization/internal/service"
	"github.com/sirupsen/logrus"
)

const logTagStartWebservice = "[Start]"

func Start(conf *config.Config, logger *logrus.Entry) {
	db, err := component.InitMariaDB(&component.InitMariaDBParams{
		Conf:   &conf.MariaDBConfig,
		Logger: logger,
	})

	if err != nil {
		logger.Fatalf("%s initializing maria db: %+v", logTagStartWebservice, err)
	}

	redis, err := component.InitRedis(&component.InitRedisParams{
		Conf:   &conf.RedisConfig,
		Logger: logger,
	})

	if err != nil {
		logger.Fatalf("%s initalizing redis: %+v", logTagStartWebservice, err)
	}

	ec := echo.New()
	ec.HideBanner = true
	ec.HidePort = true

	repo := repository.NewRepository(&repository.NewRepositoryParams{
		Logger:  logger,
		MariaDB: db,
		Redis:   redis,
	})

	service := service.NewService(&service.NewServiceParams{
		Logger:     logger,
		Repository: repo,
	})

	router.Init(&router.InitRouterParams{
		Logger:  logger,
		Service: service,
		Ec:      ec,
		Conf:    conf,
	})

	logger.Infof("%s starting service, listening to port: %s", logTagStartWebservice, conf.ServicePort)

	if err := ec.Start(conf.ServicePort); err != nil {
		logger.Errorf("%s starting service, cause: %+v", logTagStartWebservice, err)
	}
}
