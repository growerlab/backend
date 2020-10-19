package db

import (
	"database/sql"
	"io"

	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/utils/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	ActionInsert = "insert"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

type Event struct {
	Table  string        // table
	Action string        //
	Sql    string        // sql
	Args   []interface{} // sql args
}

type EventProcessor interface {
	Label() string
	OnEvents([]*Event) error
}

type EventHub struct {
}

type Hooker interface {
	Release()
	AsyncProcess()
	SyncProcess() error
}

type HookQueryer interface {
	Query(query interface{}, args ...interface{}) (*sql.Rows, error)
	Queryx(query interface{}, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query interface{}, args ...interface{}) *sqlx.Row
	Exec(query interface{}, args ...interface{}) (sql.Result, error)
	Select(dest interface{}, query sq.Sqlizer) error
}

var _ SqlRunner = (*SqlRunnerWithHook)(nil)

type SqlRunnerWithHook struct {
	id     string
	query  Queryer
	events []*Event
}

func NewSqlRunnerWithHook(query Queryer) *SqlRunnerWithHook {
	return &SqlRunnerWithHook{
		id:    uuid.UUID(),
		query: query,
	}
}

func (s *SqlRunnerWithHook) Put(txID string, sqlizer sq.Sqlizer) {
	if len(txID) == 0 {
		return
	}
}

func (s *SqlRunnerWithHook) Release(txID string) {

}

func (s *SqlRunnerWithHook) AsyncProcess(txID string) {

}

func (s *SqlRunnerWithHook) SyncProcess(txID string) error {
	return nil
}
func (s *SqlRunnerWithHook) Debug() bool {
	panic("implement me")
}

func (s *SqlRunnerWithHook) Logger() io.Writer {
	panic("implement me")
}

func (s *SqlRunnerWithHook) Println(query interface{}, args ...interface{}) {
	panic("implement me")
}

func (s *SqlRunnerWithHook) Queryx(query interface{}, args ...interface{}) (*sqlx.Rows, error) {
	panic("implement me")
}

func (s *SqlRunnerWithHook) QueryRowx(query interface{}, args ...interface{}) *sqlx.Row {
	panic("implement me")
}

func (s *SqlRunnerWithHook) Exec(query interface{}, args ...interface{}) (sql.Result, error) {
	panic("implement me")
}

func (s *SqlRunnerWithHook) Query(query interface{}, args ...interface{}) (*sql.Rows, error) {
	panic("implement me")
}

func (s *SqlRunnerWithHook) Select(dest interface{}, query interface{}, args ...interface{}) error {
	sqlQuery, sqlArgs, err := s.parseQueryAndArgs(query, args)
	if err != nil {
		return err
	}

	err = sqlx.Select(s.query, dest, sqlQuery, sqlArgs)
	return errors.Wrap(err, errors.SQLError())
}

func (s *SqlRunnerWithHook) parseQueryAndArgs(query interface{}, args ...interface{}) (
	sqlQuery string, sqlArgs []interface{}, err error) {

	switch q := query.(type) {
	case sq.Sqlizer:
		sqlQuery, sqlArgs, err = q.ToSql()
		if err != nil {
			err = errors.Trace(err)
			return
		}

	case string:
		sqlQuery = q
		sqlArgs = args
	default:
		err = errors.Errorf("invalid query type '%v'", q)
		return
	}
	return
}
