package service

import (
	"database/sql"
	"errors"
	"pr-reviewer/internal/repository"
)

var (
	ErrTeamExists     = errors.New("team already exists")
	ErrTeamNotFound   = errors.New("team not found")
	ErrUserNotFound   = errors.New("user not found")
	ErrPRExists       = errors.New("pull request already exists")
	ErrAuthorNotFound = errors.New("author not found")
	ErrPRNotFound     = errors.New("pull request not found")
	ErrPRMerged       = errors.New("cannot reassign on merged PR")
	ErrNotAssigned    = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate    = errors.New("no active replacement candidate in team")
)

type Service struct {
	db        *sql.DB
	repoTeams repository.TeamRepository
	repoUsers repository.UserRepository
	repoPR    repository.PullRequestRepository
}

func NewService(
	db *sql.DB,
	repoTeams repository.TeamRepository,
	repoUsers repository.UserRepository,
	repoPR repository.PullRequestRepository,
) *Service {
	return &Service{
		db:        db,
		repoTeams: repoTeams,
		repoUsers: repoUsers,
		repoPR:    repoPR,
	}
}
