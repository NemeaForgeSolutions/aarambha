package handler

import (
	"fmt"
	"net/http"

	"aarambha-agency/internal/domain"
	"aarambha-agency/internal/service"
	"aarambha-agency/internal/view"
)

type PageHandler struct {
	propertyService *service.PropertyService
	catalogService  *service.CatalogService
	engine          *view.Engine
	company         domain.CompanyInfo
}

func NewPageHandler(
	pService *service.PropertyService,
	cService *service.CatalogService,
	engine *view.Engine,
	company domain.CompanyInfo,
) *PageHandler {
	return &PageHandler{
		propertyService: pService,
		catalogService:  cService,
		engine:          engine,
		company:         company,
	}
}

func (h *PageHandler) render(w http.ResponseWriter, r *http.Request, page string, data view.PageData) {
	if data.User == nil && r != nil {
		data.User = view.UserFromContext(r.Context())
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.engine.RenderPage(w, page, data); err != nil {
		http.Error(w, fmt.Sprintf("Template rendering error: %v", err), http.StatusInternalServerError)
	}
}

// NotFound renders the branded 404 page with a 404 HTTP status code.
func (h *PageHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	pageData := view.PageData{
		Title:       "Page Not Found",
		Description: "The page you're looking for doesn't exist. Browse our luxury property portfolio or contact our advisory team.",
		ActivePage:  "",
		Robots:      "noindex, nofollow",
		Company:     h.company,
		User:        view.UserFromContext(r.Context()),
	}
	if err := h.engine.RenderPage(w, "404.html", pageData); err != nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	}
}

func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	featured := h.propertyService.GetFeatured("all", 6)
	projects := h.catalogService.GetProjects()
	testimonials := h.catalogService.GetTestimonials()
	whyChooseUs := h.catalogService.GetWhyChooseUs()

	data := struct {
		Featured      []domain.Property
		Projects      []domain.Project
		Testimonials  []domain.Testimonial
		WhyChooseUs   []domain.WhyChooseUsItem
		TotalListings int
	}{
		Featured:      featured,
		Projects:      projects,
		Testimonials:  testimonials,
		WhyChooseUs:   whyChooseUs,
		TotalListings: h.propertyService.TotalCount(),
	}

	pageData := view.PageData{
		Title:        "Premium Indian Real Estate Brokerage & Property Advisory",
		Description:  "Aarambha Agency is an ethical Indian real estate brokerage specializing in luxury homes, apartments, villas, and commercial spaces across Hyderabad, Bengaluru, Mumbai, and NCR.",
		ActivePage:   "home",
		CanonicalURL: "https://aarambha.in/",
		Company:      h.company,
		Data:         data,
	}

	h.render(w, r, "home.html", pageData)
}

func (h *PageHandler) Locations(w http.ResponseWriter, r *http.Request) {
	locations := h.catalogService.GetLocations()

	pageData := view.PageData{
		Title:        "Prime Metropolitan Real Estate Corridors & Micro-Market Guides",
		Description:  "Comprehensive capital value benchmarks, rental yields, and investment corridors across Hyderabad, Bengaluru, Mumbai, and Gurugram.",
		ActivePage:   "locations",
		CanonicalURL: "https://aarambha.in/locations",
		Company:      h.company,
		Data: struct {
			Locations []domain.LocationInsight
		}{
			Locations: locations,
		},
	}

	h.render(w, r, "locations.html", pageData)
}

func (h *PageHandler) Commercial(w http.ResponseWriter, r *http.Request) {
	properties := h.propertyService.GetCommercial()

	pageData := view.PageData{
		Title:        "Grade-A Commercial Real Estate & Pre-Leased Office Assets",
		Description:  "Explore Grade-A institutional office spaces, IT SEZ floor plates, and pre-leased commercial assets yielding 7.5% - 9.2% net returns.",
		ActivePage:   "commercial",
		CanonicalURL: "https://aarambha.in/commercial",
		Company:      h.company,
		Data: struct {
			Properties []domain.Property
		}{
			Properties: properties,
		},
	}

	h.render(w, r, "commercial.html", pageData)
}

