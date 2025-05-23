package platform

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(leagueHandler LeagueHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	r.Get("/league-table", leagueHandler.GetLeagueTable)
	r.Post("/play-week", leagueHandler.PlayWeek)
	r.Post("/reset-league", leagueHandler.ResetLeague)

	return r
}
