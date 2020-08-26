package client

import (
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

const (
	defaultClientTimeout = 10 * time.Second
)

type Client struct {
	apiURL string
	client *http.Client
}

func NewClient(apiURL string, timeout time.Duration) (*Client, error) {
	if timeout == 0 {
		timeout = defaultClientTimeout
	}
	_, err := url.Parse(apiURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Client{
		apiURL: apiURL,
		client: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (c *Client) Query(req *Request) (*Result, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	return Post(c.client, c.apiURL, req.RequestBody())
}

func (c *Client) Mutation(req *Request) (*Result, error) {
	return c.Query(req)
}

func (c *Client) Branch(repo *RepoContext) *Branch {
	return &Branch{
		client: c,
		repo:   repo,
	}
}

func (c *Client) Repository(repo *RepoContext) *Repository {
	return &Repository{
		client: c,
		repo:   repo,
	}
}
