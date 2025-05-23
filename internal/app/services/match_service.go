package services

import (
	"math/rand"
	"time"

	"github.com/muzaffertuna/football-league-sim/internal/app/models"
	"github.com/muzaffertuna/football-league-sim/internal/app/repositories"
)

type matchService struct {
	matchRepo repositories.MatchRepository
	teamRepo  repositories.TeamRepository
}

func NewMatchService(matchRepo repositories.MatchRepository, teamRepo repositories.TeamRepository) MatchService {
	return &matchService{matchRepo: matchRepo, teamRepo: teamRepo}
}

func (s *matchService) CreateMatch(homeTeamID, awayTeamID, week int) (*models.Match, error) {
	match := &models.Match{
		HomeTeamID: homeTeamID,
		AwayTeamID: awayTeamID,
		Week:       week,
		Played:     false,
	}
	if err := s.matchRepo.CreateMatch(match); err != nil {
		return nil, err
	}
	return match, nil
}

func (s *matchService) SimulateMatch(match *models.Match, homeTeam, awayTeam *models.Team) error {
	if match.Played {
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	homeAdvantage := 10
	homeChance := float64(homeTeam.Strength+homeAdvantage) / float64(homeTeam.Strength+awayTeam.Strength+homeAdvantage)
	homeGoals := 0
	awayGoals := 0

	for i := 0; i < 5; i++ {
		if rand.Float64() < homeChance {
			homeGoals++
		}
		if rand.Float64() < (1 - homeChance) {
			awayGoals++
		}
	}

	match.HomeGoals = homeGoals
	match.AwayGoals = awayGoals
	match.Played = true
	if err := s.matchRepo.UpdateMatch(match); err != nil {
		return err
	}

	homeTeam.MatchesPlayed++
	awayTeam.MatchesPlayed++
	homeTeam.GoalsFor += homeGoals
	homeTeam.GoalsAgainst += awayGoals
	awayTeam.GoalsFor += awayGoals
	awayTeam.GoalsAgainst += homeGoals

	if homeGoals > awayGoals {
		homeTeam.Points += 3
	} else if homeGoals < awayGoals {
		awayTeam.Points += 3
	} else {
		homeTeam.Points++
		awayTeam.Points++
	}

	if err := s.teamRepo.UpdateTeam(homeTeam); err != nil {
		return err
	}
	if err := s.teamRepo.UpdateTeam(awayTeam); err != nil {
		return err
	}

	return nil
}

func (s *matchService) GetMatchesByWeek(week int) ([]models.Match, error) {
	return s.matchRepo.GetMatchesByWeek(week)
}
