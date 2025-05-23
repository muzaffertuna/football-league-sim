// platform/interfaces.go
package platform

import "net/http"

type LeagueHandler interface {
	GetLeagueTable(w http.ResponseWriter, r *http.Request)
	PlayWeek(w http.ResponseWriter, r *http.Request)
	ResetLeague(w http.ResponseWriter, r *http.Request)
}
