package handler

import (
	"fmt"
	"net/http"
	"strings"

	"aarambha-agency/internal/domain"
	"aarambha-agency/internal/service"
	"aarambha-agency/internal/view"
)

type LeadHandler struct {
	leadService *service.LeadService
}

func NewLeadHandler(service *service.LeadService) *LeadHandler {
	return &LeadHandler{leadService: service}
}

func (h *LeadHandler) Submit(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	lead := domain.Lead{
		Name:          strings.TrimSpace(r.FormValue("name")),
		Phone:         strings.TrimSpace(r.FormValue("phone")),
		Email:         strings.TrimSpace(r.FormValue("email")),
		PropertyID:    strings.TrimSpace(r.FormValue("propertyId")),
		PropertyTitle: strings.TrimSpace(r.FormValue("propertyTitle")),
		Type:          strings.TrimSpace(r.FormValue("type")),
		Status:        "pending",
		PreferredDate: strings.TrimSpace(r.FormValue("preferredDate")),
		Message:       strings.TrimSpace(r.FormValue("message")),
	}

	u := view.UserFromContext(r.Context())
	if u != nil {
		lead.UserID = &u.ID
		lead.CreatedBy = &u.ID
		lead.UpdatedBy = &u.ID
		if lead.Email == "" {
			lead.Email = u.Email
		}
	}

	result, err := h.leadService.SubmitLead(lead)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
			<div class="p-4 bg-red-50 text-red-700 border border-red-200 rounded-xl text-xs space-y-2">
				<p class="font-bold">%s</p>
				<button onclick="window.location.reload()" class="text-xs font-semibold underline text-red-800">Try Again</button>
			</div>
		`, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	heading := fmt.Sprintf("Inquiry Confirmed, %s!", result.Lead.Name)
	message := "Our senior property advisor has received your details and will connect with verified documentation within 2 business hours."
	if result.Lead.Type == "sell_request" {
		heading = fmt.Sprintf("Sell Request Submitted, %s!", result.Lead.Name)
		message = "Our acquisitions & asset advisory team has received your property details. An administrator will review your sell request and reach out shortly."
	}

	fmt.Fprintf(w, `
		<div class="p-6 bg-emerald-50 border border-emerald-200 rounded-2xl text-center space-y-3 animate-fade-in">
			<div class="w-12 h-12 rounded-full bg-emerald-100 text-emerald-600 flex items-center justify-center mx-auto">
				<svg class="w-6 h-6" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>
			</div>
			<h4 class="text-base font-bold text-navy-900 font-serif">%s</h4>
			<p class="text-xs text-slate-600 leading-relaxed max-w-sm mx-auto">
				%s
			</p>
			<div class="pt-2">
				<a href="%s" target="_blank" rel="noopener" class="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-xs uppercase tracking-wider transition-all shadow-md">
					<span>Chat on WhatsApp Now</span>
				</a>
			</div>
		</div>
	`, heading, message, result.WhatsAppURL)
}
