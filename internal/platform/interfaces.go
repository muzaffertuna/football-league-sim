package platform

import "net/http"

// LeagueHandlerContract router'ın LeagueHandler'dan beklediği metotları tanımlar.
// Bu arayüz, platform paketinin doğrudan handlers paketine bağımlılığını kırar.
type LeagueHandlerContract interface {
	GetLeagueTable(w http.ResponseWriter, r *http.Request)
	PlayWeek(w http.ResponseWriter, r *http.Request)
	ResetLeague(w http.ResponseWriter, r *http.Request)
	GetFixture(w http.ResponseWriter, r *http.Request)
	SimulateAllWeeks(w http.ResponseWriter, r *http.Request)
}
