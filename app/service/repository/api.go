package repository

import (
	"github.com/growerlab/backend/app/model/db"
	"github.com/growerlab/backend/app/model/server"
	"github.com/growerlab/codev-svc/client"
)

func NewApi(srv *server.Server, repoPath, repoName string) (*SVCApi, error) {
	return getClient(srv, repoName, repoName)
}

func NewApiFromSrvID(srvID int64, repoPath, repoName string) (*SVCApi, error) {
	return getClientFromServerID(srvID, repoPath, repoName)
}

func getClientFromServerID(srvID int64, repoPath, repoName string) (*SVCApi, error) {
	srv, err := server.GetServer(db.DB, srvID)
	if err != nil {
		return nil, err
	}
	return getClient(srv, repoPath, repoName)
}

func getClient(srv *server.Server, repoPath, repoName string) (*SVCApi, error) {
	c, err := client.NewClient(srv.URL(), 0) // default 10s timeout
	if err != nil {
		return nil, err
	}
	return &SVCApi{
		c: c,
		repo: &client.RepoContext{
			Path: repoPath,
			Name: repoName,
		},
	}, nil
}

type SVCApi struct {
	c    *client.Client
	repo *client.RepoContext
}

func (s *SVCApi) Repository() *client.Repository {
	return s.c.Repository(s.repo)
}
