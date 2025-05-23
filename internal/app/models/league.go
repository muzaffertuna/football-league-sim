package models

// League, lig tablosunu temsil eder
type League struct {
    Teams  []Team  `json:"teams"`
    Matches []Match `json:"matches"`
    CurrentWeek int `json:"current_week"`
}