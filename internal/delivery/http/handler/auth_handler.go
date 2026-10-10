package handler

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"aarambha-agency/internal/domain"
	"aarambha-agency/internal/service"
	"aarambha-agency/internal/view"
)

type AuthHandler struct {
	authService  *service.AuthService
	viewEngine   *view.Engine
	company      domain.CompanyInfo
	isProduction bool
}

func NewAuthHandler(
	aService *service.AuthService,
	vEngine *view.Engine,
	company domain.CompanyInfo,
	isProduction bool,
) *AuthHandler {
	return &AuthHandler{
		authService:  aService,
		viewEngine:   vEngine,
		company:      company,
		isProduction: isProduction,
	}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// If already authenticated, redirect to appropriate landing page
	if user := view.UserFromContext(r.Context()); user != nil {
		if user.IsAdmin() {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	redirect := r.URL.Query().Get("redirect")
	errorMsg := r.URL.Query().Get("error")
	msg := r.URL.Query().Get("msg")
	var successMsg string
	if msg == "logged_out" {
		successMsg = "You have been securely signed out."
	} else if msg == "reg_restricted" {
		successMsg = "Registration completed. Note: Portal sign-in is currently reserved for administrators and registered agents only."
	}

	pageData := view.PageData{
		Title:       "Sign In to Advisory Account",
		Description: "Access Aarambha Agency luxury properties portfolio, client dashboards and advisory management",
		ActivePage:  "login",
		Company:     h.company,
		Data: map[string]interface{}{
			"Redirect":       redirect,
			"Error":          errorMsg,
			"SuccessMessage": successMsg,
			"Email":          "",
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.viewEngine.RenderPage(w, "login.html", pageData); err != nil {
		log.Printf("Error rendering login page: %v", err)
		http.Error(w, "Failed to render login page", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) LoginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	redirect := strings.TrimSpace(r.FormValue("redirect"))

	user, err := h.authService.Authenticate(email, password)
	if err != nil {
		errorMsg := "Invalid email or password. Please verify your credentials."
		if errors.Is(err, domain.ErrLoginRestricted) {
			errorMsg = "Portal login is currently restricted to authorized administrators and registered agents only."
		} else if errors.Is(err, domain.ErrUserInactive) {
			errorMsg = "Your account is currently inactive. Please contact administration."
		}

		pageData := view.PageData{
			Title:       "Sign In to Advisory Account",
			Description: "Access Aarambha Agency luxury properties portfolio, client dashboards and advisory management",
			ActivePage:  "login",
			Company:     h.company,
			Data: map[string]interface{}{
				"Redirect": redirect,
				"Error":    errorMsg,
				"Email":    email,
			},
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = h.viewEngine.RenderPage(w, "login.html", pageData)
		return
	}

	// Set session cookie
	h.authService.SetSessionCookie(w, user, h.isProduction)

	// Safe redirection check
	if redirect != "" && strings.HasPrefix(redirect, "/") && !strings.HasPrefix(redirect, "//") {
		// Non-admin attempting to redirect to /admin should go to /profile instead
		if strings.HasPrefix(redirect, "/admin") && !user.IsAdmin() {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, redirect, http.StatusSeeOther)
		return
	}

	if user.IsAdmin() {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.authService.ClearSessionCookie(w, h.isProduction)
	http.Redirect(w, r, "/login?msg=logged_out", http.StatusSeeOther)
}

func (h *AuthHandler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	user := view.UserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login?redirect=/profile", http.StatusSeeOther)
		return
	}

	var successMsg string
	if r.URL.Query().Get("msg") == "welcome" {
		successMsg = "Welcome to Aarambha Agency! Your account has been initialized."
	} else if r.URL.Query().Get("msg") == "updated" {
		successMsg = "Your profile details have been successfully updated."
	}

	errorMsg := r.URL.Query().Get("error")

	pageData := view.PageData{
		Title:       "Client Profile & Advisory Desk",
		Description: "Manage your advisory profile and verified luxury searches",
		ActivePage:  "profile",
		Company:     h.company,
		User:        user,
		Data: map[string]interface{}{
			"SuccessMessage": successMsg,
			"Error":          errorMsg,
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.viewEngine.RenderPage(w, "profile.html", pageData); err != nil {
		log.Printf("Error rendering profile page: %v", err)
		http.Error(w, "Failed to render profile page", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) ProfileUpdate(w http.ResponseWriter, r *http.Request) {
	user := view.UserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login?redirect=/profile", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	phone := strings.TrimSpace(r.FormValue("phone"))

	updatedUser, err := h.authService.UpdateProfile(user.ID, name, phone)
	if err != nil {
		pageData := view.PageData{
			Title:       "Client Profile & Advisory Desk",
			Description: "Manage your advisory profile and verified luxury searches",
			ActivePage:  "profile",
			Company:     h.company,
			User:        user,
			Data: map[string]interface{}{
				"Error": "Failed to update profile: " + err.Error(),
			},
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = h.viewEngine.RenderPage(w, "profile.html", pageData)
		return
	}

	// Update user in context and re-issue session cookie
	h.authService.SetSessionCookie(w, updatedUser, h.isProduction)
	http.Redirect(w, r, "/profile?msg=updated", http.StatusSeeOther)
}
