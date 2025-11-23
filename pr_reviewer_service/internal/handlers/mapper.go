package handlers

import "pr-reviewer/internal/domain"

func toDomainTeam(req createTeamRequest) domain.Team {
	members := make([]domain.User, len(req.Members))
	for i, m := range req.Members {
		members[i] = domain.User{
			ID:       m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}

	return domain.Team{
		Name:    req.TeamName,
		Members: members,
	}
}
