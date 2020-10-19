package context

import (
	"github.com/growerlab/backend/app/model/db"
)

type Context struct {
	Type   int   `json:"type"`
	Param1 int64 `json:"param1"`
	Param2 int64 `json:"param2"`
}

type DBContext struct {
	Src   db.HookQueryer
	MemDB *db.MemDBClient
}
