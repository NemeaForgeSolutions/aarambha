package service

import (
	"aarambha-agency/internal/domain"
)

type AdminService struct {
	propertyRepo domain.PropertyRepository
	leadRepo     domain.LeadRepository
}

func NewAdminService(pRepo domain.PropertyRepository, lRepo domain.LeadRepository) *AdminService {
	return &AdminService{
		propertyRepo: pRepo,
		leadRepo:     lRepo,
	}
}

func (s *AdminService) GetStats() domain.AdminStats {
	props := s.propertyRepo.FindAll()
	leads := s.leadRepo.FindAll()

	var totalValue int64
	featuredCount := 0
	pendingSellRequests := 0

	for _, p := range props {
		totalValue += p.Price
		if p.Featured {
			featuredCount++
		}
	}

	for _, l := range leads {
		if (l.Type == "sell_request" || l.Type == "valuation_request") && (l.Status == "pending" || l.Status == "") {
			pendingSellRequests++
		}
	}

	return domain.AdminStats{
		TotalProperties:            len(props),
		TotalLeads:                 len(leads),
		PendingSellRequests:        pendingSellRequests,
		FeaturedCount:              featuredCount,
		TotalPortfolioValue:        totalValue,
		TotalPortfolioValueDisplay: FormatIndianCurrency(totalValue),
	}
}

func (s *AdminService) GetRecentLeads() []domain.Lead {
	leads := s.leadRepo.FindAll()
	reversed := make([]domain.Lead, len(leads))
	for i, l := range leads {
		reversed[len(leads)-1-i] = l
	}
	if len(reversed) > 15 {
		return reversed[:15]
	}
	return reversed
}

func (s *AdminService) GetSellRequests() []domain.Lead {
	leads := s.leadRepo.FindAll()
	var sellRequests []domain.Lead
	for _, l := range leads {
		if l.Type == "sell_request" || l.Type == "valuation_request" {
			if l.Status == "" {
				l.Status = "pending"
			}
			sellRequests = append(sellRequests, l)
		}
	}
	// Sort newest first
	reversed := make([]domain.Lead, len(sellRequests))
	for i, l := range sellRequests {
		reversed[len(sellRequests)-1-i] = l
	}
	return reversed
}

func (s *AdminService) UpdateLeadStatus(id string, status string) error {
	switch status {
	case "pending", "approved", "rejected":
		return s.leadRepo.UpdateStatus(id, status)
	default:
		return domain.ErrNotFound
	}
}

func (s *AdminService) GetLeadByID(id string) (*domain.Lead, error) {
	return s.leadRepo.FindByID(id)
}
