package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/repository"
	"strings"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) Upsert(ctx context.Context, db repository.Querier, users []domain.User) error {
	if len(users) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(users))
	valueArgs := make([]any, 0, len(users)*4)

	for i, u := range users {
		n := i * 4
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", n+1, n+2, n+3, n+4))
		valueArgs = append(valueArgs, u.ID, u.Username, u.IsActive, u.TeamName)
	}

	query := fmt.Sprintf(`
		INSERT INTO users (id, username, is_active, team_name)
		VALUES %s
		ON CONFLICT (id) DO UPDATE
		SET username = EXCLUDED.username,
			is_active = EXCLUDED.is_active,
			team_name = EXCLUDED.team_name
	`, strings.Join(valueStrings, ","))

	_, err := db.ExecContext(ctx, query, valueArgs...)
	return err
}

func (r *UserRepo) SetIsActive(ctx context.Context, db repository.Querier, userID string, isActive bool) (*domain.User, error) {
	query := `
		UPDATE users
		SET is_active = $2
		WHERE id = $1
		RETURNING id, username, is_active, team_name
	`

	var u domain.User
	err := db.QueryRowContext(ctx, query, userID, isActive).Scan(&u.ID, &u.Username, &u.IsActive, &u.TeamName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update user active status: %w", err)
	}

	return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, db repository.Querier, userID string) (*domain.User, error) {
	query := "SELECT id, username, is_active, team_name FROM users WHERE id = $1"
	var u domain.User
	err := db.QueryRowContext(ctx, query, userID).Scan(&u.ID, &u.Username, &u.IsActive, &u.TeamName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) GetActiveCandidates(ctx context.Context, db repository.Querier, teamName string, excludeUserIDs []string) ([]domain.User, error) {
	query := "SELECT id, username, is_active, team_name FROM users WHERE team_name = $1 AND is_active = true"
	args := []any{teamName}

	if len(excludeUserIDs) > 0 {
		placeholders := make([]string, len(excludeUserIDs))
		for i, id := range excludeUserIDs {
			// +2, так как $1 занят teamName
			placeholders[i] = fmt.Sprintf("$%d", i+2)
			args = append(args, id)
		}
		query += fmt.Sprintf(" AND id NOT IN (%s)", strings.Join(placeholders, ", "))
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.IsActive, &u.TeamName); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
