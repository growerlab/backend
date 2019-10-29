package graphql

import (
	"github.com/growerlab/backend/app/common/env"
)

type Session struct {
	e *env.Environment
}

func NewSession(userToken string) *Session {
	e := env.NewEnvironment()
	e.Set(env.VarUserToken, userToken)
	return &Session{
		e: e,
	}
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
