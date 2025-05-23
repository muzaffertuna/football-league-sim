package models

type Match struct {
	ID         int  `json:"id"`
	HomeTeamID int  `json:"home_team_id"`
	AwayTeamID int  `json:"away_team_id"`
	HomeGoals  int  `json:"home_goals"`
	AwayGoals  int  `json:"away_goals"`
	Week       int  `json:"week"`
	Played     bool `json:"played"`
}
