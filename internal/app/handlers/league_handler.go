package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/muzaffertuna/football-league-sim/internal/app/services"
	"github.com/muzaffertuna/football-league-sim/internal/pkg/logger" // <--- Düzeltildi: Yeni logger paketi

	// _ "github.com/muzaffertuna/football-league-sim/internal/platform" // <--- Kaldırıldı

	_ "github.com/muzaffertuna/football-league-sim/docs" // Swagger dokümantasyonu için
)

// @title Football League Simulation API
// @version 1.0
// @description Premier Lig simülasyonu için REST API
// @host localhost:8080
// @BasePath /
type LeagueHandler struct {
	leagueSvc services.LeagueService
	logger    *logger.Logger // <--- Düzeltildi: *platform.Logger yerine *logger.Logger
}

func NewLeagueHandler(leagueSvc services.LeagueService, logger *logger.Logger) *LeagueHandler { // <--- Düzeltildi
	return &LeagueHandler{leagueSvc: leagueSvc, logger: logger}
}

// @Summary Lig tablosunu getirir
// @Description Mevcut lig tablosunu puan sırasına göre döndürür
// @Tags league
// @Produce json
// @Success 200 {object} models.League
// @Failure 500 {string} string "Internal server error"
// @Router /league-table [get]
func (h *LeagueHandler) GetLeagueTable(w http.ResponseWriter, r *http.Request) {
	league, err := h.leagueSvc.GetLeagueTable()
	if err != nil {
		h.logger.Error("Failed to get league table: " + err.Error())
		http.Error(w, "Failed to get league table", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(league); err != nil {
		h.logger.Error("Failed to encode league table: " + err.Error())
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// @Summary Mevcut haftayı oynatır
// @Description Ligin güncel haftasını simüle eder ve ligi günceller
// @Tags league
// @Produce plain
// @Success 200 {string} string "Week played successfully"
// @Failure 500 {string} string "Internal server error"
// @Router /play-week [post]
func (h *LeagueHandler) PlayWeek(w http.ResponseWriter, r *http.Request) {
	week, err := h.leagueSvc.GetCurrentWeek()
	if err != nil {
		h.logger.Error("Failed to get current week: " + err.Error())
		http.Error(w, "Failed to get current week", http.StatusInternalServerError)
		return
	}

	if week > 6 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("League has already completed"))
		return
	}

	if err := h.leagueSvc.PlayWeek(week); err != nil {
		h.logger.Error("Failed to play week: " + err.Error())
		http.Error(w, "Failed to play week", http.StatusInternalServerError)
		return
	}

	matches, err := h.leagueSvc.GetMatchesByWeek(week)
	if err != nil {
		h.logger.Error("Failed to get matches: " + err.Error())
		http.Error(w, "Failed to get matches", http.StatusInternalServerError)
		return
	}

	result := fmt.Sprintf("Week %d Results:\n", week)
	for _, match := range matches {
		homeTeam, err := h.leagueSvc.GetTeamByID(match.HomeTeamID)
		if err != nil {
			h.logger.Error("Failed to get home team: " + err.Error())
			http.Error(w, "Failed to get home team", http.StatusInternalServerError)
			return
		}
		awayTeam, err := h.leagueSvc.GetTeamByID(match.AwayTeamID)
		if err != nil {
			h.logger.Error("Failed to get away team: " + err.Error())
			http.Error(w, "Failed to get away team", http.StatusInternalServerError)
			return
		}
		result += fmt.Sprintf("%s vs %s: %d - %d\n", homeTeam.Name, awayTeam.Name, match.HomeGoals, match.AwayGoals)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

// @Summary Ligi sıfırlar
// @Description Tüm takımları ve maçları sıfırlayarak ligi yeniden başlatır
// @Tags league
// @Produce plain
// @Success 200 {string} string "League reset successfully"
// @Failure 500 {string} string "Internal server error"
// @Router /reset-league [post]
func (h *LeagueHandler) ResetLeague(w http.ResponseWriter, r *http.Request) {
	if err := h.leagueSvc.ResetLeague(); err != nil {
		h.logger.Error("Failed to reset league: " + err.Error())
		http.Error(w, "Failed to reset league", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("League reset successfully"))
}

// @Summary Tüm ligi simüle eder
// @Description Ligdeki kalan tüm haftaları otomatik olarak simüle eder ve sonuçları döndürür
// @Tags league
// @Produce json
// @Success 200 {array} models.Match "Tüm simüle edilmiş maçların sonuçları"
// @Failure 500 {string} string "Internal server error"
// @Router /simulate-all-weeks [post]
func (h *LeagueHandler) SimulateAllWeeks(w http.ResponseWriter, r *http.Request) {
	simulatedMatches, err := h.leagueSvc.SimulateAllWeeks()
	if err != nil {
		if err.Error() == "league has already completed. current week: 7" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("League simulation completed previously. Use /reset-league to start a new season."))
			return
		}
		h.logger.Error("Failed to simulate all weeks: " + err.Error())
		http.Error(w, "Failed to simulate all weeks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(simulatedMatches); err != nil {
		h.logger.Error("Failed to encode simulated matches: " + err.Error())
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
