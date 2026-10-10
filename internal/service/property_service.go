package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"aarambha-agency/internal/domain"
)

type PropertyService struct {
	propertyRepo domain.PropertyRepository
	catalogRepo  domain.CatalogRepository
}

func NewPropertyService(pRepo domain.PropertyRepository, cRepo domain.CatalogRepository) *PropertyService {
	return &PropertyService{
		propertyRepo: pRepo,
		catalogRepo:  cRepo,
	}
}

func (s *PropertyService) ListProperties(f domain.PropertyFilter) []domain.Property {
	return s.propertyRepo.FindByFilter(f)
}

func (s *PropertyService) GetFeatured(propType string, limit int) []domain.Property {
	return s.propertyRepo.FindFeatured(propType, limit)
}

func (s *PropertyService) GetDetail(slugOrID string) (*domain.Property, *domain.Agent, []domain.Property, error) {
	prop, err := s.propertyRepo.FindBySlug(slugOrID)
	if err != nil {
		return nil, nil, nil, err
	}

	var agent *domain.Agent
	if prop.AgentID != "" {
		if a, aErr := s.catalogRepo.FindAgentByID(prop.AgentID); aErr == nil {
			agent = a
		}
	}

	similar := s.propertyRepo.FindSimilar(prop, 3)

	return prop, agent, similar, nil
}

func (s *PropertyService) GetByID(idOrSlug string) (*domain.Property, error) {
	return s.propertyRepo.FindBySlug(idOrSlug)
}

func (s *PropertyService) GetCommercial() []domain.Property {
	return s.propertyRepo.FindByFilter(domain.PropertyFilter{
		Type: "Commercial",
	})
}

func (s *PropertyService) GetRentals() []domain.Property {
	return s.propertyRepo.FindByFilter(domain.PropertyFilter{
		Transaction: "Rent",
	})
}

func (s *PropertyService) TotalCount() int {
	return s.propertyRepo.Count()
}

func (s *PropertyService) Create(prop *domain.Property) error {
	if strings.TrimSpace(prop.Title) == "" {
		return errors.New("property title is required")
	}
	if strings.TrimSpace(prop.City) == "" {
		return errors.New("city is required")
	}
	if prop.Price <= 0 {
		return errors.New("price must be greater than zero")
	}

	if strings.TrimSpace(prop.ID) == "" {
		cityCode := "PROP"
		if len(prop.City) >= 3 {
			cityCode = strings.ToUpper(prop.City[:3])
		}
		prop.ID = fmt.Sprintf("PROP-%s-%d", cityCode, time.Now().Unix()%100000)
	}

	if strings.TrimSpace(prop.Slug) == "" {
		prop.Slug = Slugify(prop.Title) + "-" + strings.ToLower(prop.ID)
	}

	if strings.TrimSpace(prop.PriceDisplay) == "" {
		prop.PriceDisplay = FormatIndianCurrency(prop.Price)
		if strings.EqualFold(prop.Transaction, "Rent") {
			prop.PriceDisplay += "/month"
		}
	}

	if strings.TrimSpace(prop.PricePerSqFt) == "" && prop.Area > 0 {
		sqFtRate := prop.Price / int64(prop.Area)
		prop.PricePerSqFt = fmt.Sprintf("₹%s/sq.ft.", formatWithCommas(sqFtRate))
	}

	if strings.TrimSpace(prop.Image) == "" {
		prop.Image = "https://images.unsplash.com/photo-1600585154340-be6161a56a0c?auto=format&fit=crop&w=1200&q=85"
	}

	if prop.CarpetArea <= 0 && prop.Area > 0 {
		prop.CarpetArea = int(float64(prop.Area) * 0.8)
	}

	return s.propertyRepo.Create(prop)
}

func (s *PropertyService) Update(prop *domain.Property) error {
	if strings.TrimSpace(prop.ID) == "" {
		return errors.New("property ID is required for update")
	}
	if strings.TrimSpace(prop.Title) == "" {
		return errors.New("property title is required")
	}

	if strings.TrimSpace(prop.PriceDisplay) == "" {
		prop.PriceDisplay = FormatIndianCurrency(prop.Price)
		if strings.EqualFold(prop.Transaction, "Rent") {
			prop.PriceDisplay += "/month"
		}
	}

	if strings.TrimSpace(prop.PricePerSqFt) == "" && prop.Area > 0 {
		sqFtRate := prop.Price / int64(prop.Area)
		prop.PricePerSqFt = fmt.Sprintf("₹%s/sq.ft.", formatWithCommas(sqFtRate))
	}

	return s.propertyRepo.Update(prop)
}

func (s *PropertyService) Delete(id string) error {
	return s.propertyRepo.Delete(id)
}

func (s *PropertyService) ToggleFeatured(id string) (bool, error) {
	return s.propertyRepo.ToggleFeatured(id)
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9]+`)

func Slugify(s string) string {
	s = strings.ToLower(s)
	s = nonAlphanumericRegex.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func FormatIndianCurrency(amount int64) string {
	if amount >= 10000000 {
		val := float64(amount) / 10000000.0
		str := fmt.Sprintf("%.2f", val)
		str = strings.TrimSuffix(str, ".00")
		return fmt.Sprintf("₹%s Cr", str)
	} else if amount >= 100000 {
		val := float64(amount) / 100000.0
		str := fmt.Sprintf("%.2f", val)
		str = strings.TrimSuffix(str, ".00")
		return fmt.Sprintf("₹%s Lakh", str)
	}
	return fmt.Sprintf("₹%d", amount)
}

func formatWithCommas(n int64) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	last3 := s[len(s)-3:]
	remaining := s[:len(s)-3]
	var parts []string
	for len(remaining) > 2 {
		parts = append([]string{remaining[len(remaining)-2:]}, parts...)
		remaining = remaining[:len(remaining)-2]
	}
	if len(remaining) > 0 {
		parts = append([]string{remaining}, parts...)
	}
	return strings.Join(parts, ",") + "," + last3
}
