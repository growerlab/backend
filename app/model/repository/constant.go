package repository

type RepoStatus int

const (
	StatusPrivate RepoStatus = 0
	StatusPublic  RepoStatus = 1
	StatusAll     RepoStatus = -1
)
