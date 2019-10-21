package db

import (
	"database/sql"
	"fmt"
	"runtime/debug"

	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/utils/conf"
	"github.com/growerlab/backend/app/utils/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	DB *sqlx.DB
)

func InitDatabase() error {
	var config = conf.GetConf()
	var err error
	DB, err = sqlx.Connect("postgres", config.Database.URL)
	if err != nil {
		return errors.Sql(err)
	}
	return nil
}

func Transaction(txFn func(tx *sql.Tx) error) (err error) {
	var tx *sql.Tx
	tx, err = DB.Begin()
	if err != nil {
		return errors.Sql(err)
	}

	defer func() {
		if p := recover(); p != nil {
			logger.Warn("%s: %s", p, debug.Stack())
			switch p.(type) {
			case error:
				err = p.(error)
			default:
				err = fmt.Errorf("%s", p)
			}
		}
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = errors.Sql(tx.Commit())
	}()

	return txFn(tx)
}
