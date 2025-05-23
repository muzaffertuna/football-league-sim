package models

type Team struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Strength      int    `json:"strength"`
	Points        int    `json:"points"`
	GoalsFor      int    `json:"goals_for"`
	GoalsAgainst  int    `json:"goals_against"`
	MatchesPlayed int    `json:"matches_played"`
}

func (t *Team) GoalDifference() int {
	return t.GoalsFor - t.GoalsAgainst
}
