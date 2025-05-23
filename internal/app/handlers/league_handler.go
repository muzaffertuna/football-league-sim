package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/muzaffertuna/football-league-sim/internal/app/services"
	"github.com/muzaffertuna/football-league-sim/internal/platform"

	_ "github.com/muzaffertuna/football-league-sim/docs" // Swagger dokümantasyonu için
)

// @title Football League Simulation API
// @version 1.0
// @description Premier Lig simülasyonu için REST API
// @host localhost:8080
// @BasePath /
type LeagueHandler struct {
    leagueSvc services.LeagueService
    logger    *platform.Logger
}

func NewLeagueHandler(leagueSvc services.LeagueService, logger *platform.Logger) *LeagueHandler {
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

// @Summary Belirli bir haftayı oynatır
// @Description Belirtilen hafta için maçları simüle eder ve ligi günceller
// @Tags league
// @Produce plain
// @Param week query int true "Hafta numarası"
// @Success 200 {string} string "Week played successfully"
// @Failure 400 {string} string "Invalid week parameter"
// @Failure 500 {string} string "Internal server error"
// @Router /play-week [post]
func (h *LeagueHandler) PlayWeek(w http.ResponseWriter, r *http.Request) {
    weekStr := r.URL.Query().Get("week")
    week, err := strconv.Atoi(weekStr)
    if err != nil {
        h.logger.Error("Invalid week parameter: " + err.Error())
        http.Error(w, "Invalid week parameter", http.StatusBadRequest)
        return
    }

    if err := h.leagueSvc.PlayWeek(week); err != nil {
        h.logger.Error("Failed to play week: " + err.Error())
        http.Error(w, "Failed to play week", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Week played successfully"))
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