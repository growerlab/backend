package resolver

import (
	"context"

	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/user"
	"github.com/growerlab/backend/app/service/graphql"
)

var (
	ErrNotFoundSession = errors.New("not found session")
)

func CurrentUser(ctx context.Context) (*user.User, error) {
	session, ok := Session(ctx)
	if ok {
		return nil, ErrNotFoundSession
	}
	_ = session.UserToken()

	return nil, nil
}

func Session(ctx context.Context) (*graphql.Session, bool) {
	sess, ok := ctx.Value("session").(*graphql.Session)
	return sess, ok
}
