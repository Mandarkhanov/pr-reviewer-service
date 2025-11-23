package repository

import (
	"context"
	"pr-reviewer/internal/domain"
)

type TeamRepository interface {
	Create(ctx context.Context, db Querier, team domain.Team) error
	GetByName(ctx context.Context, db Querier, name string) (*domain.Team, error)
}

type UserRepository interface {
	Upsert(ctx context.Context, db Querier, users []domain.User) error
	SetIsActive(ctx context.Context, db Querier, userID string, isActive bool) (*domain.User, error)
	GetByID(ctx context.Context, db Querier, userID string) (*domain.User, error)
	GetActiveCandidates(ctx context.Context, db Querier, teamName string, excludeUserIDs []string) ([]domain.User, error)
}

type PullRequestRepository interface {
	Exists(ctx context.Context, db Querier, id string) (bool, error)
	Create(ctx context.Context, db Querier, pr domain.PullRequest) error
	GetByID(ctx context.Context, db Querier, id string) (*domain.PullRequest, error)
	SetStatus(ctx context.Context, db Querier, id string, status domain.PRStatus) error
	ReplaceReviewer(ctx context.Context, db Querier, prID, oldReviewerID, newReviewerID string) error
	GetByReviewerID(ctx context.Context, db Querier, reviewerID string) ([]domain.PullRequestShort, error)
}
