package view

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"

	"aarambha-agency/internal/domain"
	"aarambha-agency/web"
)

type Engine struct {
	pages      map[string]*template.Template
	adminPages map[string]*template.Template
	partials   *template.Template
}

var funcMap = template.FuncMap{
	"formatCurrency": func(amount int64) string {
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
	},
	"formatNumber": func(n int) string {
		return fmt.Sprintf("%d", n)
	},
	"lower":     strings.ToLower,
	"upper":     strings.ToUpper,
	"contains":  strings.Contains,
	"hasPrefix": strings.HasPrefix,
	"split":     strings.Split,
	"join":      strings.Join,
	"safeHTML": func(s string) template.HTML {
		return template.HTML(s)
	},
	"whatsappURL": func(text string) string {
		return fmt.Sprintf("https://wa.me/917386985852?text=%s", url.QueryEscape(text))
	},
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
	"sliceProperties": func(list []domain.Property, start, end int) []domain.Property {
		if start < 0 {
			start = 0
		}
		if start >= len(list) {
			return nil
		}
		if end > len(list) || end <= 0 {
			end = len(list)
		}
		return list[start:end]
	},
	"sliceProjects": func(list []domain.Project, start, end int) []domain.Project {
		if start < 0 {
			start = 0
		}
		if start >= len(list) {
			return nil
		}
		if end > len(list) || end <= 0 {
			end = len(list)
		}
		return list[start:end]
	},
	"sliceAgents": func(list []domain.Agent, start, end int) []domain.Agent {
		if start < 0 {
			start = 0
		}
		if start >= len(list) {
			return nil
		}
		if end > len(list) || end <= 0 {
			end = len(list)
		}
		return list[start:end]
	},
}

func NewEngine() (*Engine, error) {
	engine := &Engine{
		pages:      make(map[string]*template.Template),
		adminPages: make(map[string]*template.Template),
	}

	partials, err := template.New("partials").Funcs(funcMap).ParseFS(
		web.TemplateFS,
		"template/partials/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse partials: %w", err)
	}
	engine.partials = partials

	pageNames := []string{
		"home.html",
		"properties.html",
		"property_detail.html",
		"locations.html",
		"commercial.html",
		"rent.html",
		"sell.html",
		"new_projects.html",
		"agents.html",
		"insights.html",
		"about.html",
		"contact.html",
		"login.html",
		"profile.html",
		"404.html",
	}

	for _, name := range pageNames {
		tmpl, err := template.New("base.html").Funcs(funcMap).ParseFS(
			web.TemplateFS,
			"template/layouts/base.html",
			"template/partials/*.html",
			"template/pages/"+name,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", name, err)
		}
		engine.pages[name] = tmpl
	}

	adminPageNames := []string{
		"dashboard.html",
		"property_form.html",
		"agents.html",
		"agent_form.html",
	}

	for _, name := range adminPageNames {
		tmpl, err := template.New("admin_base.html").Funcs(funcMap).ParseFS(
			web.TemplateFS,
			"template/layouts/admin_base.html",
			"template/partials/*.html",
			"template/admin/"+name,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to parse admin template %s: %w", name, err)
		}
		engine.adminPages[name] = tmpl
	}

	return engine, nil
}

func (e *Engine) RenderPage(w io.Writer, pageName string, data PageData) error {
	tmpl, exists := e.pages[pageName]
	if !exists {
		return fmt.Errorf("page template %s not found", pageName)
	}
	if hw, ok := w.(http.ResponseWriter); ok {
		if hw.Header().Get("Content-Type") == "" {
			hw.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
	}
	return tmpl.ExecuteTemplate(w, "base.html", data)
}

// RenderNotFound writes a 404 status and renders the branded 404 page.
// It is safe to call from any handler that has access to the view engine.
func (e *Engine) RenderNotFound(w http.ResponseWriter, r *http.Request, company domain.CompanyInfo) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	data := PageData{
		Title:       "Page Not Found",
		Description: "The page you're looking for doesn't exist. Browse our luxury property portfolio or contact our advisory team.",
		Robots:      "noindex, nofollow",
		Company:     company,
		User:        UserFromContext(r.Context()),
	}
	if err := e.RenderPage(w, "404.html", data); err != nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	}
}

func (e *Engine) RenderAdminPage(w io.Writer, pageName string, data PageData) error {
	tmpl, exists := e.adminPages[pageName]
	if !exists {
		return fmt.Errorf("admin page template %s not found", pageName)
	}
	if hw, ok := w.(http.ResponseWriter); ok {
		if hw.Header().Get("Content-Type") == "" {
			hw.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
	}
	return tmpl.ExecuteTemplate(w, "admin_base.html", data)
}

func (e *Engine) RenderPartial(w io.Writer, partialName string, data interface{}) error {
	if hw, ok := w.(http.ResponseWriter); ok {
		if hw.Header().Get("Content-Type") == "" {
			hw.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
	}
	var buf bytes.Buffer
	if err := e.partials.ExecuteTemplate(&buf, partialName, data); err != nil {
		return err
	}
	_, err := w.Write(buf.Bytes())
	return err
}

func IsHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
