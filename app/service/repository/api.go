package repository

import (
	"github.com/growerlab/backend/app/model/db"
	"github.com/growerlab/backend/app/model/repository"
	"github.com/growerlab/backend/app/model/server"
	"github.com/growerlab/codev-svc/client"
)

func NewApi(srv *server.Server, repo *repository.Repository) (*SVCApi, error) {
	return getClient(srv, repo)
}

func NewApiFromSrvID(srvID int64, repo *repository.Repository) (*SVCApi, error) {
	return getClientFromServerID(srvID, repo)
}

func getClientFromServerID(srvID int64, repo *repository.Repository) (*SVCApi, error) {
	srv, err := server.GetServer(db.DB, srvID)
	if err != nil {
		return nil, err
	}
	return getClient(srv, repo)
}

func getClient(srv *server.Server,  repo *repository.Repository) (*SVCApi, error) {
	c, err := client.NewClient(srv.URL(), 0) // default 10s timeout
	if err != nil {
		return nil, err
	}
	return &SVCApi{
		c: c,
		repo: &client.RepoContext{
			Path: repo.ServerPath,
			Name: repo.Path,
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
