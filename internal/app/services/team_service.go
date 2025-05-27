package services

import (
	"github.com/muzaffertuna/football-league-sim/internal/app/models"
	"github.com/muzaffertuna/football-league-sim/internal/app/repositories"
)

type teamService struct {
	teamRepo repositories.TeamRepository
}

func NewTeamService(teamRepo repositories.TeamRepository) TeamService {
	return &teamService{teamRepo: teamRepo}
}

func (s *teamService) CreateTeam(name string, strength int) (*models.Team, error) {
	team := &models.Team{
		Name:          name,
		Strength:      strength,
		Points:        0,
		GoalsFor:      0,
		GoalsAgainst:  0,
		MatchesPlayed: 0,
	}
	if err := s.teamRepo.CreateTeam(team); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *teamService) GetTeamByID(id int) (*models.Team, error) {
	return s.teamRepo.GetTeamByID(id)
}

func (s *teamService) GetAllTeams() ([]models.Team, error) {
	return s.teamRepo.GetAllTeams()
}
