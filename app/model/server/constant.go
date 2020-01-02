package server

import (
	sq "github.com/Masterminds/squirrel"
)

type statusType int

const (
	StatusClosed  statusType = 0 // 关闭的
	StatusNormal  statusType = 1 // 正常的
	StatusSuspend statusType = 2 // 暂停的
)

var (
	SqlNormal        = sq.Eq{"deleted_at": nil}
	SqlStatusClosed  = sq.Eq{"status": StatusClosed}
	SqlStatusNormal  = sq.Eq{"status": StatusNormal}
	SqlStatusSuspend = sq.Eq{"status": StatusSuspend}
	SqlStatusDeleted = sq.NotEq{"deleted_at": nil}
)