func (h *PageHandler) Rent(w http.ResponseWriter, r *http.Request) {
	properties := h.propertyService.GetRentals()

	pageData := view.PageData{
		Title:        "Luxury Executive Rentals & Fully Furnished Homes",
		Description:  "Explore verified high-end rental apartments and villas in prime gated communities close to major tech clusters.",
		ActivePage:   "rent",
		CanonicalURL: "https://aarambha.in/rent",
		Company:      h.company,
		Data: struct {
			Properties []domain.Property
		}{
			Properties: properties,
		},
	}

	h.render(w, r, "rent.html", pageData)
}

func (h *PageHandler) Sell(w http.ResponseWriter, r *http.Request) {
	pageData := view.PageData{
		Title:        "Sell or Lease Your Luxury Property | Free Market Valuation",
		Description:  "Discreet, confidential seller representation for high-value villas, sky penthouses, and commercial parcels with institutional buyer access.",
		ActivePage:   "sell",
		CanonicalURL: "https://aarambha.in/sell",
		Company:      h.company,
		Data:         nil,
	}

	h.render(w, r, "sell.html", pageData)
}

func (h *PageHandler) NewProjects(w http.ResponseWriter, r *http.Request) {
	projects := h.catalogService.GetProjects()

	pageData := view.PageData{
		Title:        "RERA Approved Pre-Launch & Under-Construction Developer Projects",
		Description:  "Discover curated early-stage developer launches from Prestige, Godrej, and My Home Group with verified escrow accounts and construction milestones.",
		ActivePage:   "new-projects",
		CanonicalURL: "https://aarambha.in/new-projects",
		Company:      h.company,
		Data: struct {
			Projects []domain.Project
		}{
			Projects: projects,
		},
	}

	h.render(w, r, "new_projects.html", pageData)
}

func (h *PageHandler) Agents(w http.ResponseWriter, r *http.Request) {
	agents := h.catalogService.GetAgents()

	pageData := view.PageData{
		Title:        "Licensed Advisory Leaders & Verified Real Estate Brokers",
		Description:  "Meet our senior property consultants across Hyderabad, Bengaluru, and Mumbai with verified RERA credentials and hyper-local expertise.",
		ActivePage:   "agents",
		CanonicalURL: "https://aarambha.in/agents",
		Company:      h.company,
		Data: struct {
			Agents []domain.Agent
		}{
			Agents: agents,
		},
	}

	h.render(w, r, "agents.html", pageData)
}

func (h *PageHandler) Insights(w http.ResponseWriter, r *http.Request) {
	insights := h.catalogService.GetInsights()

	pageData := view.PageData{
		Title:        "Indian Real Estate Insights, RERA Checklists & Tax Guides",
		Description:  "Expert regulatory insights, home loan tax deduction guides, NRI property buying regulations, and quarterly corridor absorption reports.",
		ActivePage:   "insights",
		CanonicalURL: "https://aarambha.in/insights",
		Company:      h.company,
		Data: struct {
			Insights []domain.InsightArticle
		}{
			Insights: insights,
		},
	}

	h.render(w, r, "insights.html", pageData)
}

func (h *PageHandler) About(w http.ResponseWriter, r *http.Request) {
	pageData := view.PageData{
		Title:        "About Aarambha Agency | Ethical Fiduciary Advisory",
		Description:  "Founded in 2025, Aarambha Agency was built to restore trust and transparency to Indian property transactions with strict 40-point legal title verification.",
		ActivePage:   "about",
		CanonicalURL: "https://aarambha.in/about",
		Company:      h.company,
		Data:         nil,
	}

	h.render(w, r, "about.html", pageData)
}

func (h *PageHandler) Contact(w http.ResponseWriter, r *http.Request) {
	pageData := view.PageData{
		Title:        "Contact Aarambha Agency Advisory Desks | Hyderabad, Bengaluru, Mumbai",
		Description:  "Get in touch with our licensed real estate advisors or visit our regional headquarters in Hitec City, Indiranagar, or BKC.",
		ActivePage:   "contact",
		CanonicalURL: "https://aarambha.in/contact",
		Company:      h.company,
		Data:         nil,
	}

	h.render(w, r, "contact.html", pageData)
}
