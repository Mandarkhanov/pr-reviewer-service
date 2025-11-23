package service

import (
	"context"
	"math/rand"
	"pr-reviewer/internal/domain"
	"time"
)

func (s *Service) CreatePR(ctx context.Context, prID, prName, authorID string) (*domain.PullRequest, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	exists, err := s.repoPR.Exists(ctx, tx, prID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPRExists
	}

	author, err := s.repoUsers.GetByID(ctx, tx, authorID)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, ErrAuthorNotFound
	}

	candidates, err := s.repoUsers.GetActiveCandidates(ctx, tx, author.TeamName, []string{authorID})
	if err != nil {
		return nil, err
	}

	reviewers := selectRandomReviewers(candidates, 2)

	pr := domain.PullRequest{
		ID:        prID,
		Name:      prName,
		AuthorID:  authorID,
		Status:    domain.PRStatusOpen,
		Reviewers: reviewers,
	}

	if err := s.repoPR.Create(ctx, tx, pr); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &pr, nil
}

func selectRandomReviewers(candidates []domain.User, limit int) []domain.User {
	if len(candidates) == 0 {
		return []domain.User{}
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	if len(candidates) > limit {
		return candidates[:limit]
	}
	return candidates
}

func (s *Service) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	pr, err := s.repoPR.GetByID(ctx, tx, prID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, ErrPRNotFound
	}

	if pr.Status == domain.PRStatusMerged {
		return pr, nil
	}

	if err := s.repoPR.SetStatus(ctx, tx, prID, domain.PRStatusMerged); err != nil {
		return nil, err
	}

	pr.Status = domain.PRStatusMerged
	now := time.Now()
	pr.MergedAt = &now

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *Service) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*domain.PullRequest, *domain.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	pr, err := s.repoPR.GetByID(ctx, tx, prID)
	if err != nil {
		return nil, nil, err
	}
	if pr == nil {
		return nil, nil, ErrPRNotFound
	}

	if pr.Status == domain.PRStatusMerged {
		return nil, nil, ErrPRMerged
	}

	isAssigned := false
	currentReviewerIDs := []string{pr.AuthorID}

	for _, r := range pr.Reviewers {
		currentReviewerIDs = append(currentReviewerIDs, r.ID)
		if r.ID == oldUserID {
			isAssigned = true
		}
	}

	if !isAssigned {
		return nil, nil, ErrNotAssigned
	}

	oldReviewerUser, err := s.repoUsers.GetByID(ctx, tx, oldUserID)
	if err != nil {
		return nil, nil, err
	}
	if oldReviewerUser == nil {
		return nil, nil, ErrUserNotFound
	}

	candidates, err := s.repoUsers.GetActiveCandidates(ctx, tx, oldReviewerUser.TeamName, currentReviewerIDs)
	if err != nil {
		return nil, nil, err
	}

	if len(candidates) == 0 {
		return nil, nil, ErrNoCandidate
	}

	newReviewer := selectRandomReviewers(candidates, 1)[0]

	if err := s.repoPR.ReplaceReviewer(ctx, tx, prID, oldUserID, newReviewer.ID); err != nil {
		return nil, nil, err
	}

	for i, r := range pr.Reviewers {
		if r.ID == oldUserID {
			pr.Reviewers[i] = newReviewer
			break
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	return pr, &newReviewer, nil
}

func (s *Service) GetUserReviews(ctx context.Context, userID string) ([]domain.PullRequestShort, error) {
	return s.repoPR.GetByReviewerID(ctx, s.db, userID)
}
