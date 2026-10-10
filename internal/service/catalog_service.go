package service

import (
	"aarambha-agency/internal/domain"
)

type CatalogService struct {
	catalogRepo domain.CatalogRepository
}

func NewCatalogService(repo domain.CatalogRepository) *CatalogService {
	return &CatalogService{catalogRepo: repo}
}

func (s *CatalogService) GetAgents() []domain.Agent {
	return s.catalogRepo.FindAgents()
}

func (s *CatalogService) GetAgentByID(id string) (*domain.Agent, error) {
	return s.catalogRepo.FindAgentByID(id)
}

func (s *CatalogService) CreateAgent(agent *domain.Agent) error {
	return s.catalogRepo.CreateAgent(agent)
}

func (s *CatalogService) UpdateAgent(agent *domain.Agent) error {
	return s.catalogRepo.UpdateAgent(agent)
}

func (s *CatalogService) DeleteAgent(id string) error {
	return s.catalogRepo.DeleteAgent(id)
}

func (s *CatalogService) GetProjects() []domain.Project {
	return s.catalogRepo.FindProjects()
}

func (s *CatalogService) GetLocations() []domain.LocationInsight {
	return s.catalogRepo.FindLocations()
}

func (s *CatalogService) GetInsights() []domain.InsightArticle {
	return s.catalogRepo.FindInsights()
}

func (s *CatalogService) GetTestimonials() []domain.Testimonial {
	return s.catalogRepo.FindTestimonials()
}

func (s *CatalogService) GetWhyChooseUs() []domain.WhyChooseUsItem {
	return s.catalogRepo.FindWhyChooseUs()
}
