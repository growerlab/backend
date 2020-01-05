package client

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type Branch struct {
	client APISubmitter
	repo   *RepoContext
}

func (b *Branch) Info() (defaultBranch string, branches []string, err error) {
	body := `
{
	repo {
		default_branch {
			name
		}
		branches {
			name
		}
	}
}`
	var result *Result
	result, err = b.client.Query(NewRequest(body, b.repo, nil))
	if err != nil {
		return
	}

	branches = make([]string, 0)
	result.DataPath.Get("repo.branches.#.name").ForEach(func(_, value gjson.Result) bool {
		branches = append(branches, value.String())
		return true
	})

	defaultBranch = result.DataPath.Get("repo.default_branch.name").String()
	return
}

func (b *Branch) Delete(branchName string) (err error) {
	body := `
mutation DeleteBranch {
	deleteBranch(branchName: "%s")
}`
	body = fmt.Sprintf(body, branchName)

	var result *Result
	result, err = b.client.Mutation(NewRequest(body, b.repo, nil))
	if err != nil {
		return
	}
	if !result.DataPath.Get("deleteBranch").Bool() {
		err = errors.New("delete branch was faild")
		return
	}
	return
}
