package repositories

import "github.com/muzaffertuna/football-league-sim/internal/app/models"

type TeamRepository interface {
	CreateTeam(team *models.Team) error
	GetTeamByID(id int) (*models.Team, error)
	GetAllTeams() ([]models.Team, error)
	UpdateTeam(team *models.Team) error
}

type MatchRepository interface {
	CreateMatch(match *models.Match) error
	GetMatchByID(id int) (*models.Match, error)
	GetMatchesByWeek(week int) ([]models.Match, error)
	UpdateMatch(match *models.Match) error
	DeleteAllMatches() error
	GetAllMatches() ([]models.Match, error)
	GetPlayedMatches() ([]models.Match, error)
	GetMaxWeekPlayed() (int, error)
}

type LeagueRepository interface {
	GetLeague() (*models.League, error)
	SaveLeague(league *models.League) error
}
