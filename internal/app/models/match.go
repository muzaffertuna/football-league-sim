package models

// Match, bir maçı temsil eder
type Match struct {
    ID         int    `json:"id"`
    HomeTeamID int    `json:"home_team_id"`
    AwayTeamID int    `json:"away_team_id"`
    HomeGoals  int    `json:"home_goals"`
    AwayGoals  int    `json:"away_goals"`
    Week       int    `json:"week"` // Maçın oynandığı hafta
    Played     bool   `json:"played"` // Maç oynandı mı?
}