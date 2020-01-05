package client

import (
	"encoding/json"
	"errors"

	"github.com/tidwall/gjson"
)

type GraphQLError struct {
	Message   string          `json:"message"`
	Locations json.RawMessage `json:"locations,omitempty"`
}

type Result struct {
	Data     json.RawMessage `json:"data"`
	DataPath gjson.Result
	Errors   []*GraphQLError `json:"errors,omitempty"`
}

type Request struct {
	Body      string                 // query/mutation body
	RepoName  string                 // repo name
	RepoPath  string                 // repo path
	Variables map[string]interface{} // query variables
}

func NewRequest(body string, repo *RepoContext, vars map[string]interface{}) *Request {
	return &Request{
		Body:      body,
		RepoName:  repo.Name,
		RepoPath:  repo.Path,
		Variables: vars,
	}
}

func (r *Request) Validate() error {
	if len(r.RepoName) == 0 {
		return errors.New("repo name is required")
	}
	if len(r.RepoPath) == 0 {
		return errors.New("repo path is required")
	}
	if len(r.Body) == 0 {
		return errors.New("body is required")
	}
	return nil
}

func (r *Request) RequestBody() map[string]interface{} {
	if r.Variables == nil {
		r.Variables = map[string]interface{}{}
	}
	r.Variables["path"] = r.RepoPath
	r.Variables["name"] = r.RepoName
	bodyMap := map[string]interface{}{
		"query":     r.Body,
		"variables": r.Variables,
	}
	return bodyMap
}

type APISubmitter interface {
	Query(req *Request) (*Result, error)
	Mutation(req *Request) (*Result, error)
}

type RepoContext struct {
	Name string // repo name in Path
	Path string // repo abs Path
}
