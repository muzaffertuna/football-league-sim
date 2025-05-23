package repositories

import (
	"database/sql"

	"github.com/muzaffertuna/football-league-sim/internal/app/models"
	"github.com/muzaffertuna/football-league-sim/internal/database"
)

type matchRepository struct {
	db *database.DB
}

func NewMatchRepository(db *database.DB) MatchRepository {
	return &matchRepository{db: db}
}

func (r *matchRepository) CreateMatch(match *models.Match) error {
	query := `
        INSERT INTO Matches (HomeTeamID, AwayTeamID, HomeGoals, AwayGoals, Week, Played)
        VALUES (@p1, @p2, @p3, @p4, @p5, @p6);
        SELECT SCOPE_IDENTITY();`
	var id int
	err := r.db.QueryRow(query,
		sql.Named("p1", match.HomeTeamID),
		sql.Named("p2", match.AwayTeamID),
		sql.Named("p3", match.HomeGoals),
		sql.Named("p4", match.AwayGoals),
		sql.Named("p5", match.Week),
		sql.Named("p6", match.Played),
	).Scan(&id)
	if err != nil {
		return err
	}
	match.ID = id
	return nil
}

func (r *matchRepository) GetMatchByID(id int) (*models.Match, error) {
	query := `
        SELECT ID, HomeTeamID, AwayTeamID, HomeGoals, AwayGoals, Week, Played
        FROM Matches
        WHERE ID = @p1`
	match := &models.Match{}
	err := r.db.QueryRow(query, sql.Named("p1", id)).Scan(
		&match.ID,
		&match.HomeTeamID,
		&match.AwayTeamID,
		&match.HomeGoals,
		&match.AwayGoals,
		&match.Week,
		&match.Played,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return match, nil
}

func (r *matchRepository) GetMatchesByWeek(week int) ([]models.Match, error) {
	query := `
        SELECT ID, HomeTeamID, AwayTeamID, HomeGoals, AwayGoals, Week, Played
        FROM Matches
        WHERE Week = @p1`
	rows, err := r.db.Query(query, sql.Named("p1", week))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := []models.Match{}
	for rows.Next() {
		match := models.Match{}
		if err := rows.Scan(
			&match.ID,
			&match.HomeTeamID,
			&match.AwayTeamID,
			&match.HomeGoals,
			&match.AwayGoals,
			&match.Week,
			&match.Played,
		); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	return matches, nil
}

func (r *matchRepository) UpdateMatch(match *models.Match) error {
	query := `
        UPDATE Matches
        SET HomeTeamID = @p1, AwayTeamID = @p2, HomeGoals = @p3, AwayGoals = @p4, Week = @p5, Played = @p6
        WHERE ID = @p7`
	_, err := r.db.Exec(query,
		sql.Named("p1", match.HomeTeamID),
		sql.Named("p2", match.AwayTeamID),
		sql.Named("p3", match.HomeGoals),
		sql.Named("p4", match.AwayGoals),
		sql.Named("p5", match.Week),
		sql.Named("p6", match.Played),
		sql.Named("p7", match.ID),
	)
	return err
}
