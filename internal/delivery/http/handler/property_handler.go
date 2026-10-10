package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"aarambha-agency/internal/domain"
	"aarambha-agency/internal/service"
	"aarambha-agency/internal/view"
)

type PropertyHandler struct {
	propertyService *service.PropertyService
	engine          *view.Engine
	company         domain.CompanyInfo
}

func NewPropertyHandler(
	pService *service.PropertyService,
	engine *view.Engine,
	company domain.CompanyInfo,
) *PropertyHandler {
	return &PropertyHandler{
		propertyService: pService,
		engine:          engine,
		company:         company,
	}
}

func parseFilter(r *http.Request) domain.PropertyFilter {
	q := r.URL.Query()
	bhk, _ := strconv.Atoi(q.Get("bhk"))
	minPrice, _ := strconv.ParseInt(q.Get("minPrice"), 10, 64)
	maxPrice, _ := strconv.ParseInt(q.Get("maxPrice"), 10, 64)

	return domain.PropertyFilter{
		Transaction: q.Get("transaction"),
		City:        q.Get("city"),
		Type:        q.Get("type"),
		Bedrooms:    bhk,
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
		Possession:  q.Get("possession"),
		Furnishing:  q.Get("furnishing"),
		ReraOnly:    q.Get("rera") == "true" || q.Get("rera") == "1",
		Query:       strings.TrimSpace(q.Get("q")),
		Sort:        q.Get("sort"),
	}
}

func (h *PropertyHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := parseFilter(r)
	properties := h.propertyService.ListProperties(filter)

	pageData := view.PageData{
		Title:        "Discover Luxury Properties & Commercial Real Estate",
		Description:  "Explore verified luxury apartments, gated villas, penthouses, and Grade-A commercial office spaces across Hyderabad, Bengaluru, Mumbai, and NCR.",
		ActivePage:   "buy",
		CanonicalURL: "https://aarambha.in/properties",
		Company:      h.company,
		Data: struct {
			Properties []domain.Property
			Filter     domain.PropertyFilter
		}{
			Properties: properties,
			Filter:     filter,
		},
		IsHTMX: view.IsHTMX(r),
		User:   view.UserFromContext(r.Context()),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.engine.RenderPage(w, "properties.html", pageData); err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
	}
}

func (h *PropertyHandler) Search(w http.ResponseWriter, r *http.Request) {
	filter := parseFilter(r)
	properties := h.propertyService.ListProperties(filter)

	if len(properties) == 0 {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
			<div class="col-span-full py-16 text-center bg-white rounded-2xl border border-slate-200 p-8 space-y-4">
				<div class="w-16 h-16 rounded-full bg-amber-500/10 text-amber-600 flex items-center justify-center mx-auto text-2xl font-bold">
					!
				</div>
				<h3 class="text-xl font-bold text-navy-900 font-serif">No Matching Properties Found</h3>
				<p class="text-slate-500 text-sm max-w-md mx-auto">
					We couldn't find any properties matching your exact filter criteria. Try broadening your city or price filters.
				</p>
				<a href="/properties" class="inline-block px-6 py-2.5 rounded-xl bg-navy-900 text-white font-bold text-xs uppercase tracking-wider">
					Clear All Filters
				</a>
			</div>
		`)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, prop := range properties {
		if err := h.engine.RenderPartial(w, "property_card", prop); err != nil {
			http.Error(w, fmt.Sprintf("Partial error: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func (h *PropertyHandler) Featured(w http.ResponseWriter, r *http.Request) {
	propType := r.URL.Query().Get("type")
	featured := h.propertyService.GetFeatured(propType, 6)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if len(featured) == 0 {
		fmt.Fprintf(w, `<div class="col-span-full py-16 text-center text-ivory-600 dark:text-ivory-400">
			<svg class="w-10 h-10 mx-auto mb-3 text-champagne-500/60" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/></svg>
			<p class="text-sm font-semibold text-ivory-900 dark:text-ivory-200">No properties currently featured in this category</p>
			<p class="text-xs text-ivory-500 mt-1">Please check back shortly or explore our full collection.</p>
		</div>`)
		return
	}

	for _, prop := range featured {
		if err := h.engine.RenderPartial(w, "property_card", prop); err != nil {
			http.Error(w, fmt.Sprintf("Partial error: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func (h *PropertyHandler) Detail(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		slug = r.URL.Query().Get("id")
	}

	prop, agent, similar, err := h.propertyService.GetDetail(slug)
	if err != nil {
		h.engine.RenderNotFound(w, r, h.company)
		return
	}

	pageData := view.PageData{
		Title:        prop.Title,
		Description:  prop.Description,
		ActivePage:   "buy",
		CanonicalURL: "https://aarambha.in/properties/" + prop.Slug,
		OGType:       "article",
		OGImage:      prop.Image,
		Company:      h.company,
		Data: struct {
			Property *domain.Property
			Agent    *domain.Agent
			Similar  []domain.Property
		}{
			Property: prop,
			Agent:    agent,
			Similar:  similar,
		},
		IsHTMX: view.IsHTMX(r),
		User:   view.UserFromContext(r.Context()),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.engine.RenderPage(w, "property_detail.html", pageData); err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
	}
}
