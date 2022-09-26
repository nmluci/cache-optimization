package component

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nmluci/cache-optimization/internal/config"

	"github.com/sirupsen/logrus"
)

type InitMariaDBParams struct {
	Conf   *config.MariaDBConfig
	Logger *logrus.Entry
}

const logTagInitMariaDB = "[InitMariaDB]"

func InitMariaDB(params *InitMariaDBParams) (db *sql.DB, err error) {
	dataSource := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true",
		params.Conf.Username, params.Conf.Password,
		params.Conf.Address, params.Conf.DBName,
	)

	for i := 10; i > 0; i-- {
		db, err = sql.Open("mysql", dataSource)
		if err == nil {
			break
		}

		params.Logger.Errorf("%s error init opening db for %s: %+v, retrying in 1 second", logTagInitMariaDB, dataSource, err)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return
	}

	for i := 20; i > 0; i-- {
		err = db.Ping()
		if err == nil {
			break
		}

		params.Logger.Errorf("%s error ping db for %s: %+v, retrying in 1 second", logTagInitMariaDB, dataSource, err)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return
	}

	params.Logger.Infof("%s db init successfully", logTagInitMariaDB)
	return
}
