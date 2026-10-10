package memory

import (
	"fmt"
	"sync"
	"time"

	"aarambha-agency/internal/domain"
)

type LeadRepo struct {
	mu    sync.RWMutex
	leads []domain.Lead
}

func NewLeadRepo() *LeadRepo {
	return &LeadRepo{
		leads: make([]domain.Lead, 0),
	}
}

func (r *LeadRepo) Save(lead *domain.Lead) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	lead.ID = fmt.Sprintf("LEAD-%d", time.Now().UnixNano())
	lead.CreatedAt = time.Now()
	if lead.Status == "" {
		lead.Status = "pending"
	}
	r.leads = append(r.leads, *lead)
	return nil
}

func (r *LeadRepo) UpdateStatus(id string, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, l := range r.leads {
		if l.ID == id {
			r.leads[i].Status = status
			return nil
		}
	}
	return domain.ErrNotFound
}

func (r *LeadRepo) FindByID(id string) (*domain.Lead, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, l := range r.leads {
		if l.ID == id {
			item := l
			return &item, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (r *LeadRepo) FindAll() []domain.Lead {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]domain.Lead, len(r.leads))
	copy(res, r.leads)
	return res
}
