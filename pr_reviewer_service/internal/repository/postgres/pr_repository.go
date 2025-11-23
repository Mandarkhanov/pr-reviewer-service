package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/repository"
	"time"
)

type PRRepo struct{}

func NewPRRepo() *PRRepo {
	return &PRRepo{}
}

func (r *PRRepo) Exists(ctx context.Context, db repository.Querier, id string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pull_requests WHERE id = $1)"
	err := db.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (r *PRRepo) Create(ctx context.Context, db repository.Querier, pr domain.PullRequest) error {
	queryPR := `INSERT INTO pull_requests (id, name, author_id, status, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.ExecContext(ctx, queryPR, pr.ID, pr.Name, pr.AuthorID, pr.Status, time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert PR: %w", err)
	}

	if len(pr.Reviewers) > 0 {
		queryRev := `INSERT INTO pull_requests_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2)`
		for _, reviewer := range pr.Reviewers {
			_, err := db.ExecContext(ctx, queryRev, pr.ID, reviewer.ID)
			if err != nil {
				return fmt.Errorf("failed to insert reviewer %s: %w", reviewer.ID, err)
			}
		}
	}
	return nil
}

func (r *PRRepo) GetByID(ctx context.Context, db repository.Querier, id string) (*domain.PullRequest, error) {
	queryPR := `
		SELECT id, name, author_id, status, created_at, merged_at 
		FROM pull_requests 
		WHERE id = $1
	`
	var pr domain.PullRequest
	var createdAt sql.NullTime
	var mergedAt sql.NullTime

	err := db.QueryRowContext(ctx, queryPR, id).Scan(
		&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &createdAt, &mergedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get PR: %w", err)
	}
	pr.CreatedAt = createdAt.Time
	if mergedAt.Valid {
		pr.MergedAt = &mergedAt.Time
	} else {
		pr.MergedAt = nil
	}

	queryReviewers := `
		SELECT u.id, u.username, u.is_active, u.team_name
		FROM users u
		JOIN pull_requests_reviewers prr ON u.id = prr.reviewer_id
		WHERE prr.pull_request_id = $1
	`
	rows, err := db.QueryContext(ctx, queryReviewers, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviewers: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.IsActive, &u.TeamName); err != nil {
			return nil, err
		}
		pr.Reviewers = append(pr.Reviewers, u)
	}

	return &pr, nil
}

func (r *PRRepo) SetStatus(ctx context.Context, db repository.Querier, id string, status domain.PRStatus) error {
	var query string
	if status == domain.PRStatusMerged {
		query = "UPDATE pull_requests SET status = $1, merged_at = NOW() WHERE id = $2"
	} else {
		query = "UPDATE pull_requests SET status = $1 WHERE id = $2"
	}

	_, err := db.ExecContext(ctx, query, status, id)
	return err
}

func (r *PRRepo) ReplaceReviewer(ctx context.Context, db repository.Querier, prID, oldReviewerID, newReviewerID string) error {
	queryDel := "DELETE FROM pull_requests_reviewers WHERE pull_request_id = $1 AND reviewer_id = $2"
	_, err := db.ExecContext(ctx, queryDel, prID, oldReviewerID)
	if err != nil {
		return err
	}

	queryIns := "INSERT INTO pull_requests_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2)"
	_, err = db.ExecContext(ctx, queryIns, prID, newReviewerID)
	return err
}

func (r *PRRepo) GetByReviewerID(ctx context.Context, db repository.Querier, reviewerID string) ([]domain.PullRequestShort, error) {
	query := `
		SELECT pr.id, pr.name, pr.author_id, pr.status
		FROM pull_requests pr
		JOIN pull_requests_reviewers prr ON pr.id = prr.pull_request_id
		WHERE prr.reviewer_id = $1
	`

	rows, err := db.QueryContext(ctx, query, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user reviews: %w", err)
	}
	defer rows.Close()

	// Инициализируем как пустой слайс, чтобы JSON был [], если записей нет
	result := []domain.PullRequestShort{}

	for rows.Next() {
		var pr domain.PullRequestShort
		if err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		result = append(result, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
