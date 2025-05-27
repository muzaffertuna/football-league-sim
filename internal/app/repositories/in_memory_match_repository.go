package repositories

import (
	"fmt"
	"sync"

	"github.com/muzaffertuna/football-league-sim/internal/app/models"
)

// InMemoryMatchRepository MatchRepository arayüzünü bellek içi (in-memory) olarak uygular.
type InMemoryMatchRepository struct {
	mu      sync.RWMutex
	matches map[int]models.Match
	nextID  int
}

// NewInMemoryMatchRepository bellek içi maç deposunun yeni bir örneğini oluşturur.
func NewInMemoryMatchRepository() *InMemoryMatchRepository {
	return &InMemoryMatchRepository{
		matches: make(map[int]models.Match),
		nextID:  1,
	}
}

// CreateMatch yeni bir maç oluşturur. İMZA Orijinal MatchRepository ile AYNI: *models.Match alır.
func (r *InMemoryMatchRepository) CreateMatch(match *models.Match) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if match.ID == 0 {
		match.ID = r.nextID
		r.nextID++
	} else {
		if match.ID >= r.nextID {
			r.nextID = match.ID + 1
		}
	}
	r.matches[match.ID] = *match
	return nil
}

func (r *InMemoryMatchRepository) GetMatchByID(id int) (*models.Match, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	match, ok := r.matches[id]
	if !ok {
		return nil, nil
	}
	return &match, nil
}

func (r *InMemoryMatchRepository) GetMatchesByWeek(week int) ([]models.Match, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var weekMatches []models.Match
	for _, match := range r.matches {
		if match.Week == week {
			weekMatches = append(weekMatches, match)
		}
	}
	return weekMatches, nil
}

func (r *InMemoryMatchRepository) GetAllMatches() ([]models.Match, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var allMatches []models.Match
	for _, match := range r.matches {
		allMatches = append(allMatches, match)
	}
	return allMatches, nil
}

func (r *InMemoryMatchRepository) UpdateMatch(match *models.Match) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.matches[match.ID]; !ok {
		return fmt.Errorf("match with ID %d not found for update", match.ID)
	}
	r.matches[match.ID] = *match
	return nil
}

func (r *InMemoryMatchRepository) DeleteAllMatches() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.matches = make(map[int]models.Match)
	r.nextID = 1
	return nil
}

func (r *InMemoryMatchRepository) GetPlayedMatches() ([]models.Match, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var playedMatches []models.Match
	for _, match := range r.matches {
		if match.Played {
			playedMatches = append(playedMatches, match)
		}
	}
	return playedMatches, nil
}

func (r *InMemoryMatchRepository) GetMaxWeekPlayed() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	maxWeek := 0
	for _, match := range r.matches {
		if match.Played && match.Week > maxWeek {
			maxWeek = match.Week
		}
	}
	return maxWeek, nil
}
