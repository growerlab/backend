package db

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/jmoiron/sqlx"
)

type Queryer interface {
	sqlx.Queryer
	sqlx.Execer
}

var _ Queryer = (*DBBase)(nil)

type DBBase struct {
	Queryer //

	debug  bool      //
	logger io.Writer //
}

func (d *DBBase) Logger() io.Writer {
	return d.logger
}

func (d *DBBase) Debug() bool {
	return d.debug
}

func (d *DBBase) Println(query string, args ...interface{}) {
	if d.debug {
		_, _ = fmt.Fprint(d.logger, fmt.Sprintf("%c[%d;%d;%dm%s%c[0m ", 0x1B, 1, 0, 36, query, 0x1B))
		if len(args) > 0 {
			_, _ = fmt.Fprint(d.logger, args, "\n")
		} else {
			_, _ = fmt.Fprint(d.logger, "\n")
		}
	}
}

func (d *DBBase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	d.Println(query, args...)
	return d.Queryer.Query(query, args...)
}

func (d *DBBase) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	d.Println(query, args...)
	return d.Queryer.Queryx(query, args...)
}

func (d *DBBase) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	d.Println(query, args...)
	return d.Queryer.QueryRowx(query, args...)
}

func (d *DBBase) Exec(query string, args ...interface{}) (sql.Result, error) {
	d.Println(query, args...)
	return d.Queryer.Exec(query, args...)
}
