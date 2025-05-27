package models

type League struct {
	Teams                   []Team       `json:"teams"`
	Matches                 []Match      `json:"matches"`
	CurrentWeek             int          `json:"current_week"`
	ChampionshipPredictions []Prediction `json:"championshipPredictions"`
}
