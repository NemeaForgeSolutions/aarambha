package memory

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"aarambha-agency/internal/domain"
	"aarambha-agency/internal/pkg/uuid"
)

type PropertyRepo struct {
	mu           sync.RWMutex
	companyInfo  domain.CompanyInfo
	agents       []domain.Agent
	properties   []domain.Property
	newProjects  []domain.Project
	locations    []domain.LocationInsight
	insights     []domain.InsightArticle
	testimonials []domain.Testimonial
	whyChooseUs  []domain.WhyChooseUsItem
}

func NewPropertyRepo() (*PropertyRepo, error) {
	return &PropertyRepo{
		companyInfo:  domain.CompanyInfo{},
		agents:       []domain.Agent{},
		properties:   []domain.Property{},
		newProjects:  []domain.Project{},
		locations:    []domain.LocationInsight{},
		insights:     []domain.InsightArticle{},
		testimonials: []domain.Testimonial{},
		whyChooseUs:  []domain.WhyChooseUsItem{},
	}, nil
}

func (r *PropertyRepo) FindAll() []domain.Property {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]domain.Property, len(r.properties))
	copy(res, r.properties)
	return res
}

func (r *PropertyRepo) FindByFilter(f domain.PropertyFilter) []domain.Property {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matches []domain.Property
	for _, p := range r.properties {
		if p.Matches(f) {
			matches = append(matches, p)
		}
	}

	switch f.Sort {
	case "price_asc":
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Price < matches[j].Price
		})
	case "price_desc":
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Price > matches[j].Price
		})
	case "area_desc":
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Area > matches[j].Area
		})
	default:
		sort.Slice(matches, func(i, j int) bool {
			if matches[i].Featured != matches[j].Featured {
				return matches[i].Featured
			}
			return matches[i].ID > matches[j].ID
		})
	}

	return matches
}

func (r *PropertyRepo) FindBySlug(slugOrID string) (*domain.Property, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.properties {
		if p.Slug == slugOrID || p.ID == slugOrID {
			item := p
			return &item, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (r *PropertyRepo) FindFeatured(propType string, limit int) []domain.Property {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []domain.Property
	for _, p := range r.properties {
		if propType != "" && !strings.EqualFold(propType, "all") && !strings.EqualFold(p.Type, propType) {
			continue
		}
		list = append(list, p)
	}

	if limit > 0 && len(list) > limit {
		return list[:limit]
	}
	return list
}

func (r *PropertyRepo) FindSimilar(current *domain.Property, limit int) []domain.Property {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var similar []domain.Property
	for _, p := range r.properties {
		if p.ID == current.ID {
			continue
		}
		if p.City == current.City || p.Type == current.Type {
			similar = append(similar, p)
			if len(similar) >= limit {
				break
			}
		}
	}
	return similar
}

func (r *PropertyRepo) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.properties)
}

func (r *PropertyRepo) FindAgents() []domain.Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]domain.Agent, len(r.agents))
	copy(res, r.agents)
	return res
}

func (r *PropertyRepo) FindAgentByID(id string) (*domain.Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, a := range r.agents {
		if a.ID == id {
			item := a
			return &item, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (r *PropertyRepo) CreateAgent(agent *domain.Agent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if agent.ID == "" {
		agent.ID = uuid.NewV7()
	}
	for _, a := range r.agents {
		if a.ID == agent.ID {
			return fmt.Errorf("agent with ID %s already exists", agent.ID)
		}
	}
	r.agents = append(r.agents, *agent)
	return nil
}

func (r *PropertyRepo) UpdateAgent(agent *domain.Agent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, a := range r.agents {
		if a.ID == agent.ID {
			r.agents[i] = *agent
			return nil
		}
	}
	return domain.ErrNotFound
}

func (r *PropertyRepo) DeleteAgent(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, a := range r.agents {
		if a.ID == id {
			r.agents = append(r.agents[:i], r.agents[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (r *PropertyRepo) FindProjects() []domain.Project {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]domain.Project, len(r.newProjects))
	copy(res, r.newProjects)
	return res
}

func (r *PropertyRepo) FindLocations() []domain.LocationInsight {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]domain.LocationInsight, len(r.locations))
	copy(res, r.locations)
	return res
}

func (r *PropertyRepo) FindInsights() []domain.InsightArticle {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]domain.InsightArticle, len(r.insights))
	copy(res, r.insights)
	return res
}

func (r *PropertyRepo) FindTestimonials() []domain.Testimonial {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]domain.Testimonial, len(r.testimonials))
	copy(res, r.testimonials)
	return res
}

func (r *PropertyRepo) FindWhyChooseUs() []domain.WhyChooseUsItem {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]domain.WhyChooseUsItem, len(r.whyChooseUs))
	copy(res, r.whyChooseUs)
	return res
}

func (r *PropertyRepo) Create(prop *domain.Property) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range r.properties {
		if p.ID == prop.ID {
			return fmt.Errorf("property with ID %s already exists", prop.ID)
		}
	}
	r.properties = append([]domain.Property{*prop}, r.properties...)
	return nil
}

func (r *PropertyRepo) Update(prop *domain.Property) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, p := range r.properties {
		if p.ID == prop.ID {
			r.properties[i] = *prop
			return nil
		}
	}
	return domain.ErrNotFound
}

func (r *PropertyRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, p := range r.properties {
		if p.ID == id {
			r.properties = append(r.properties[:i], r.properties[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (r *PropertyRepo) ToggleFeatured(id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, p := range r.properties {
		if p.ID == id {
			r.properties[i].Featured = !r.properties[i].Featured
			return r.properties[i].Featured, nil
		}
	}
	return false, domain.ErrNotFound
}
