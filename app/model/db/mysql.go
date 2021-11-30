package db

import (
	"database/sql"
	"fmt"
	"io"
	"runtime/debug"

	_ "github.com/go-sql-driver/mysql"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/utils/conf"
	"github.com/growerlab/backend/app/utils/logger"
	"github.com/jmoiron/sqlx"
)

var (
	// DB 带sql日志输出的封装
	DB *DBQuery
)

func InitDatabase() error {
	var err error
	var config = conf.GetConf()
	DB, err = DoInitDatabase(config.Database.URL, config.Debug)
	return err
}

func DoInitDatabase(databaseURL string, debug bool) (*DBQuery, error) {
	var err error
	var sqlxDB *sqlx.DB

	sqlxDB, err = sqlx.Connect("mysql", databaseURL)
	if err != nil {
		return nil, errors.Wrap(err, errors.SQLError())
	}

	d := &DBQuery{
		Ext:    sqlxDB,
		debug:  debug,
		logger: logger.LogWriter,
	}
	return d, nil
}

func Transact(txFn func(tx sqlx.Ext) error) (err error) {
	d := DB.Ext.(*sqlx.DB)
	txa := d.MustBegin()

	defer func() {
		if p := recover(); p != nil {
			logger.Warn("%s: %s", p, debug.Stack())
			switch x := p.(type) {
			case error:
				err = x
			default:
				err = fmt.Errorf("%s", x)
			}
		}
		if err != nil {
			DB.Println("rollback")
			_ = txa.Rollback()
			return
		}
		DB.Println("commit")
		err = errors.Trace(txa.Commit())
	}()
	return txFn(txa)
}

type DBQuery struct {
	sqlx.Ext

	debug  bool
	logger io.Writer
}

func (d *DBQuery) Println(query string, args ...interface{}) {
	if d.debug {
		_, _ = fmt.Fprint(d.logger, fmt.Sprintf("%c[%d;%d;%dm%s%c[0m ", 0x1B, 1, 0, 36, query, 0x1B))
		if len(args) > 0 {
			_, _ = fmt.Fprint(d.logger, args, "\n")
		} else {
			_, _ = fmt.Fprint(d.logger, "\n")
		}
	}
}

func (d *DBQuery) Query(query string, args ...interface{}) (*sql.Rows, error) {
	d.Println(query, args...)
	return d.Ext.Query(query, args...)
}

func (d *DBQuery) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	d.Println(query, args...)
	return d.Ext.Queryx(query, args...)
}

func (d *DBQuery) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	d.Println(query, args...)
	return d.Ext.QueryRowx(query, args...)
}

func (d *DBQuery) Exec(query string, args ...interface{}) (sql.Result, error) {
	d.Println(query, args...)
	return d.Ext.Exec(query, args...)
}
