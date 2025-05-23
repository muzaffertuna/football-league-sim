package services

import (
	"sort"

	"github.com/muzaffertuna/football-league-sim/internal/app/models"
	"github.com/muzaffertuna/football-league-sim/internal/app/repositories"
)

type leagueService struct {
	leagueRepo repositories.LeagueRepository
	matchRepo  repositories.MatchRepository
	matchSvc   MatchService
	teamRepo   repositories.TeamRepository
	teamSvc    TeamService
}

func NewLeagueService(leagueRepo repositories.LeagueRepository, matchRepo repositories.MatchRepository, matchSvc MatchService, teamRepo repositories.TeamRepository, teamSvc TeamService) LeagueService {
	return &leagueService{leagueRepo: leagueRepo, matchRepo: matchRepo, matchSvc: matchSvc, teamRepo: teamRepo, teamSvc: teamSvc}
}

func (s *leagueService) PlayWeek(week int) error {
	matches, err := s.matchRepo.GetMatchesByWeek(week)
	if err != nil {
		return err
	}

	for _, match := range matches {
		if match.Played {
			continue
		}

		homeTeam, err := s.teamSvc.GetTeamByID(match.HomeTeamID)
		if err != nil {
			return err
		}
		awayTeam, err := s.teamSvc.GetTeamByID(match.AwayTeamID)
		if err != nil {
			return err
		}

		if err := s.matchSvc.SimulateMatch(&match, homeTeam, awayTeam); err != nil {
			return err
		}
	}

	league, err := s.leagueRepo.GetLeague()
	if err != nil {
		return err
	}
	league.CurrentWeek = week + 1
	return s.leagueRepo.SaveLeague(league)
}

func (s *leagueService) GetLeagueTable() (*models.League, error) {
	league, err := s.leagueRepo.GetLeague()
	if err != nil {
		return nil, err
	}

	sort.Slice(league.Teams, func(i, j int) bool {
		if league.Teams[i].Points != league.Teams[j].Points {
			return league.Teams[i].Points > league.Teams[j].Points
		}
		if league.Teams[i].GoalDifference() != league.Teams[j].GoalDifference() {
			return league.Teams[i].GoalDifference() > league.Teams[j].GoalDifference()
		}
		return league.Teams[i].GoalsFor > league.Teams[j].GoalsFor
	})

	return league, nil
}

func (s *leagueService) ResetLeague() error {
	teams, err := s.teamSvc.GetAllTeams()
	if err != nil {
		return err
	}

	for _, team := range teams {
		team.Points = 0
		team.GoalsFor = 0
		team.GoalsAgainst = 0
		team.MatchesPlayed = 0
		if err := s.teamRepo.UpdateTeam(&team); err != nil {
			return err
		}
	}

	league, err := s.leagueRepo.GetLeague()
	if err != nil {
		return err
	}
	league.CurrentWeek = 1
	league.Matches = nil
	return s.leagueRepo.SaveLeague(league)
}
