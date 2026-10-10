package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"aarambha-agency/internal/domain"
	"aarambha-agency/internal/service"
	"aarambha-agency/internal/view"
)

type AdminHandler struct {
	propertyService *service.PropertyService
	adminService    *service.AdminService
	catalogService  *service.CatalogService
	viewEngine      *view.Engine
	company         domain.CompanyInfo
}

func NewAdminHandler(
	pService *service.PropertyService,
	aService *service.AdminService,
	cService *service.CatalogService,
	vEngine *view.Engine,
	company domain.CompanyInfo,
) *AdminHandler {
	return &AdminHandler{
		propertyService: pService,
		adminService:    aService,
		catalogService:  cService,
		viewEngine:      vEngine,
		company:         company,
	}
}

func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	stats := h.adminService.GetStats()
	properties := h.propertyService.ListProperties(domain.PropertyFilter{})
	leads := h.adminService.GetRecentLeads()
	sellRequests := h.adminService.GetSellRequests()
	agents := h.catalogService.GetAgents()

	pageData := view.PageData{
		Title:       "Admin Portal - Real Estate Management",
		Description: "Aarambha Agency Internal Inventory and Advisory Management System",
		ActivePage:  "admin-dashboard",
		Company:     h.company,
		User:        view.UserFromContext(r.Context()),
		Data: map[string]interface{}{
			"Stats":        stats,
			"Properties":   properties,
			"Leads":        leads,
			"SellRequests": sellRequests,
			"Agents":       agents,
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.viewEngine.RenderAdminPage(w, "dashboard.html", pageData); err != nil {
		log.Printf("Error rendering admin dashboard: %v", err)
		http.Error(w, "Failed to render dashboard", http.StatusInternalServerError)
	}
}

func (h *AdminHandler) NewProperty(w http.ResponseWriter, r *http.Request) {
	pageData := view.PageData{
		Title:       "Add New Property Listing",
		Description: "Add a new luxury real estate listing to the catalog",
		ActivePage:  "admin-new-property",
		Company:     h.company,
		User:        view.UserFromContext(r.Context()),
		Data: map[string]interface{}{
			"IsEdit": false,
			"Property": &domain.Property{
				City:        "Hyderabad",
				State:       "Telangana",
				Transaction: "Sale",
				Type:        "Apartment",
				Facing:      "East",
				Possession:  "Ready to Move",
				Furnishing:  "Semi-Furnished",
				ReraStatus:  "RERA Approved",
				Verified:    true,
				Bedrooms:    3,
				Bathrooms:   3,
			},
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.viewEngine.RenderAdminPage(w, "property_form.html", pageData); err != nil {
		log.Printf("Error rendering new property form: %v", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
	}
}

func (h *AdminHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	prop := h.parsePropertyFromForm(r)
	u := view.UserFromContext(r.Context())
	if u != nil {
		prop.CreatedBy = &u.ID
		prop.UpdatedBy = &u.ID
	}

	if err := h.propertyService.Create(&prop); err != nil {
		pageData := view.PageData{
			Title:      "Add New Property Listing",
			ActivePage: "admin-new-property",
			Company:    h.company,
			User:       view.UserFromContext(r.Context()),
			Data: map[string]interface{}{
				"IsEdit":   false,
				"Property": &prop,
				"Error":    err.Error(),
			},
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = h.viewEngine.RenderAdminPage(w, "property_form.html", pageData)
		return
	}

	http.Redirect(w, r, "/admin?msg=created", http.StatusSeeOther)
}

func (h *AdminHandler) EditProperty(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	prop, err := h.propertyService.GetByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	pageData := view.PageData{
		Title:       "Edit Property: " + prop.Title,
		Description: "Modify existing listing details",
		ActivePage:  "admin-dashboard",
		Company:     h.company,
		User:        view.UserFromContext(r.Context()),
		Data: map[string]interface{}{
			"IsEdit":   true,
			"Property": prop,
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.viewEngine.RenderAdminPage(w, "property_form.html", pageData); err != nil {
		log.Printf("Error rendering edit property form: %v", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
	}
}

func (h *AdminHandler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing, err := h.propertyService.GetByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	updated := h.parsePropertyFromForm(r)
	updated.ID = existing.ID
	if updated.Slug == "" {
		updated.Slug = existing.Slug
	}
	if len(updated.Gallery) == 0 {
		updated.Gallery = existing.Gallery
	}
	if len(updated.Nearby) == 0 {
		updated.Nearby = existing.Nearby
	}
	if updated.AgentID == "" {
		updated.AgentID = existing.AgentID
	}
	updated.CreatedBy = existing.CreatedBy
	u := view.UserFromContext(r.Context())
	if u != nil {
		updated.UpdatedBy = &u.ID
	}

	if err := h.propertyService.Update(&updated); err != nil {
		pageData := view.PageData{
			Title:      "Edit Property: " + updated.Title,
			ActivePage: "admin-dashboard",
			Company:    h.company,
			User:       view.UserFromContext(r.Context()),
			Data: map[string]interface{}{
				"IsEdit":   true,
				"Property": &updated,
				"Error":    err.Error(),
			},
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = h.viewEngine.RenderAdminPage(w, "property_form.html", pageData)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *AdminHandler) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.propertyService.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if view.IsHTMX(r) {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *AdminHandler) ToggleFeatured(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := h.propertyService.ToggleFeatured(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prop, err := h.propertyService.GetByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if view.IsHTMX(r) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := h.viewEngine.RenderPartial(w, "admin_featured_badge", prop); err != nil {
			log.Printf("Error rendering featured badge partial: %v", err)
			http.Error(w, "Render error", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *AdminHandler) parsePropertyFromForm(r *http.Request) domain.Property {
	price, _ := strconv.ParseInt(strings.TrimSpace(r.FormValue("price")), 10, 64)
	bedrooms, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("bedrooms")))
	bathrooms, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("bathrooms")))
	area, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("area")))
	carpetArea, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("carpetArea")))

	var amenities []string
	rawAmenities := r.FormValue("amenities")
	if rawAmenities != "" {
		for _, a := range strings.Split(rawAmenities, ",") {
			cleaned := strings.TrimSpace(a)
			if cleaned != "" {
				amenities = append(amenities, cleaned)
			}
		}
	}

	return domain.Property{
		ID:           strings.TrimSpace(r.FormValue("id")),
		Title:        strings.TrimSpace(r.FormValue("title")),
		Transaction:  strings.TrimSpace(r.FormValue("transaction")),
		Type:         strings.TrimSpace(r.FormValue("type")),
		SubType:      strings.TrimSpace(r.FormValue("subType")),
		Developer:    strings.TrimSpace(r.FormValue("developer")),
		City:         strings.TrimSpace(r.FormValue("city")),
		Location:     strings.TrimSpace(r.FormValue("location")),
		State:        strings.TrimSpace(r.FormValue("state")),
		Facing:       strings.TrimSpace(r.FormValue("facing")),
		Floor:        strings.TrimSpace(r.FormValue("floor")),
		Parking:      strings.TrimSpace(r.FormValue("parking")),
		Price:        price,
		PriceDisplay: strings.TrimSpace(r.FormValue("priceDisplay")),
		Deposit:      strings.TrimSpace(r.FormValue("deposit")),
		Bedrooms:     bedrooms,
		Bathrooms:    bathrooms,
		Area:         area,
		CarpetArea:   carpetArea,
		AreaUnit:     "sq.ft.",
		Furnishing:   strings.TrimSpace(r.FormValue("furnishing")),
		Possession:   strings.TrimSpace(r.FormValue("possession")),
		ReraStatus:   strings.TrimSpace(r.FormValue("reraStatus")),
		ReraNumber:   strings.TrimSpace(r.FormValue("reraNumber")),
		Image:        strings.TrimSpace(r.FormValue("image")),
		Description:  strings.TrimSpace(r.FormValue("description")),
		Amenities:    amenities,
		Featured:     r.FormValue("featured") == "true",
		Verified:     r.FormValue("verified") == "true",
		IsNewLaunch:  r.FormValue("isNewLaunch") == "true",
	}
}

func (h *AdminHandler) AgentsList(w http.ResponseWriter, r *http.Request) {
	agents := h.catalogService.GetAgents()

	pageData := view.PageData{
		Title:       "Advisory Panel & Agents Management",
		Description: "Licensed senior real estate fiduciaries and advisors",
		ActivePage:  "admin-agents",
		Company:     h.company,
		User:        view.UserFromContext(r.Context()),
		Data: map[string]interface{}{
			"Agents": agents,
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.viewEngine.RenderAdminPage(w, "agents.html", pageData); err != nil {
		log.Printf("Error rendering admin agents page: %v", err)
		http.Error(w, "Failed to render agents", http.StatusInternalServerError)
	}
}

func (h *AdminHandler) NewAgent(w http.ResponseWriter, r *http.Request) {
	pageData := view.PageData{
		Title:       "Register Licensed Advisor",
		Description: "Add a new real estate fiduciary consultant",
		ActivePage:  "admin-agents",
		Company:     h.company,
		User:        view.UserFromContext(r.Context()),
		Data: map[string]interface{}{
			"IsEdit": false,
			"Agent": &domain.Agent{
				City:           "Hyderabad",
				ActiveListings: 3,
				Experience:     "10+ Years",
				Photo:          "https://images.unsplash.com/photo-1560250097-0b93528c311a?auto=format&fit=crop&w=600&q=80",
			},
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.viewEngine.RenderAdminPage(w, "agent_form.html", pageData); err != nil {
		log.Printf("Error rendering new agent form: %v", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
	}
}

func (h *AdminHandler) CreateAgent(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	agent := h.parseAgentFromForm(r)
	u := view.UserFromContext(r.Context())
	if u != nil {
		agent.CreatedBy = &u.ID
		agent.UpdatedBy = &u.ID
	}

	if err := h.catalogService.CreateAgent(&agent); err != nil {
		pageData := view.PageData{
			Title:      "Register Licensed Advisor",
			ActivePage: "admin-agents",
			Company:    h.company,
			User:       view.UserFromContext(r.Context()),
			Data: map[string]interface{}{
				"IsEdit": false,
				"Agent":  &agent,
				"Error":  err.Error(),
			},
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = h.viewEngine.RenderAdminPage(w, "agent_form.html", pageData)
		return
	}

	http.Redirect(w, r, "/admin/agents?msg=created", http.StatusSeeOther)
}

func (h *AdminHandler) EditAgent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	agent, err := h.catalogService.GetAgentByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	pageData := view.PageData{
		Title:       "Edit Advisor: " + agent.Name,
		Description: "Modify fiduciary credentials and contact details",
		ActivePage:  "admin-agents",
		Company:     h.company,
		User:        view.UserFromContext(r.Context()),
		Data: map[string]interface{}{
			"IsEdit": true,
			"Agent":  agent,
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.viewEngine.RenderAdminPage(w, "agent_form.html", pageData); err != nil {
		log.Printf("Error rendering edit agent form: %v", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
	}
}

func (h *AdminHandler) UpdateAgent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing, err := h.catalogService.GetAgentByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	updated := h.parseAgentFromForm(r)
	updated.ID = existing.ID
	updated.CreatedBy = existing.CreatedBy
	u := view.UserFromContext(r.Context())
	if u != nil {
		updated.UpdatedBy = &u.ID
	}

	if err := h.catalogService.UpdateAgent(&updated); err != nil {
		pageData := view.PageData{
			Title:      "Edit Advisor: " + updated.Name,
			ActivePage: "admin-agents",
			Company:    h.company,
			User:       view.UserFromContext(r.Context()),
			Data: map[string]interface{}{
				"IsEdit": true,
				"Agent":  &updated,
				"Error":  err.Error(),
			},
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = h.viewEngine.RenderAdminPage(w, "agent_form.html", pageData)
		return
	}

	http.Redirect(w, r, "/admin/agents?msg=updated", http.StatusSeeOther)
}

func (h *AdminHandler) DeleteAgent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.catalogService.DeleteAgent(id); err != nil {
		log.Printf("Error deleting agent %s: %v", id, err)
	}
	http.Redirect(w, r, "/admin/agents?msg=deleted", http.StatusSeeOther)
}

func (h *AdminHandler) parseAgentFromForm(r *http.Request) domain.Agent {
	activeListings, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("activeListings")))

	var areas []string
	rawAreas := r.FormValue("areasServed")
	if rawAreas != "" {
		for _, a := range strings.Split(rawAreas, ",") {
			cleaned := strings.TrimSpace(a)
			if cleaned != "" {
				areas = append(areas, cleaned)
			}
		}
	}

	var languages []string
	rawLanguages := r.FormValue("languages")
	if rawLanguages != "" {
		for _, l := range strings.Split(rawLanguages, ",") {
			cleaned := strings.TrimSpace(l)
			if cleaned != "" {
				languages = append(languages, cleaned)
			}
		}
	}

	photo := strings.TrimSpace(r.FormValue("photo"))
	if photo == "" {
		photo = "https://images.unsplash.com/photo-1560250097-0b93528c311a?auto=format&fit=crop&w=600&q=80"
	}

	return domain.Agent{
		ID:             strings.TrimSpace(r.FormValue("id")),
		Name:           strings.TrimSpace(r.FormValue("name")),
		Role:           strings.TrimSpace(r.FormValue("role")),
		City:           strings.TrimSpace(r.FormValue("city")),
		AreasServed:    areas,
		Specialization: strings.TrimSpace(r.FormValue("specialization")),
		Experience:     strings.TrimSpace(r.FormValue("experience")),
		Languages:      languages,
		ReraID:         strings.TrimSpace(r.FormValue("reraId")),
		Photo:          photo,
		Phone:          strings.TrimSpace(r.FormValue("phone")),
		WhatsApp:       strings.TrimSpace(r.FormValue("whatsapp")),
		Email:          strings.TrimSpace(r.FormValue("email")),
		ActiveListings: activeListings,
		Bio:            strings.TrimSpace(r.FormValue("bio")),
	}
}

func (h *AdminHandler) UpdateSellRequestStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		id = r.FormValue("id")
	}

	status := strings.ToLower(strings.TrimSpace(r.FormValue("status")))
	if status == "" {
		status = strings.ToLower(strings.TrimSpace(r.URL.Query().Get("status")))
	}

	if status != "approved" && status != "rejected" && status != "pending" {
		http.Error(w, "Invalid status: must be approved, rejected, or pending", http.StatusBadRequest)
		return
	}

	if err := h.adminService.UpdateLeadStatus(id, status); err != nil {
		log.Printf("Error updating sell request status %s to %s: %v", id, status, err)
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	// If request came from HTMX, return the updated table row
	if r.Header.Get("HX-Request") != "" {
		lead, err := h.adminService.GetLeadByID(id)
		if err == nil && lead != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			h.renderSellRequestRow(w, *lead)
			return
		}
	}

	http.Redirect(w, r, "/admin?tab=sell-requests&msg=status_updated", http.StatusSeeOther)
}

func (h *AdminHandler) renderSellRequestRow(w http.ResponseWriter, l domain.Lead) {
	dateStr := "Recent"
	if !l.CreatedAt.IsZero() {
		dateStr = l.CreatedAt.Format("02 Jan 2006, 15:04")
	}

	statusBadge := ""
	actionButtons := ""

	cleanTitle := url.QueryEscape(l.PropertyTitle)
	createListingURL := fmt.Sprintf("/admin/properties/new?title=%s", cleanTitle)

	switch l.Status {
	case "approved":
		statusBadge = `<span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-bold bg-emerald-100 dark:bg-emerald-950/70 text-emerald-800 dark:text-emerald-300 border border-emerald-300/40 shadow-sm">
			<svg class="w-3.5 h-3.5 text-emerald-600 dark:text-emerald-400" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>
			Approved
		</span>`
		actionButtons = fmt.Sprintf(`
			<div class="flex items-center justify-end gap-2">
				<a href="%s" class="px-2.5 py-1.5 bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-navy-950 rounded-lg text-xs font-bold transition-all shadow-sm flex items-center gap-1 whitespace-nowrap">
					<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4"/></svg>
					<span>Create Listing</span>
				</a>
				<button 
					type="button" 
					hx-post="/admin/sell-requests/%s/status?status=rejected" 
					hx-target="#sell-request-row-%s" 
					hx-swap="outerHTML" 
					class="px-2.5 py-1.5 bg-rose-50 dark:bg-rose-950/50 hover:bg-rose-100 dark:hover:bg-rose-900/60 text-rose-700 dark:text-rose-300 border border-rose-200 dark:border-rose-800/60 rounded-lg text-xs font-semibold transition-all cursor-pointer whitespace-nowrap"
					title="Change decision to Rejected">
					Reject
				</button>
			</div>
		`, createListingURL, l.ID, l.ID)
	case "rejected":
		statusBadge = `<span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-bold bg-rose-100 dark:bg-rose-950/70 text-rose-800 dark:text-rose-300 border border-rose-300/40 shadow-sm">
			<svg class="w-3.5 h-3.5 text-rose-600 dark:text-rose-400" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/></svg>
			Rejected
		</span>`
		actionButtons = fmt.Sprintf(`
			<div class="flex items-center justify-end gap-2">
				<button 
					type="button" 
					hx-post="/admin/sell-requests/%s/status?status=approved" 
					hx-target="#sell-request-row-%s" 
					hx-swap="outerHTML" 
					class="px-3 py-1.5 bg-emerald-50 dark:bg-emerald-950/50 hover:bg-emerald-100 dark:hover:bg-emerald-900/60 text-emerald-700 dark:text-emerald-300 border border-emerald-300 dark:border-emerald-800/60 rounded-lg text-xs font-bold transition-all cursor-pointer flex items-center gap-1 whitespace-nowrap">
					<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/></svg>
					<span>Approve</span>
				</button>
			</div>
		`, l.ID, l.ID)
	default: // pending
		statusBadge = `<span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-bold bg-amber-100 dark:bg-amber-950/70 text-amber-800 dark:text-amber-300 border border-amber-300/40 shadow-sm animate-pulse">
			<svg class="w-3.5 h-3.5 text-amber-600 dark:text-amber-400" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd"/></svg>
			Pending Review
		</span>`
		actionButtons = fmt.Sprintf(`
			<div class="flex items-center justify-end gap-2">
				<button 
					type="button" 
					hx-post="/admin/sell-requests/%s/status?status=approved" 
					hx-target="#sell-request-row-%s" 
					hx-swap="outerHTML" 
					class="px-3 py-1.5 bg-gradient-to-r from-emerald-600 to-emerald-700 hover:from-emerald-500 hover:to-emerald-600 text-white rounded-lg text-xs font-bold shadow-sm transition-all active:scale-95 cursor-pointer flex items-center gap-1 whitespace-nowrap">
					<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/></svg>
					<span>Approve</span>
				</button>
				<button 
					type="button" 
					hx-post="/admin/sell-requests/%s/status?status=rejected" 
					hx-target="#sell-request-row-%s" 
					hx-swap="outerHTML" 
					class="px-2.5 py-1.5 bg-rose-50 dark:bg-rose-950/50 hover:bg-rose-100 dark:hover:bg-rose-900/60 text-rose-700 dark:text-rose-300 border border-rose-200 dark:border-rose-800/60 rounded-lg text-xs font-semibold transition-all active:scale-95 cursor-pointer whitespace-nowrap">
					Reject
				</button>
			</div>
		`, l.ID, l.ID, l.ID, l.ID)
	}

	emailSnippet := ""
	if l.Email != "" {
		emailSnippet = fmt.Sprintf(`<span>•</span><a href="mailto:%s" class="hover:text-amber-600 dark:hover:text-amber-400">%s</a>`, l.Email, l.Email)
	}

	propSnippet := l.PropertyTitle
	if propSnippet == "" {
		propSnippet = "Unspecified Asset"
	}

	msgSnippet := l.Message
	if msgSnippet == "" {
		msgSnippet = `<span class="text-slate-400 italic">No additional notes</span>`
	}

	fmt.Fprintf(w, `
		<tr id="sell-request-row-%s" class="hover:bg-slate-50/80 dark:hover:bg-navy-800/50 transition-colors">
			<td class="py-3.5 px-4 sm:px-6 text-xs text-slate-500 dark:text-slate-400 whitespace-nowrap">
				%s
			</td>
			<td class="py-3.5 px-4">
				<div class="font-bold text-navy-950 dark:text-white">%s</div>
				<div class="text-xs text-slate-500 dark:text-slate-400 flex items-center gap-2 mt-0.5">
					<a href="tel:%s" class="hover:text-amber-600 dark:hover:text-amber-400 font-mono font-medium">%s</a>
					%s
				</div>
			</td>
			<td class="py-3.5 px-4 text-xs font-medium text-slate-800 dark:text-slate-200 max-w-xs">
				<span class="font-semibold text-navy-950 dark:text-white block">%s</span>
				<p class="text-[11px] text-slate-500 dark:text-slate-400 mt-0.5 line-clamp-2">%s</p>
			</td>
			<td class="py-3.5 px-4 whitespace-nowrap">
				%s
			</td>
			<td class="py-3.5 px-4 text-right whitespace-nowrap">
				%s
			</td>
		</tr>
	`, l.ID, dateStr, l.Name, l.Phone, l.Phone, emailSnippet, propSnippet, msgSnippet, statusBadge, actionButtons)
}
