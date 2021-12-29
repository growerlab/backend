package base

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/jmoiron/sqlx"
)

type Model struct {
	src sqlx.Ext

	table   string
	alias   string
	columns []string

	defaultTerms       sq.And
	ignoreDefaultTerms bool
}

func NewModel(src sqlx.Ext, table string) *Model {
	return &Model{
		src:   src,
		table: table,
	}
}

func (m *Model) Alias(a string) *Model {
	m.alias = a
	return m
}

func (m *Model) Columns(columns []string) *Model {
	m.columns = columns
	return m
}

func (m *Model) IgnoreDefaultTerms() {
	m.ignoreDefaultTerms = true
}

func (m *Model) Select(term sq.Sqlizer) *Selector {
	where := sq.And{}
	if m.defaultTerms != nil && !m.ignoreDefaultTerms {
		where = append(where, m.defaultTerms...)
	}
	if term != nil {
		where = append(where, term)
	}

	table := m.getTable()
	builder := sq.Select(m.columns...).From(table).Where(where)

	return &Selector{
		src:     m.src,
		builder: builder,
	}
}

func (m *Model) getTable() string {
	if len(m.alias) > 0 {
		return fmt.Sprintf("%s AS %s", m.table, m.alias)
	}
	return m.table
}

func (m *Model) Update(set map[string]interface{}, term sq.Sqlizer) *Updater {
	if len(set) == 0 {
		panic("'set' must required")
	}

	where := sq.And{}
	where = append(where, term)

	builder := sq.Update(m.table).SetMap(set).Where(where)

	return &Updater{
		src:     m.src,
		table:   m.table,
		values:  set,
		builder: builder,
	}
}

func (m *Model) Insert(values []interface{}) *Inserter {
	if len(values) == 0 {
		panic("'values' must required")
	}

	builder := sq.Insert(m.table).Columns(m.columns...)

	return &Inserter{
		src:     m.src,
		table:   m.table,
		values:  values,
		builder: builder,
	}
}

// BatchInsert 批量插入，不会触发hook
func (m *Model) BatchInsert(size int, getValuesFn func(int) []interface{}) error {
	const maxValues = 1000
	valueBucket := make([][]interface{}, 0, maxValues)

	batchInsertFunc := func(mulValues [][]interface{}) error {
		builder := sq.Insert(m.table).Columns(m.columns...)
		for _, values := range mulValues {
			builder = builder.Values(values...)
		}
		_, err := builder.RunWith(m.src).Exec()
		return errors.Wrap(err, errors.SQLError())
	}

	for i := 0; i < size; i++ {
		values := getValuesFn(i)
		valueBucket = append(valueBucket, values)
		if len(valueBucket) >= maxValues || i == size-1 {
			if err := batchInsertFunc(valueBucket); err != nil {
				return err
			}
			valueBucket = valueBucket[:0]
		}
	}

	return nil
}

type Inserter struct {
	src     sqlx.Ext
	table   string
	columns []string
	values  []interface{}
	builder sq.InsertBuilder
}

func (i *Inserter) Exec() error {

	// 单个插入数据
	set := make(map[string]interface{})
	for v := range i.values {
		set[i.columns[v]] = i.values[v]
	}
	// hook before
	err := hook.Effect(i.src, i.table, ActionCreate, TenseBefore, set)
	if err != nil {
		return errors.Trace(err)
	}

	i.builder = i.builder.Values(i.values...)
	_, err = i.builder.RunWith(i.src).Exec()
	if err != nil {
		return errors.Wrap(err, "insert error")
	}

	// hook after
	err = hook.Effect(i.src, i.table, ActionCreate, TenseAfter, set)
	return errors.Trace(err)
}

type Updater struct {
	src     sqlx.Ext
	table   string
	values  map[string]interface{}
	builder sq.UpdateBuilder
}

func (u *Updater) Exec() error {
	// hook before
	err := hook.Effect(u.src, u.table, ActionUpdate, TenseBefore, u.values)
	if err != nil {
		return errors.Trace(err)
	}

	_, err = u.builder.RunWith(u.src).Exec()
	if err != nil {
		return errors.Wrap(err, errors.SQLError())
	}

	// hook after
	err = hook.Effect(u.src, u.table, ActionUpdate, TenseAfter, u.values)
	return errors.Trace(err)
}

type Selector struct {
	src     sqlx.Ext
	builder sq.SelectBuilder
}

func (s *Selector) BuildSQL(fn func(builder sq.SelectBuilder) sq.SelectBuilder) *Selector {
	s.builder = fn(s.builder)
	return s
}

func (s *Selector) Query(dest interface{}) error {
	query, args, err := s.builder.ToSql()
	if err != nil {
		return errors.Trace(err)
	}
	err = sqlx.Select(s.src, dest, query, args...)
	return errors.Wrap(err, errors.SQLError())
}
