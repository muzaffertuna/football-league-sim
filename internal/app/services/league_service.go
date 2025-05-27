package services

import (
	"fmt"
	"sort"
	"sync"
	"time"

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

func (s *leagueService) SimulateAllWeeks() ([]models.Match, error) {
	var allSimulatedMatches []models.Match
	const totalWeeks = 6

	currentWeek, err := s.GetCurrentWeek()
	if err != nil {
		return nil, fmt.Errorf("failed to get current week: %w", err)
	}

	if currentWeek > totalWeeks {
		allMatches, err := s.matchRepo.GetAllMatches()
		if err != nil {
			return nil, err
		}
		return allMatches, fmt.Errorf("league has already completed. current week: %d", currentWeek)
	}

	for week := currentWeek; week <= totalWeeks; week++ {
		err := s.PlayWeek(week)
		if err != nil {
			return nil, fmt.Errorf("failed to play week %d: %w", week, err)
		}

		playedMatches, err := s.matchRepo.GetMatchesByWeek(week)
		if err != nil {
			return nil, fmt.Errorf("failed to get matches for played week %d: %w", week, err)
		}
		allSimulatedMatches = append(allSimulatedMatches, playedMatches...)

	}

	return allSimulatedMatches, nil
}

func (s *leagueService) GetLeagueTable() (*models.League, error) {
	fmt.Println("GetLeagueTable: Starting...")
	teams, err := s.teamSvc.GetAllTeams()
	if err != nil {
		fmt.Printf("GetLeagueTable: Error getting all teams: %v\n", err)
		return nil, err
	}
	fmt.Println("GetLeagueTable: Teams retrieved.")

	allMatches, err := s.matchRepo.GetAllMatches()
	if err != nil {
		fmt.Printf("GetLeagueTable: Error getting all matches: %v\n", err)
		return nil, err
	}
	fmt.Println("GetLeagueTable: Matches retrieved.")

	sort.Slice(teams, func(i, j int) bool {
		if teams[i].Points != teams[j].Points {
			return teams[i].Points > teams[j].Points
		}
		if teams[i].GoalDifference() != teams[j].GoalDifference() {
			return teams[i].GoalDifference() > teams[j].GoalDifference()
		}
		return teams[i].GoalsFor > teams[j].GoalsFor
	})
	fmt.Println("GetLeagueTable: Teams sorted.")

	const numSimulationsForTable = 1000
	fmt.Printf("GetLeagueTable: Calling PredictOutcomes with %d simulations.\n", numSimulationsForTable)
	predictionResult, err := s.PredictOutcomes(numSimulationsForTable)
	if err != nil {
		fmt.Printf("GetLeagueTable: Error predicting outcomes: %v\n", err)
		return nil, fmt.Errorf("failed to predict outcomes for league table: %w", err)
	}
	fmt.Println("GetLeagueTable: PredictOutcomes completed successfully.")

	league := &models.League{
		Teams:                   teams,
		Matches:                 allMatches,
		CurrentWeek:             s.currentWeek,
		ChampionshipPredictions: predictionResult.ChampionshipPredictions,
	}
	fmt.Println("GetLeagueTable: League table constructed. Returning.")
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

//--------------------------------------------------------------------

// ... (mevcut LeagueService arayüzü ve leagueService struct tanımları) ...

func (s *leagueService) PredictOutcomes(numSimulations int) (models.PredictionResult, error) {
	fmt.Printf("PredictOutcomes: Starting %d simulations...\n", numSimulations)
	startTime := time.Now() // Debug için zaman tutucu

	initialTeams, err := s.teamRepo.GetAllTeams()
	if err != nil {
		fmt.Printf("PredictOutcomes: Failed to get initial teams: %v\n", err)
		return models.PredictionResult{}, fmt.Errorf("failed to get initial teams for prediction: %w", err)
	}
	initialMatches, err := s.matchRepo.GetAllMatches()
	if err != nil {
		fmt.Printf("PredictOutcomes: Failed to get initial matches: %v\n", err)
		return models.PredictionResult{}, fmt.Errorf("failed to get initial matches for prediction: %w", err)
	}
	initialCurrentWeek := s.currentWeek
	fmt.Printf("PredictOutcomes: Initial league state captured (current week: %d).\n", initialCurrentWeek)

	championshipCounts := &sync.Map{}

	var wg sync.WaitGroup
	resultsChan := make(chan int, numSimulations)
	errorChan := make(chan error, numSimulations) // Hataları yakalamak için kanal

	for i := 0; i < numSimulations; i++ {
		wg.Add(1)
		go func(simIndex int) { // simIndex parametre olarak eklendi
			defer wg.Done()
			fmt.Printf("PredictOutcomes: Goroutine %d started (Sim #%d).\n", simIndex, simIndex+1)

			tempTeamRepo := repositories.NewInMemoryTeamRepository()
			tempMatchRepo := repositories.NewInMemoryMatchRepository()

			// Takım verilerini kopyala
			for _, team := range initialTeams {
				copiedTeam := models.Team{ // Deep copy
					ID:            team.ID,
					Name:          team.Name,
					Strength:      team.Strength,
					Points:        team.Points,
					GoalsFor:      team.GoalsFor,
					GoalsAgainst:  team.GoalsAgainst,
					MatchesPlayed: team.MatchesPlayed,
					Wins:          team.Wins,
					Draws:         team.Draws,
					Loses:         team.Loses,
				}
				if err := tempTeamRepo.CreateTeam(&copiedTeam); err != nil {
					errMsg := fmt.Sprintf("PredictOutcomes: Sim %d failed to copy team %d: %v", simIndex, team.ID, err)
					fmt.Println(errMsg)
					errorChan <- fmt.Errorf(errMsg)
					return
				}
			}
			fmt.Printf("PredictOutcomes: Sim %d teams copied.\n", simIndex)

			// Maç verilerini kopyala
			for _, match := range initialMatches {
				copiedMatch := models.Match{ // Deep copy
					ID:         match.ID,
					HomeTeamID: match.HomeTeamID,
					AwayTeamID: match.AwayTeamID,
					Week:       match.Week,
					HomeGoals:  match.HomeGoals,
					AwayGoals:  match.AwayGoals,
					Played:     match.Played,
				}
				if err := tempMatchRepo.CreateMatch(&copiedMatch); err != nil {
					errMsg := fmt.Sprintf("PredictOutcomes: Sim %d failed to copy match %d: %v", simIndex, match.ID, err)
					fmt.Println(errMsg)
					errorChan <- fmt.Errorf(errMsg)
					return
				}
			}
			fmt.Printf("PredictOutcomes: Sim %d matches copied.\n", simIndex)

			tempTeamSvc := NewTeamService(tempTeamRepo)
			tempMatchSvc := NewMatchService(tempMatchRepo, tempTeamRepo)

			tempLeagueSvc := &leagueService{
				matchRepo:   tempMatchRepo,
				matchSvc:    tempMatchSvc,
				teamRepo:    tempTeamRepo,
				teamSvc:     tempTeamSvc,
				currentWeek: initialCurrentWeek,
			}

			fmt.Printf("PredictOutcomes: Sim %d calling SimulateAllWeeks...\n", simIndex)
			_, simErr := tempLeagueSvc.SimulateAllWeeks()
			if simErr != nil && simErr.Error() != fmt.Sprintf("league has already completed. current week: %d", tempLeagueSvc.currentWeek) {
				errMsg := fmt.Sprintf("PredictOutcomes: Sim %d simulation failed: %v", simIndex, simErr)
				fmt.Println(errMsg)
				errorChan <- fmt.Errorf(errMsg)
				return
			}
			fmt.Printf("PredictOutcomes: Sim %d SimulateAllWeeks completed.\n", simIndex)

			// Şampiyonu belirlemek için tempLeagueSvc'nin GetLeagueTable'ını çağırmak yerine,
			// doğrudan tempTeamRepo'dan takımları alıp sıralayalım.
			finalTeams, simErr := tempTeamRepo.GetAllTeams() // Doğrudan repo'dan al
			if simErr != nil {
				errMsg := fmt.Sprintf("PredictOutcomes: Sim %d failed to get final teams from tempRepo: %v", simIndex, simErr)
				fmt.Println(errMsg)
				errorChan <- fmt.Errorf(errMsg)
				return
			}

			// Takımları sırala (şimdiki GetLeagueTable mantığının aynısı)
			sort.Slice(finalTeams, func(i, j int) bool {
				if finalTeams[i].Points != finalTeams[j].Points {
					return finalTeams[i].Points > finalTeams[j].Points
				}
				if finalTeams[i].GoalDifference() != finalTeams[j].GoalDifference() {
					return finalTeams[i].GoalDifference() > finalTeams[j].GoalDifference()
				}
				return finalTeams[i].GoalsFor > finalTeams[j].GoalsFor
			})

			if len(finalTeams) > 0 {
				resultsChan <- finalTeams[0].ID
				fmt.Printf("PredictOutcomes: Sim %d identified winner Team ID: %d. Completed.\n", simIndex, finalTeams[0].ID)
			} else {
				fmt.Printf("PredictOutcomes: Sim %d found no teams in final state. This should not happen if league has teams.\n", simIndex)
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("PredictOutcomes: All %d goroutines finished. Closing channels.\n", numSimulations)
	close(resultsChan)
	close(errorChan)

	select {
	case err := <-errorChan:
		if err != nil {
			fmt.Printf("PredictOutcomes: Main error channel received critical error from goroutine: %v\n", err)
			return models.PredictionResult{}, err // Kritik hata varsa hemen dön
		}
	default:
		// Hata yok
	}
	fmt.Println("PredictOutcomes: Processing simulation results from resultsChan.")

	// CHAMPIONSHIP COUNTS'A EKLEME KISMINI KONTROL ET
	for teamID := range resultsChan {
		actualCount, ok := championshipCounts.Load(teamID)
		if !ok {
			// Eğer yoksa, 0 ile başlat
			actualCount = 0
		}
		championshipCounts.Store(teamID, actualCount.(int)+1)
		fmt.Printf("PredictOutcomes: Counted team %d. Current count: %d\n", teamID, actualCount.(int)+1)
	}
	fmt.Println("PredictOutcomes: Finished counting championship wins.")

	var championshipPredictions []models.Prediction
	fmt.Println("PredictOutcomes: Populating championship predictions slice.")
	championshipCounts.Range(func(key, value interface{}) bool {
		teamID := key.(int)
		count := value.(int)

		team, err := s.teamRepo.GetTeamByID(teamID) // Orijinal repo'dan takım bilgisini al
		if err != nil {
			// **BURADA HATA YAKALANDIĞINDA DÖNGÜYÜ DURDURAN return false İFADESİ SORUN OLABİLİR.**
			// Hatayı logla ama döngüyü devam ettirerek diğer tahminleri göstermeye çalışalım.
			// Eğer teamRepo'da takım yoksa, bu bir problemdir ve loglanmalı.
			fmt.Printf("PredictOutcomes: ERROR: Failed to get team by ID %d for prediction result: %v. This team's prediction will be skipped.\n", teamID, err)
			return true // Hata olsa bile diğer takımlar için devam et
		}

		fmt.Printf("PredictOutcomes: Adding prediction for Team %s (ID: %d), wins: %d\n", team.Name, teamID, count)
		championshipPredictions = append(championshipPredictions, models.Prediction{
			TeamID:                 teamID,
			TeamName:               team.Name,
			ChampionshipLikelihood: float64(count) / float64(numSimulations) * 100,
		})
		return true
	})

	fmt.Printf("PredictOutcomes: Number of predictions gathered: %d\n", len(championshipPredictions))

	if len(championshipPredictions) > 0 {
		sort.Slice(championshipPredictions, func(i, j int) bool {
			return championshipPredictions[i].ChampionshipLikelihood > championshipPredictions[j].ChampionshipLikelihood
		})
		fmt.Println("PredictOutcomes: Championship predictions sorted.")
	} else {
		fmt.Println("PredictOutcomes: No championship predictions to sort (slice is empty).")
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("PredictOutcomes: Completed %d simulations in %s. Returning prediction result.\n", numSimulations, elapsedTime)
	return models.PredictionResult{
		ChampionshipPredictions: championshipPredictions,
	}, nil
}
