package services

import (
	"fmt"
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

	if len(matches) == 0 {
		return fmt.Errorf("no matches found for week %d", week)
	}

	for _, match := range matches {
		if match.Played {
			return fmt.Errorf("week %d has already been played", week)
		}
	}

	for i := range matches {
		match := &matches[i]
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

		if err := s.matchSvc.SimulateMatch(match, homeTeam, awayTeam); err != nil {
			return err
		}

		// Maç sonucuna göre takım istatistiklerini güncelle
		if match.HomeGoals > match.AwayGoals {
			homeTeam.Wins++
			awayTeam.Loses++
		} else if match.HomeGoals < match.AwayGoals {
			homeTeam.Loses++
			awayTeam.Wins++
		} else {
			homeTeam.Draws++
			awayTeam.Draws++
		}

		// Takımların güncellenmiş hallerini kaydet
		if err := s.teamRepo.UpdateTeam(homeTeam); err != nil {
			return err
		}
		if err := s.teamRepo.UpdateTeam(awayTeam); err != nil {
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
		team.Wins = 0  // Yeni eklendi
		team.Draws = 0 // Yeni eklendi
		team.Loses = 0 // Yeni eklendi
		if err := s.teamRepo.UpdateTeam(&team); err != nil {
			return err
		}
	}

	if err := s.matchRepo.DeleteAllMatches(); err != nil {
		return err
	}

	league, err := s.leagueRepo.GetLeague()
	if err != nil {
		return err
	}
	league.CurrentWeek = 1
	league.Matches = nil

	if err := s.generateMatches(teams); err != nil {
		return err
	}

	return s.leagueRepo.SaveLeague(league)
}

func (s *leagueService) generateMatches(teams []models.Team) error {
	if len(teams) != 4 {
		return fmt.Errorf("expected 4 teams, got %d", len(teams))
	}

	// 4 takım için round-robin fikstürü: 6 maç (ilk yarı) + 6 maç (rövanş) = 12 maç
	// Her hafta 2 maç, toplam 6 hafta
	matches := []struct {
		homeTeamID int
		awayTeamID int
		week       int
	}{
		// 1. Hafta
		{teams[0].ID, teams[1].ID, 1},
		{teams[2].ID, teams[3].ID, 1},
		// 2. Hafta
		{teams[0].ID, teams[2].ID, 2},
		{teams[1].ID, teams[3].ID, 2},
		// 3. Hafta
		{teams[0].ID, teams[3].ID, 3},
		{teams[1].ID, teams[2].ID, 3},
		// 4. Hafta (rövanşlar)
		{teams[1].ID, teams[0].ID, 4},
		{teams[3].ID, teams[2].ID, 4},
		// 5. Hafta
		{teams[2].ID, teams[0].ID, 5},
		{teams[3].ID, teams[1].ID, 5},
		// 6. Hafta
		{teams[3].ID, teams[0].ID, 6},
		{teams[2].ID, teams[1].ID, 6},
	}

	for _, m := range matches {
		match := &models.Match{
			HomeTeamID: m.homeTeamID,
			AwayTeamID: m.awayTeamID,
			Week:       m.week,
			Played:     false,
		}
		if err := s.matchRepo.CreateMatch(match); err != nil {
			return err
		}
	}

	return nil
}

func (s *leagueService) GetMatchesByWeek(week int) ([]models.Match, error) {
	return s.matchRepo.GetMatchesByWeek(week)
}

func (s *leagueService) GetTeamByID(id int) (*models.Team, error) {
	return s.teamSvc.GetTeamByID(id)
}
