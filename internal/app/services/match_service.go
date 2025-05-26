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
		return nil // Zaten oynanmışsa hiçbir şey yapma
	}

	rand.Seed(time.Now().UnixNano())
	homeAdvantage := 10 // Ev sahibi avantajı
	// Gol şansları
	homeChance := float64(homeTeam.Strength+homeAdvantage) / float64(homeTeam.Strength+awayTeam.Strength+homeAdvantage)
	// awayChance := float64(awayTeam.Strength) / float64(homeTeam.Strength+awayTeam.Strength+homeAdvantage) // Away için ayrı bir şansa gerek yok, 1-homeChance yeterli

	homeGoals := 0
	awayGoals := 0

	// Her takım için 5'er şut, her şut için gol olup olmadığını kontrol et
	for i := 0; i < 5; i++ {
		if rand.Float64() < homeChance {
			homeGoals++
		}
		// Deplasman takımının gol şansı: 1 - ev sahibi gol şansı
		// Bu şekilde goller birbirini dengelemez ve daha gerçekçi olur.
		if rand.Float64() < (1 - homeChance) { // Bu satırda daha önce sorun yoktu
			awayGoals++
		}
	}

	match.HomeGoals = homeGoals
	match.AwayGoals = awayGoals
	match.Played = true
	// Maç sonucunu repository'e kaydet
	if err := s.matchRepo.UpdateMatch(match); err != nil {
		return err
	}

	// TAKIM İSTATİSTİKLERİNİ GÜNCELLEME SADECE BURADA YAPILACAK!
	homeTeam.MatchesPlayed++
	awayTeam.MatchesPlayed++

	homeTeam.GoalsFor += homeGoals
	homeTeam.GoalsAgainst += awayGoals
	awayTeam.GoalsFor += awayGoals
	awayTeam.GoalsAgainst += homeGoals

	if homeGoals > awayGoals {
		homeTeam.Wins++
		awayTeam.Loses++
		homeTeam.Points += 3 // Galibiyet 3 puan
	} else if homeGoals < awayGoals {
		homeTeam.Loses++
		awayTeam.Wins++
		awayTeam.Points += 3 // Galibiyet 3 puan
	} else { // Beraberlik
		homeTeam.Draws++
		awayTeam.Draws++
		homeTeam.Points += 1 // Beraberlik 1 puan
		awayTeam.Points += 1 // Beraberlik 1 puan
	}

	// Takımların güncellenmiş hallerini repository'e kaydet
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
