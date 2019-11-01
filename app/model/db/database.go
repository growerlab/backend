package db

import (
	"fmt"
	"runtime/debug"

	"github.com/Masterminds/squirrel"
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

	// pgsql placeholder
	squirrel.StatementBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return nil
}

func Transact(txFn func(tx *sqlx.Tx) error) (err error) {
	var tx *sqlx.Tx
	tx, err = DB.Beginx()
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
