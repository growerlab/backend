package db

import (
	"database/sql"
	"fmt"
	"io"
	"runtime/debug"

	"github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/utils/conf"
	"github.com/growerlab/backend/app/utils/logger"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	// 带sql日志输出的封装
	DB *DBQuery
)

func InitDatabase() error {
	var config = conf.GetConf()
	var err error
	var sqlxDB *sqlx.DB

	sqlxDB, err = sqlx.Connect("pgx", config.Database.URL)
	if err != nil {
		return errors.Sql(err)
	}

	DB = &DBQuery{
		sqlxDB: sqlxDB,
		debug:  config.Debug,
		logger: logger.LogWriter,
	}

	// pgsql placeholder
	squirrel.StatementBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return nil
}

func Transact(txFn func(*DBTx) error) (err error) {
	tx := DB.Begin()

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
			logger.Error("%+v", err)
			_ = tx.Rollback()
			return
		}
		err = errors.Sql(tx.Commit())
	}()

	return txFn(tx)
}

// 带日志输出的db封装
type DBQuery struct {
	sqlxDB *sqlx.DB

	debug  bool
	logger io.Writer
}

func (d *DBQuery) Println(query string, args ...interface{}) {
	if d.debug {
		fmt.Fprint(d.logger, fmt.Sprintf("%c[%d;%d;%dm%s%c[0m ", 0x1B, 1, 0, 36, query, 0x1B))
		if len(args) > 0 {
			fmt.Fprint(d.logger, args, "\n")
		} else {
			fmt.Fprint(d.logger, "\n")
		}
	}
}

func (d *DBQuery) Query(query string, args ...interface{}) (*sql.Rows, error) {
	d.Println(query, args...)
	return d.sqlxDB.Query(query, args...)
}

func (d *DBQuery) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	d.Println(query, args...)
	return d.sqlxDB.Queryx(query, args...)
}

func (d *DBQuery) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	d.Println(query, args...)
	return d.sqlxDB.QueryRowx(query, args...)
}

func (d *DBQuery) Exec(query string, args ...interface{}) (sql.Result, error) {
	d.Println(query, args...)
	return d.sqlxDB.Exec(query, args...)
}

func (d *DBQuery) Begin() *DBTx {
	d.Println("BEGIN")
	tx := d.sqlxDB.MustBegin()
	return &DBTx{
		Execer:  tx,
		Queryer: tx,
		tx:      tx,
		logger:  d.logger,
	}
}

type DBTx struct {
	sqlx.Execer
	sqlx.Queryer

	tx     *sqlx.Tx
	logger io.Writer
}

func (d *DBTx) Rollback() error {
	fmt.Fprintln(d.logger, "ROLLBACK")
	return d.tx.Rollback()
}

func (d *DBTx) Commit() error {
	fmt.Fprintln(d.logger, "COMMIT")
	return d.tx.Commit()
}
