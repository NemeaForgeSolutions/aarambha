package service

import (
	"fmt"
	"net/url"

	"aarambha-agency/internal/domain"
)

type LeadService struct {
	leadRepo domain.LeadRepository
}

func NewLeadService(repo domain.LeadRepository) *LeadService {
	return &LeadService{leadRepo: repo}
}

type LeadResult struct {
	Lead        domain.Lead
	WhatsAppURL string
}

func (s *LeadService) SubmitLead(lead domain.Lead) (*LeadResult, error) {
	if err := lead.Validate(); err != nil {
		return nil, err
	}

	if err := s.leadRepo.Save(&lead); err != nil {
		return nil, fmt.Errorf("failed to save lead: %w", err)
	}

	waMsg := fmt.Sprintf("Hi Aarambha, I submitted an inquiry. Name: %s, Phone: %s.", lead.Name, lead.Phone)
	if lead.PropertyTitle != "" {
		waMsg = fmt.Sprintf("Hi Aarambha, I just submitted an inquiry regarding %s (%s). Name: %s.", lead.PropertyTitle, lead.PropertyID, lead.Name)
	}
	waURL := fmt.Sprintf("https://wa.me/917386985852?text=%s", url.QueryEscape(waMsg))

	return &LeadResult{
		Lead:        lead,
		WhatsAppURL: waURL,
	}, nil
}
