package database

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/muzaffertuna/football-league-sim/internal/platform"
)

type DB struct {
	*sql.DB
	logger *platform.Logger
}

func ConnectMSSQL(connectionString string, logger *platform.Logger) (*DB, error) {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to connect to MSSQL: %v", err))
		return nil, fmt.Errorf("failed to connect to MSSQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		logger.Error(fmt.Sprintf("failed to ping MSSQL: %v", err))
		return nil, fmt.Errorf("failed to ping MSSQL: %w", err)
	}

	logger.Info("Successfully connected to MSSQL")
	return &DB{db, logger}, nil
}
