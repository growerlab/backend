package graphql

import "github.com/growerlab/backend/app/common/env"

type Session struct {
	userID int
	e      env.Environment
}

func NewSession(userID int) *Session {
	e := env.Environment{}
	return &Session{
		userID: userID,
		e:      e,
	}
}

func (s *Session) Env() env.Environment {
	return s.e
}
func (s *Session) UserID() int { // current user
	return s.userID
}
