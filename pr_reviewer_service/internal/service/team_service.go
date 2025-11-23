package service

import (
	"context"
	"pr-reviewer/internal/domain"
)

func (s *Service) CreateTeam(ctx context.Context, team domain.Team) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.repoTeams.Create(ctx, tx, team); err != nil {
		return err
	}

	for i := range team.Members {
		team.Members[i].TeamName = team.Name
	}

	if err := s.repoUsers.Upsert(ctx, tx, team.Members); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Service) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	team, err := s.repoTeams.GetByName(ctx, s.db, name)
	if err != nil {
		return nil, err
	}

	if team == nil {
		return nil, ErrTeamNotFound
	}

	return team, nil
}
