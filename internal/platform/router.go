package platform

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	// "github.com/muzaffertuna/football-league-sim/internal/app/handlers" // <--- BU SATIR KALDIRILDI!
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter fonksiyonunun LeagueHandlerContract arayüzünü alması gerekiyor.
func NewRouter(leagueHandler LeagueHandlerContract) http.Handler { // <--- Düzeltildi: *handlers.LeagueHandler yerine LeagueHandlerContract
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	r.Get("/league-table", leagueHandler.GetLeagueTable)
	r.Post("/play-week", leagueHandler.PlayWeek)
	r.Post("/reset-league", leagueHandler.ResetLeague)
	// r.Get("/fixture", leagueHandler.GetFixture)
	r.Post("/simulate-all-weeks", leagueHandler.SimulateAllWeeks)

	return r
}
