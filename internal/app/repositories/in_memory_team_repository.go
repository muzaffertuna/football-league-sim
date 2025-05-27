package repositories

import (
	"fmt"
	"sync"

	"github.com/muzaffertuna/football-league-sim/internal/app/models"
)

type InMemoryTeamRepository struct {
	mu     sync.RWMutex
	teams  map[int]models.Team
	nextID int
}

func NewInMemoryTeamRepository() *InMemoryTeamRepository {
	return &InMemoryTeamRepository{
		teams:  make(map[int]models.Team),
		nextID: 1,
	}
}

func (r *InMemoryTeamRepository) CreateTeam(team *models.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if team.ID == 0 {
		team.ID = r.nextID
		r.nextID++
	} else {
		if team.ID >= r.nextID {
			r.nextID = team.ID + 1
		}
	}
	r.teams[team.ID] = *team
	return nil
}

func (r *InMemoryTeamRepository) GetTeamByID(id int) (*models.Team, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	team, ok := r.teams[id]
	if !ok {
		return nil, nil
	}
	return &team, nil
}

func (r *InMemoryTeamRepository) GetAllTeams() ([]models.Team, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var allTeams []models.Team
	for _, team := range r.teams {
		allTeams = append(allTeams, team)
	}
	return allTeams, nil
}

func (r *InMemoryTeamRepository) UpdateTeam(team *models.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.teams[team.ID]; !ok {
		return fmt.Errorf("team with ID %d not found for update", team.ID)
	}
	r.teams[team.ID] = *team
	return nil
}

func (r *InMemoryTeamRepository) DeleteAllTeams() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.teams = make(map[int]models.Team)
	r.nextID = 1
	return nil
}
