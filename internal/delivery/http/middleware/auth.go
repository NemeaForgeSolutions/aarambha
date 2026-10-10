package middleware

import (
	"net/http"
	"net/url"

	"aarambha-agency/internal/service"
	"aarambha-agency/internal/view"
)

// Authenticate extracts session cookie and injects *domain.User into context if valid.
func Authenticate(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := authService.GetUserFromRequest(r)
			if err == nil && user != nil {
				ctx := view.WithUser(r.Context(), user)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAuth enforces that a user must be logged in.
func RequireAuth(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := view.UserFromContext(r.Context())
			if user == nil {
				// Try fetching directly from cookie in case Authenticate wasn't executed
				var err error
				user, err = authService.GetUserFromRequest(r)
				if err == nil && user != nil {
					ctx := view.WithUser(r.Context(), user)
					r = r.WithContext(ctx)
				}
			}

			if user == nil {
				target := r.URL.RequestURI()
				loginURL := "/login"
				if target != "" && target != "/" {
					loginURL = "/login?redirect=" + url.QueryEscape(target)
				}

				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", loginURL)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				http.Redirect(w, r, loginURL, http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin enforces that a user must be logged in with the 'admin' role.
func RequireAdmin(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := view.UserFromContext(r.Context())
			if user == nil {
				var err error
				user, err = authService.GetUserFromRequest(r)
				if err == nil && user != nil {
					ctx := view.WithUser(r.Context(), user)
					r = r.WithContext(ctx)
				}
			}

			if user == nil {
				target := r.URL.RequestURI()
				loginURL := "/login?redirect=" + url.QueryEscape(target)
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", loginURL)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				http.Redirect(w, r, loginURL, http.StatusSeeOther)
				return
			}

			if !user.IsAdmin() {
				if r.Header.Get("HX-Request") == "true" {
					http.Error(w, "Forbidden: Admin privileges required", http.StatusForbidden)
					return
				}
				http.Redirect(w, r, "/profile?error=Admin+privileges+required", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
