package main

import (
	"fmt"
	"net/http"

	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlserver"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/microsoft/go-mssqldb" // sqlserver sürücüsü
	"github.com/muzaffertuna/football-league-sim/config"
	"github.com/muzaffertuna/football-league-sim/internal/app/handlers"
	"github.com/muzaffertuna/football-league-sim/internal/app/repositories"
	"github.com/muzaffertuna/football-league-sim/internal/app/services"
	"github.com/muzaffertuna/football-league-sim/internal/database"
	"github.com/muzaffertuna/football-league-sim/internal/pkg/logger"
	"github.com/muzaffertuna/football-league-sim/internal/platform"
)

func runMigrations(db *sql.DB, dbURL string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	driver, err := sqlserver.WithInstance(db, &sqlserver.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlserver",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func main() {
	// Logger'ı oluştur
	logger := logger.NewLogger()

	// Config'i yükle
	cfg := config.LoadConfig()

	// Veritabanına bağlan
	db, err := database.ConnectMSSQL(cfg.DBConnectionString, logger)
	if err != nil {
		logger.Error("Database connection failed: " + err.Error())
		return
	}
	defer db.Close()

	// Migration'ları çalıştır
	if err := runMigrations(db.DB, cfg.DBConnectionString); err != nil {
		logger.Error("Failed to run migrations: " + err.Error())
		return
	}

	// Repository'leri oluştur
	teamRepo := repositories.NewTeamRepository(db)
	matchRepo := repositories.NewMatchRepository(db)
	// leagueRepo := repositories.NewLeagueRepository(teamRepo, matchRepo) // Bu satır artık kullanılmıyor ve yorum satırı yapılmalı veya silinmeli

	// Servisleri oluştur
	teamSvc := services.NewTeamService(teamRepo)
	matchSvc := services.NewMatchService(matchRepo, teamRepo)

	// SADECE BU KISIM GÜNCELLENDİ:
	// services.NewLeagueService artık leagueRepo almıyor ve bir hata döndürüyor.
	leagueSvc, err := services.NewLeagueService(matchRepo, matchSvc, teamRepo, teamSvc)
	if err != nil {
		logger.Error("Failed to initialize league service: " + err.Error())
		return
	}

	// Handler'ları oluştur
	leagueHandler := handlers.NewLeagueHandler(leagueSvc, logger)

	// Router'ı oluştur
	router := platform.NewRouter(leagueHandler)

	// Sunucuyu başlat
	logger.Info("Starting server on " + cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, router); err != nil {
		logger.Error("Server failed to start: " + err.Error())
	}
}
