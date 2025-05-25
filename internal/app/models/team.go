package models

type Team struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Strength      int    `json:"strength"`
	Points        int    `json:"points"`
	GoalsFor      int    `json:"goals_for"`
	GoalsAgainst  int    `json:"goals_against"`
	MatchesPlayed int    `json:"matches_played"`
	Wins          int    `json:"wins"`  // Yeni eklendi
	Draws         int    `json:"draws"` // Yeni eklendi
	Loses         int    `json:"loses"` // Yeni eklendi
}

func (t *Team) GoalDifference() int {
	return t.GoalsFor - t.GoalsAgainst
}
