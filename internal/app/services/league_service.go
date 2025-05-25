package services

import (
	"fmt"
	"sort"

	"github.com/muzaffertuna/football-league-sim/internal/app/models"
	"github.com/muzaffertuna/football-league-sim/internal/app/repositories"
)

type leagueService struct {
	// leagueRepo repositories.LeagueRepository // Artık doğrudan League tablosu kullanmadığımız için kaldırıldı
	matchRepo   repositories.MatchRepository
	matchSvc    MatchService
	teamRepo    repositories.TeamRepository
	teamSvc     TeamService
	currentWeek int // Ligin güncel haftasını tutacak alan
}

// NewLeagueService Constructor'ı güncellendi ve başlangıç haftası hesaplaması eklendi
func NewLeagueService(matchRepo repositories.MatchRepository, matchSvc MatchService, teamRepo repositories.TeamRepository, teamSvc TeamService) (LeagueService, error) {
	ls := &leagueService{
		matchRepo: matchRepo,
		matchSvc:  matchSvc,
		teamRepo:  teamRepo,
		teamSvc:   teamSvc,
	}

	// Uygulama başladığında CurrentWeek'i hesapla
	err := ls.initializeCurrentWeek()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize current week: %w", err)
	}

	return ls, nil
}

// initializeCurrentWeek başlangıçta mevcut haftayı hesaplar ve s.currentWeek'e atar
func (s *leagueService) initializeCurrentWeek() error {
	maxPlayedWeek, err := s.matchRepo.GetMaxWeekPlayed()
	if err != nil {
		return err
	}

	// Eğer hiç maç oynanmamışsa, CurrentWeek 1'dir.
	if maxPlayedWeek == 0 {
		s.currentWeek = 1
		return nil
	}

	// maxPlayedWeek'teki tüm maçların oynanıp oynanmadığını kontrol et
	matchesInMaxWeek, err := s.matchRepo.GetMatchesByWeek(maxPlayedWeek)
	if err != nil {
		return err
	}

	allPlayedInMaxWeek := true
	for _, match := range matchesInMaxWeek {
		if !match.Played {
			allPlayedInMaxWeek = false
			break
		}
	}

	if allPlayedInMaxWeek {
		// Eğer en yüksek haftadaki tüm maçlar oynandıysa, bir sonraki haftaya geç.
		s.currentWeek = maxPlayedWeek + 1
	} else {
		// Aksi takdirde, hala en yüksek haftadayız.
		s.currentWeek = maxPlayedWeek
	}
	return nil
}

// GetCurrentWeek ligin güncel haftasını döndürür
func (s *leagueService) GetCurrentWeek() (int, error) {
	return s.currentWeek, nil
}

func (s *leagueService) PlayWeek(week int) error {
	// Doğrudan servis içindeki currentWeek'i kullan
	if week != s.currentWeek {
		return fmt.Errorf("it's not week %d, current week is %d", week, s.currentWeek)
	}

	matches, err := s.matchRepo.GetMatchesByWeek(week)
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		return fmt.Errorf("no matches found for week %d", week)
	}

	for _, match := range matches {
		if match.Played {
			// Bu kontrol aslında yukarıdaki 'week != s.currentWeek' ile kısmen ele alınıyor
			// Ancak emin olmak için bırakılabilir veya kaldırılabilir.
			return fmt.Errorf("week %d has already been played", week)
		}
	}

	for i := range matches {
		match := &matches[i]
		// Maç zaten oynandıysa döngüyü atla
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

	// Haftayı başarıyla oynadıktan sonra currentWeek'i bir artır
	s.currentWeek = week + 1
	return nil
}

func (s *leagueService) GetLeagueTable() (*models.League, error) {
	// Takımları al
	teams, err := s.teamSvc.GetAllTeams()
	if err != nil {
		return nil, err
	}

	// Maçları al (opsiyonel, eğer League modelinde maçları da göstermek istiyorsan)
	// CurrentWeek'i burada doğrudan league.CurrentWeek'e atayabiliriz
	allMatches, err := s.matchRepo.GetAllMatches()
	if err != nil {
		return nil, err
	}

	sort.Slice(teams, func(i, j int) bool {
		if teams[i].Points != teams[j].Points {
			return teams[i].Points > teams[j].Points
		}
		if teams[i].GoalDifference() != teams[j].GoalDifference() {
			return teams[i].GoalDifference() > teams[j].GoalDifference()
		}
		return teams[i].GoalsFor > teams[j].GoalsFor
	})

	league := &models.League{
		Teams:   teams,
		Matches: allMatches,
		// CurrentWeek'i buradan set ediyoruz
		CurrentWeek: s.currentWeek,
	}

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
		team.Wins = 0
		team.Draws = 0
		team.Loses = 0
		if err := s.teamRepo.UpdateTeam(&team); err != nil {
			return err
		}
	}

	if err := s.matchRepo.DeleteAllMatches(); err != nil {
		return err
	}

	if err := s.generateMatches(teams); err != nil {
		return err
	}

	// Ligi sıfırladıktan sonra currentWeek'i 1'e geri getir
	s.currentWeek = 1
	return nil
}

func (s *leagueService) generateMatches(teams []models.Team) error {
	if len(teams) != 4 {
		return fmt.Errorf("expected 4 teams, got %d", len(teams))
	}

	matches := []struct {
		homeTeamID int
		awayTeamID int
		week       int
	}{
		{teams[0].ID, teams[1].ID, 1},
		{teams[2].ID, teams[3].ID, 1},
		{teams[0].ID, teams[2].ID, 2},
		{teams[1].ID, teams[3].ID, 2},
		{teams[0].ID, teams[3].ID, 3},
		{teams[1].ID, teams[2].ID, 3},
		{teams[1].ID, teams[0].ID, 4},
		{teams[3].ID, teams[2].ID, 4},
		{teams[2].ID, teams[0].ID, 5},
		{teams[3].ID, teams[1].ID, 5},
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
