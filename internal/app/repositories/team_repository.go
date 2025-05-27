package repositories

import (
	"database/sql"

	"github.com/muzaffertuna/football-league-sim/internal/app/models"
	"github.com/muzaffertuna/football-league-sim/internal/database"
)

type teamRepository struct {
	db *database.DB
}

func NewTeamRepository(db *database.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) CreateTeam(team *models.Team) error {
	query := `
		INSERT INTO Teams (Name, Strength, Points, GoalsFor, GoalsAgainst, MatchesPlayed, Wins, Draws, Loses)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9);
		SELECT SCOPE_IDENTITY();`
	var id int
	err := r.db.QueryRow(query,
		sql.Named("p1", team.Name),
		sql.Named("p2", team.Strength),
		sql.Named("p3", team.Points),
		sql.Named("p4", team.GoalsFor),
		sql.Named("p5", team.GoalsAgainst),
		sql.Named("p6", team.MatchesPlayed),
		sql.Named("p7", team.Wins),  // Yeni eklendi
		sql.Named("p8", team.Draws), // Yeni eklendi
		sql.Named("p9", team.Loses), // Yeni eklendi
	).Scan(&id)
	if err != nil {
		return err
	}
	team.ID = id
	return nil
}

func (r *teamRepository) GetTeamByID(id int) (*models.Team, error) {
	query := `
		SELECT ID, Name, Strength, Points, GoalsFor, GoalsAgainst, MatchesPlayed, Wins, Draws, Loses
		FROM Teams
		WHERE ID = @p1`
	team := &models.Team{}
	err := r.db.QueryRow(query, sql.Named("p1", id)).Scan(
		&team.ID,
		&team.Name,
		&team.Strength,
		&team.Points,
		&team.GoalsFor,
		&team.GoalsAgainst,
		&team.MatchesPlayed,
		&team.Wins,  // Yeni eklendi
		&team.Draws, // Yeni eklendi
		&team.Loses, // Yeni eklendi
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (r *teamRepository) GetAllTeams() ([]models.Team, error) {
	query := `
		SELECT ID, Name, Strength, Points, GoalsFor, GoalsAgainst, MatchesPlayed, Wins, Draws, Loses
		FROM Teams`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teams := []models.Team{}
	for rows.Next() {
		team := models.Team{}
		if err := rows.Scan(
			&team.ID,
			&team.Name,
			&team.Strength,
			&team.Points,
			&team.GoalsFor,
			&team.GoalsAgainst,
			&team.MatchesPlayed,
			&team.Wins,
			&team.Draws,
			&team.Loses,
		); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func (r *teamRepository) UpdateTeam(team *models.Team) error {
	query := `
		UPDATE Teams
		SET Name = @p1, Strength = @p2, Points = @p3, GoalsFor = @p4, GoalsAgainst = @p5, MatchesPlayed = @p6, Wins = @p7, Draws = @p8, Loses = @p9
		WHERE ID = @p10`
	_, err := r.db.Exec(query,
		sql.Named("p1", team.Name),
		sql.Named("p2", team.Strength),
		sql.Named("p3", team.Points),
		sql.Named("p4", team.GoalsFor),
		sql.Named("p5", team.GoalsAgainst),
		sql.Named("p6", team.MatchesPlayed),
		sql.Named("p7", team.Wins),
		sql.Named("p8", team.Draws),
		sql.Named("p9", team.Loses),
		sql.Named("p10", team.ID),
	)
	return err
}
