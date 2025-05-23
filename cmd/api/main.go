package main

import (
	"net/http"

	"github.com/muzaffertuna/football-league-sim/config"
	"github.com/muzaffertuna/football-league-sim/internal/app/handlers"
	"github.com/muzaffertuna/football-league-sim/internal/app/repositories"
	"github.com/muzaffertuna/football-league-sim/internal/app/services"
	"github.com/muzaffertuna/football-league-sim/internal/database"
	"github.com/muzaffertuna/football-league-sim/internal/platform"
)

func main() {
	logger := platform.NewLogger()

	cfg := config.LoadConfig()

	db, err := database.ConnectMSSQL(cfg.DBConnectionString, logger)
	if err != nil {
		logger.Error("Database connection failed: " + err.Error())
		return
	}
	defer db.Close()

	teamRepo := repositories.NewTeamRepository(db)
	matchRepo := repositories.NewMatchRepository(db)
	leagueRepo := repositories.NewLeagueRepository(teamRepo, matchRepo)

	teamSvc := services.NewTeamService(teamRepo)
	matchSvc := services.NewMatchService(matchRepo, teamRepo)
	leagueSvc := services.NewLeagueService(leagueRepo, matchRepo, matchSvc, teamRepo, teamSvc)

	leagueHandler := handlers.NewLeagueHandler(leagueSvc, logger)

	router := platform.NewRouter(leagueHandler)

	logger.Info("Starting server on " + cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, router); err != nil {
		logger.Error("Server failed to start: " + err.Error())
	}
}
