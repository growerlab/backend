package client

import (
	"errors"
	"fmt"
)

type Repository struct {
	client APISubmitter
	repo   *RepoContext
}

func (r *Repository) Create() (err error) {
	body := `
mutation CreateRepo {
	createRepo(path: "%s", name: "%s") {
		name
	}
}
`
	body = fmt.Sprintf(body, r.repo.Path, r.repo.Name)
	result, err := r.client.Mutation(NewRequest(body, r.repo, nil))
	if err != nil {
		return
	}

	if !result.DataPath.Get("createRepo.name").Exists() {
		return errors.New("create repo faild")
	}
	return nil
}
