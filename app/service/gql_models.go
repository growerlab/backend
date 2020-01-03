// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package service

type ActivateCodePayload struct {
	Code string `json:"code"`
}

type LoginUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewRepository struct {
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	Public bool   `json:"public"`
}

type NewUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type Result struct {
	Ok bool `json:"OK"`
}

type UserToken struct {
	Token string `json:"token"`
}
