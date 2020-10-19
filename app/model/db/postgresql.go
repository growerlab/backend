package db

import (
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

	db := &DBBase{
		Queryer: sqlxDB,
		debug:   debug,
		logger:  logger.LogWriter,
	}

	d := &DBQuery{
		SqlRunner: NewSqlRunnerWithHook(db),
		db:        sqlxDB,
	}
	return d, nil
}

func Transact(txFn func(SqlRunner) error) (err error) {
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
		err = tx.Commit()
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	return txFn(tx)
}

type Transaction interface {
	Rollback() error
	Commit() error
}

type SqlRunner interface {
	HookQueryer
	Debug() bool
	Logger() io.Writer
	Println(query interface{}, args ...interface{})
}

// 带日志输出的db封装
type DBQuery struct {
	SqlRunner
	db *sqlx.DB
}

func (d *DBQuery) Begin() *DBTx {
	d.Println("BEGIN")
	tx := d.db.MustBegin()

	db := &DBBase{
		Queryer: tx,
		debug:   d.SqlRunner.Debug(),
		logger:  d.SqlRunner.Logger(),
	}

	return &DBTx{
		SqlRunner: NewSqlRunnerWithHook(db),
		tx:        tx,
	}
}

type DBTx struct {
	SqlRunner
	tx Transaction
}

func (d *DBTx) Rollback() error {
	d.SqlRunner.Println("ROLLBACK")

	if hooker, ok := d.SqlRunner.(Hooker); ok {
		defer hooker.Release()
	}
	return d.tx.Rollback()
}

func (d *DBTx) Commit() error {
	d.SqlRunner.Println("COMMIT")

	if hooker, ok := d.SqlRunner.(Hooker); ok {
		defer hooker.Release()
		defer hooker.AsyncProcess()
		err := hooker.SyncProcess()
		if err != nil {
			return err
		}
	}
	return d.tx.Commit()
}
