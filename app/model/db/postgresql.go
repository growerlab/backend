package db

import (
	"database/sql"
	"fmt"
	"io"
	"runtime/debug"

	sq "github.com/Masterminds/squirrel"
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
		return errors.Wrap(err, errors.SqlError)
	}

	DB = &DBQuery{
		dbBase: &dbBase{
			Queryer: sqlxDB,
			Execer:  sqlxDB,
			debug:   config.Debug,
			logger:  logger.LogWriter,
		},
		db: sqlxDB,
	}

	// pgsql placeholder
	sq.StatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

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
			_ = tx.Rollback()
			return
		}
		err = errors.Trace(tx.Commit())
	}()

	return txFn(tx)
}

type Transaction interface {
	Rollback() error
	Commit() error
}

type dbBase struct {
	sqlx.Queryer
	sqlx.Execer

	debug  bool
	logger io.Writer
}

func (d *dbBase) Println(query string, args ...interface{}) {
	if d.debug {
		fmt.Fprint(d.logger, fmt.Sprintf("%c[%d;%d;%dm%s%c[0m ", 0x1B, 1, 0, 36, query, 0x1B))
		if len(args) > 0 {
			fmt.Fprint(d.logger, args, "\n")
		} else {
			fmt.Fprint(d.logger, "\n")
		}
	}
}

func (d *dbBase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	d.Println(query, args...)
	return d.Queryer.Query(query, args...)
}

func (d *dbBase) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	d.Println(query, args...)
	return d.Queryer.Queryx(query, args...)
}

func (d *dbBase) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	d.Println(query, args...)
	return d.Queryer.QueryRowx(query, args...)
}

func (d *dbBase) Exec(query string, args ...interface{}) (sql.Result, error) {
	d.Println(query, args...)
	return d.Execer.Exec(query, args...)
}

// 带日志输出的db封装
type DBQuery struct {
	*dbBase
	db *sqlx.DB
}

func (d *DBQuery) Begin() *DBTx {
	d.Println("BEGIN")
	tx := d.db.MustBegin()

	return &DBTx{
		dbBase: &dbBase{
			Queryer: tx,
			Execer:  tx,
			debug:   d.dbBase.debug,
			logger:  d.dbBase.logger,
		},
		Transaction: tx,
	}
}

type DBTx struct {
	Transaction
	*dbBase
}

var _ sqlx.Queryer = (*DBTx)(nil)
var _ sqlx.Execer = (*DBTx)(nil)

func (d *DBTx) Rollback() error {
	d.Println("ROLLBACK")
	return d.Transaction.Rollback()
}

func (d *DBTx) Commit() error {
	d.Println("COMMIT")
	return d.Transaction.Commit()
}
