package service

import (
	"context"

	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/db"
	"github.com/growerlab/backend/app/model/user"
	"github.com/growerlab/backend/app/service/graphql"
)

func CurrentUser(ctx context.Context) (u *user.User, userNSID *int64, err error) {
	session, ok := Session(ctx)
	if !ok {
		return nil, nil, errors.New(errors.Unauthorize())
	}
	userToken := session.UserToken()
	u, err = user.GetUserByUserToken(db.DB, userToken)
	if err != nil {
		return nil, nil, err
	}
	if u != nil {
		return u, &u.NamespaceID, nil
	}
	return nil, nil, nil
}

func Session(ctx context.Context) (*graphql.Session, bool) {
	sess, ok := ctx.Value("session").(*graphql.Session)
	return sess, ok
}
