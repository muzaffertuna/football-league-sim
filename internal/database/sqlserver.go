package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/muzaffertuna/football-league-sim/internal/pkg/logger"
)

type DB struct {
	*sql.DB
	logger *logger.Logger
}

// ConnectMSSQL MSSQL veritabanına bağlanır ve eğer yoksa veritabanını otomatik olarak oluşturur.
func ConnectMSSQL(connectionString string, logger *logger.Logger) (*DB, error) {
	parsedURL, err := url.Parse(connectionString)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to parse connection string: %v", err))
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	dbName := parsedURL.Query().Get("database")
	if dbName == "" {
		logger.Error("Database name not found in connection string.")
		return nil, fmt.Errorf("database name not found in connection string")
	}

	masterConnectionString := strings.Replace(connectionString, fmt.Sprintf("database=%s", dbName), "database=master", 1)

	logger.Info("Attempting to connect to MSSQL master database to check/create target DB...")

	dbMaster, err := sql.Open("sqlserver", masterConnectionString)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to open master database connection: %v", err))
		return nil, fmt.Errorf("failed to open master database connection: %w", err)
	}
	defer dbMaster.Close()

	err = dbMaster.Ping()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to ping master database: %v", err))
		return nil, fmt.Errorf("failed to ping master database: %w", err)
	}
	logger.Info("Successfully connected to MSSQL master database.")

	var exists int
	query := fmt.Sprintf("SELECT 1 FROM sys.databases WHERE name = '%s'", dbName)
	err = dbMaster.QueryRow(query).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(fmt.Sprintf("Failed to check if database '%s' exists: %v", dbName, err))
		return nil, fmt.Errorf("failed to check if database '%s' exists: %w", dbName, err)
	}

	if exists == 0 {
		logger.Info(fmt.Sprintf("Database '%s' does not exist. Creating it...", dbName))
		createDBQuery := fmt.Sprintf("CREATE DATABASE %s", dbName)
		_, err = dbMaster.Exec(createDBQuery)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to create database '%s': %v", dbName, err))
			return nil, fmt.Errorf("failed to create database '%s': %w", dbName, err)
		}
		logger.Info(fmt.Sprintf("Database '%s' created successfully.", dbName))
	} else {
		logger.Info(fmt.Sprintf("Database '%s' already exists. Skipping creation.", dbName))
	}

	logger.Info(fmt.Sprintf("Attempting to connect to target database '%s'...", dbName))
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to open target database connection: %v", err))
		return nil, fmt.Errorf("failed to open target database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to ping target database: %v", err))
		return nil, fmt.Errorf("failed to ping target database: %w", err)
	}

	logger.Info("Successfully connected to MSSQL target database.")
	return &DB{db, logger}, nil
}
