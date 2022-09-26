package component

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nmluci/cache-optimization/internal/config"
	"github.com/sirupsen/logrus"
)

type InitRedisParams struct {
	Conf   *config.RedisConfig
	Logger *logrus.Entry
}

const logTagInitRedis = "[InitRedis]"

func InitRedis(params *InitRedisParams) (db *redis.Client, err error) {
	dbSource := fmt.Sprintf("%s:%s", params.Conf.Address, params.Conf.Port)

	db = redis.NewClient(&redis.Options{
		Addr:     dbSource,
		Password: params.Conf.Password,
		DB:       0,
	})

	for i := 20; i > 0; i-- {
		_, err = db.Ping(context.TODO()).Result()
		if err == nil {
			break
		}

		params.Logger.Errorf("%s error init db: %+v, retrying in 1 second", logTagInitRedis, err)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return
	}

	params.Logger.Infof("%s redis init succesfully", logTagInitRedis)
	return
}
