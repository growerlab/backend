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

func init() {
	// pgsql placeholder
	sq.StatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func InitDatabase() error {
	var err error
	var config = conf.GetConf()
	DB, err = DoInitDatabase(config.Database.URL, config.Debug)
	return err
}

func DoInitDatabase(databaseURL string, debug bool) (*DBQuery, error) {
	var err error
	var sqlxDB *sqlx.DB

	sqlxDB, err = sqlx.Connect("pgx", databaseURL)
	if err != nil {
		return nil, errors.Wrap(err, errors.SQLError())
	}

	d := &DBQuery{
		dbBase: &dbBase{
			Queryer: sqlxDB,
			debug:   debug,
			logger:  logger.LogWriter,
		},
		db: sqlxDB,
	}
	return d, nil
}

func Transact(txFn func(Queryer) error) (err error) {
	tx := DB.Begin()

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

type Queryer interface {
	sqlx.Queryer
	sqlx.Execer
}

type dbBase struct {
	Queryer

	// sql的操作
	*sqlEvent

	debug  bool
	logger io.Writer
}

func (d *dbBase) Println(query string, args ...interface{}) {
	if d.debug {
		_, _ = fmt.Fprint(d.logger, fmt.Sprintf("%c[%d;%d;%dm%s%c[0m ", 0x1B, 1, 0, 36, query, 0x1B))
		if len(args) > 0 {
			_, _ = fmt.Fprint(d.logger, args, "\n")
		} else {
			_, _ = fmt.Fprint(d.logger, "\n")
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
	return d.Queryer.Exec(query, args...)
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
			Queryer:  tx,
			debug:    d.dbBase.debug,
			logger:   d.dbBase.logger,
			sqlEvent: NewSqlEvent(),
		},
		tx: tx,
	}
}

var _ sqlx.Queryer = (*DBTx)(nil)
var _ sqlx.Execer = (*DBTx)(nil)

type DBTx struct {
	*dbBase
	tx Transaction
}

func (d *DBTx) Rollback() error {
	d.Println("ROLLBACK")
	return d.tx.Rollback()
}

func (d *DBTx) Commit() error {
	d.Println("COMMIT")
	return d.tx.Commit()
}
