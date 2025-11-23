package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/repository"
)

type TeamRepo struct{}

func NewTeamRepo() *TeamRepo {
	return &TeamRepo{}
}

func (r *TeamRepo) Create(ctx context.Context, db repository.Querier, team domain.Team) error {
	query := "INSERT INTO teams (name) VALUES ($1)"

	_, err := db.ExecContext(ctx, query, team.Name)
	if err != nil {
		return fmt.Errorf("failed to insert team: %w", err)
	}
	return nil
}

func (r *TeamRepo) GetByName(ctx context.Context, db repository.Querier, name string) (*domain.Team, error) {
	queryTeam := "SELECT name FROM teams WHERE name = $1"

	var team domain.Team
	err := db.QueryRowContext(ctx, queryTeam, name).Scan(&team.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	queryMembers := `
		SELECT id, username, is_active, team_name
		FROM users
		WHERE team_name = $1
	`
	rows, err := db.QueryContext(ctx, queryMembers, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.IsActive, &u.TeamName); err != nil {
			return nil, err
		}
		team.Members = append(team.Members, u)
	}

	return &team, nil
}
