package graphql

import "github.com/growerlab/backend/app/common/env"

type Session struct {
	e *env.Environment
}

func NewSession(userID int) *Session {
	e := env.NewEnvironment()
	return &Session{
		e: e,
	}
}

func (s *Session) Env() *env.Environment {
	return s.e
}
func (s *Session) UserID() int { // current user
	userID, _ := s.e.MustInt64(env.VarUserID)
	return int(userID)
}
