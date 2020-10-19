package permission

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/context"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/db"
)

var table = "permission"
var columns = []string{
	"id",
	"namespace_id",
	"code",
	"context_type",
	"context_param_1",
	"context_param_2",
	"user_domain_type",
	"user_domain_param",
	"created_at",
	"deleted_at",
}

func ListPermissionsByContext(src db.HookQueryer, code int, c *context.Context) ([]*Permission, error) {
	where := sq.And{
		sq.Eq{"code": code},
		sq.Eq{"context_type": c.Type},
		sq.Eq{"context_param_1": c.Param1},
		sq.Eq{"context_param_2": c.Param2},
	}
	return listPermissionByCond(src, columns, where)
}

func listPermissionByCond(src db.HookQueryer, cols []string, cond sq.Sqlizer) ([]*Permission, error) {
	sql := sq.Select(cols...).
		From(table).
		Where(cond)

	result := make([]*Permission, 0)
	err := src.Select(&result, sql)
	if err != nil {
		return nil, errors.Wrap(err, errors.SQLError())
	}
	return result, nil
}
