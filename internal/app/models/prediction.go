package models

type Prediction struct {
	TeamID                 int     `json:"team_id"`
	TeamName               string  `json:"team_name"`
	ChampionshipLikelihood float64 `json:"championship_likelihood"` // Şampiyonluk olasılığı (%)
}

type PredictionResult struct {
	ChampionshipPredictions []Prediction `json:"championship_predictions"`
}
