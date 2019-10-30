package graphql

import (
	"github.com/gin-gonic/gin"
	"github.com/growerlab/backend/app/common/env"
)

type Session struct {
	e   *env.Environment
	ctx *gin.Context
}

func NewSession(userToken string, ctx *gin.Context) *Session {
	e := env.NewEnvironment()
	e.Set(env.VarUserToken, userToken)
	return &Session{
		e:   e,
		ctx: ctx,
	}
}

func (s *Session) GetContext() *gin.Context {
	return s.ctx
}

func (s *Session) Env() *env.Environment {
	return s.e
}

func (s *Session) UserToken() string { // current user
	userID, _ := s.e.MustString(env.VarUserToken)
	return userID
}

func (s *Session) IsGuest() bool {
	token := s.UserToken()
	return len(token) == 0
}
