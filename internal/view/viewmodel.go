package view

import (
	"context"

	"aarambha-agency/internal/domain"
)

type userContextKeyType struct{}

var UserContextKey = userContextKeyType{}

// WithUser adds a domain.User to the request context.
func WithUser(ctx context.Context, u *domain.User) context.Context {
	return context.WithValue(ctx, UserContextKey, u)
}

// UserFromContext extracts the authenticated domain.User from request context.
func UserFromContext(ctx context.Context) *domain.User {
	if ctx == nil {
		return nil
	}
	if u, ok := ctx.Value(UserContextKey).(*domain.User); ok {
		return u
	}
	return nil
}

type PageData struct {
	Title        string
	Description  string
	ActivePage   string
	CanonicalURL string
	Robots       string
	OGType       string
	OGImage      string
	SiteURL      string
	Company      domain.CompanyInfo
	User         *domain.User
	Data         interface{}
	IsHTMX       bool
}
