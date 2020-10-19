package db

type SqlAction struct {
	Action  string
	Tables  []string
	Columns map[string][]string      // table => columns
	Values  map[string][]interface{} // table => values
}

type Table struct {
	Name string
	As   string
}

type SqlParser struct {
	sql  string
	args []interface{}
}

func NewSqlParser(sql string, args ...interface{}) *SqlParser {
	return &SqlParser{
		sql:  sql,
		args: args,
	}
}
