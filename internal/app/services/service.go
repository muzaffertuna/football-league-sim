package services

import "github.com/muzaffertuna/football-league-sim/internal/app/models"

type TeamService interface {
	CreateTeam(name string, strength int) (*models.Team, error)
	GetTeamByID(id int) (*models.Team, error)
	GetAllTeams() ([]models.Team, error)
}

type MatchService interface {
	CreateMatch(homeTeamID, awayTeamID, week int) (*models.Match, error)
	SimulateMatch(match *models.Match, homeTeam, awayTeam *models.Team) error
	GetMatchesByWeek(week int) ([]models.Match, error)
}

type LeagueService interface {
	PlayWeek(week int) error
	GetLeagueTable() (*models.League, error)
	ResetLeague() error
	GetMatchesByWeek(week int) ([]models.Match, error) // Yeni metod
	GetTeamByID(id int) (*models.Team, error)          // Yeni metod
	GetCurrentWeek() (int, error)                      // Yeni metot
	SimulateAllWeeks() ([]models.Match, error)         // YENİ EKLENEN METOD
}
