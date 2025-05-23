package repositories

import "github.com/muzaffertuna/football-league-sim/internal/app/models"

type leagueRepository struct {
	teamRepo  TeamRepository
	matchRepo MatchRepository
	league    *models.League
}

func NewLeagueRepository(teamRepo TeamRepository, matchRepo MatchRepository) LeagueRepository {
	return &leagueRepository{
		teamRepo:  teamRepo,
		matchRepo: matchRepo,
		league:    &models.League{CurrentWeek: 1},
	}
}

func (r *leagueRepository) GetLeague() (*models.League, error) {
	teams, err := r.teamRepo.GetAllTeams()
	if err != nil {
		return nil, err
	}
	matches, err := r.matchRepo.GetAllMatches()
	if err != nil {
		return nil, err
	}

	r.league.Teams = teams
	r.league.Matches = matches
	return r.league, nil
}

func (r *leagueRepository) SaveLeague(league *models.League) error {
	r.league = league
	return nil
}
