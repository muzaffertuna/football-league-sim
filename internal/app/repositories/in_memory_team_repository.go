package repositories

import (
	"fmt"
	"sync"

	"github.com/muzaffertuna/football-league-sim/internal/app/models"
)

// InMemoryTeamRepository TeamRepository arayüzünü bellek içi (in-memory) olarak uygular.
type InMemoryTeamRepository struct {
	mu     sync.RWMutex
	teams  map[int]models.Team
	nextID int // Otomatik ID atamak için
}

// NewInMemoryTeamRepository bellek içi takım deposunun yeni bir örneğini oluşturur.
func NewInMemoryTeamRepository() *InMemoryTeamRepository {
	return &InMemoryTeamRepository{
		teams:  make(map[int]models.Team),
		nextID: 1,
	}
}

// CreateTeam yeni bir takım oluşturur. İMZA Orijinal TeamRepository ile AYNI: *models.Team alır.
func (r *InMemoryTeamRepository) CreateTeam(team *models.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Eğer takımın ID'si zaten varsa, mevcut ID'yi kullan
	if team.ID == 0 {
		team.ID = r.nextID
		r.nextID++
	} else {
		// Eğer belirli bir ID ile oluşturuluyorsa, nextID'yi buna göre ayarla ve çakışmayı önle
		if team.ID >= r.nextID {
			r.nextID = team.ID + 1
		}
	}
	r.teams[team.ID] = *team // Pointer'dan değeri kopyala
	return nil
}

// GetTeamByID belirtilen ID'ye sahip takımı getirir.
func (r *InMemoryTeamRepository) GetTeamByID(id int) (*models.Team, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	team, ok := r.teams[id]
	if !ok {
		return nil, nil // sql.ErrNoRows yerine nil, nil dönüyoruz çünkü bellek içi hatalar daha farklı ele alınır
	}
	return &team, nil
}

// GetAllTeams tüm takımları getirir.
func (r *InMemoryTeamRepository) GetAllTeams() ([]models.Team, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var allTeams []models.Team
	for _, team := range r.teams {
		allTeams = append(allTeams, team)
	}
	return allTeams, nil
}

// UpdateTeam mevcut bir takımın bilgilerini günceller. İMZA Orijinal TeamRepository ile AYNI: *models.Team alır.
func (r *InMemoryTeamRepository) UpdateTeam(team *models.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.teams[team.ID]; !ok {
		return fmt.Errorf("team with ID %d not found for update", team.ID)
	}
	r.teams[team.ID] = *team // Pointer'dan değeri kopyala
	return nil
}

// DeleteAllTeams tüm takımları siler (sıfırlama için kullanılır).
// Bu metod orijinal TeamRepository arayüzünüzde tanımlı olmasa bile,
// InMemoryTeamRepository'de bulunmasında bir sakınca yoktur.
func (r *InMemoryTeamRepository) DeleteAllTeams() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.teams = make(map[int]models.Team)
	r.nextID = 1
	return nil
}
