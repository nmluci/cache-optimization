package component

import (
	"context"
	"fmt"
	"time"

	"github.com/nmluci/cache-optimization/internal/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type InitMongoDBParams struct {
	Conf   *config.MongoDBConfig
	Logger *logrus.Entry
}

var logTagInitMongoDB = "[InitMongoDB]"

func InitMongoDB(params *InitMongoDBParams) (db *mongo.Database, err error) {
	dataSource := fmt.Sprintf(
		"mongodb://%s:%s@%s",
		params.Conf.Username,
		params.Conf.Password,
		params.Conf.Address,
	)

	var client *mongo.Client
	for i := 10; i > 0; i-- {
		client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(dataSource))
		if err == nil {
			break
		}

		params.Logger.Errorf("%s error init connection to mongo for %s: %+v, retrying in 1 second", logTagInitMongoDB, dataSource, err)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return
	}

	for i := 20; i > 0; i-- {
		err = client.Ping(context.TODO(), readpref.PrimaryPreferred())
		if err == nil {
			break
		}

		params.Logger.Errorf("%s error ping mongo for %s: %+v, retrying in 1 second", logTagInitMongoDB, dataSource, err)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return
	}

	db = client.Database(params.Conf.DBName)
	params.Logger.Infof("%s mongo init successfully", logTagInitMongoDB)
	return
}
