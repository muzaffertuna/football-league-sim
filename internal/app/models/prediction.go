package models

// Prediction, bir takımın belirli bir senaryo (örn. şampiyonluk) için olasılığını temsil eder.
type Prediction struct {
	TeamID                 int     `json:"team_id"`
	TeamName               string  `json:"team_name"`
	ChampionshipLikelihood float64 `json:"championship_likelihood"` // Şampiyonluk olasılığı (%)
	// İstenirse buraya diğer tahminler de eklenebilir:
	// Top4Likelihood float64 `json:"top_4_likelihood"` // İlk 4'e girme olasılığı
	// RelegationLikelihood float64 `json:"relegation_likelihood"` // Küme düşme olasılığı
}

// PredictionResult, birden fazla tahmin türünü içerebilir.
type PredictionResult struct {
	ChampionshipPredictions []Prediction `json:"championship_predictions"`
	// Diğer tahmin türleri buraya eklenebilir.
}
