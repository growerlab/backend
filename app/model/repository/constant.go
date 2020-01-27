package repository

type RepoState int

const (
	StatePrivate RepoState = 0
	StatePublic  RepoState = 1
	StateAll     RepoState = -1
)
