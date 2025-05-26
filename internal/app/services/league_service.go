package services

import (
	"fmt"
	"sort"

	"github.com/muzaffertuna/football-league-sim/internal/app/models"
	"github.com/muzaffertuna/football-league-sim/internal/app/repositories"
)

type leagueService struct {
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

	// Bu döngüye gerek kalmadı çünkü initializeCurrentWeek ve currentWeek mantığı ile
	// zaten oynanmış bir haftayı tekrar oynamaya çalışmamalıyız.
	// for _, match := range matches {
	// 	if match.Played {
	// 		return fmt.Errorf("week %d has already been played", week)
	// 	}
	// }

	for i := range matches {
		match := &matches[i]
		// MatchService'in zaten oynanmış maçları atladığı varsayılır.
		// Ancak defensive programlama için burada da bir kontrol tutulabilir.
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

		// BURADA SADECE MAÇI SİMÜLE ETMEK ÇAĞRILIR.
		// Takım istatistikleri (MatchesPlayed, GoalsFor, GoalsAgainst, Points, Wins, Draws, Loses)
		// SimulateMatch fonksiyonu içinde güncellenecektir.
		if err := s.matchSvc.SimulateMatch(match, homeTeam, awayTeam); err != nil {
			return err
		}
		// Takımları güncelleme çağrılarına burada gerek yok, SimulateMatch içinde yapılıyor.
	}

	// Haftayı başarıyla oynadıktan sonra currentWeek'i bir artır
	s.currentWeek = week + 1
	return nil
}

// SimulateAllWeeks ligdeki kalan tüm haftaları sırayla simüle eder.
func (s *leagueService) SimulateAllWeeks() ([]models.Match, error) {
	var allSimulatedMatches []models.Match
	const totalWeeks = 6 // Ligi tamamlamak için gereken hafta sayısı (sabit varsayalım)

	// Ligin mevcut haftasından başla
	currentWeek, err := s.GetCurrentWeek()
	if err != nil {
		return nil, fmt.Errorf("failed to get current week: %w", err)
	}

	// Eğer lig zaten bitmişse, tüm maçları döndür
	if currentWeek > totalWeeks {
		allMatches, err := s.matchRepo.GetAllMatches()
		if err != nil {
			return nil, err
		}
		// Zaten bitmişse hata döndürmek yerine bilgilendirme mesajı verilebilir.
		return allMatches, fmt.Errorf("league has already completed. current week: %d", currentWeek)
	}

	for week := currentWeek; week <= totalWeeks; week++ {
		// PlayWeek metodunu çağırarak mevcut haftayı oynat
		err := s.PlayWeek(week)
		if err != nil {
			// Eğer haftanın oynanmasıyla ilgili bir hata olursa, dur ve hatayı döndür
			return nil, fmt.Errorf("failed to play week %d: %w", week, err)
		}

		// Oynanan haftanın maçlarını al ve genel listeye ekle
		playedMatches, err := s.matchRepo.GetMatchesByWeek(week)
		if err != nil {
			return nil, fmt.Errorf("failed to get matches for played week %d: %w", week, err)
		}
		allSimulatedMatches = append(allSimulatedMatches, playedMatches...)

		// İsteğe bağlı: Her hafta arasında kısa bir duraklama ekleyebilirsiniz.
		// time.Sleep(50 * time.Millisecond)
	}

	return allSimulatedMatches, nil
}

func (s *leagueService) GetLeagueTable() (*models.League, error) {
	// Takımları al
	teams, err := s.teamSvc.GetAllTeams()
	if err != nil {
		return nil, err
	}

	// Maçları al (opsiyonel, eğer League modelinde maçları da göstermek istiyorsan)
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
		Teams:       teams,
		Matches:     allMatches,
		CurrentWeek: s.currentWeek, // CurrentWeek'i buradan set ediyoruz
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

// calculatePoints fonksiyonu artık doğrudan puanları hesapladığı için eski şekilde eklenmiyor.
// PlayWeek içinde doğrudan galibiyete 3, beraberliğe 1 puan ekleniyor.
// Bu fonksiyonu başka bir yerde kullanmıyorsanız silebilirsiniz.
func calculatePoints(homeGoals, awayGoals int) int {
	if homeGoals > awayGoals {
		return 3 // Galibiyet
	} else if homeGoals == awayGoals {
		return 1 // Beraberlik
	} else {
		return 0 // Mağlubiyet
	}
}
