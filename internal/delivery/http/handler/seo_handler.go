package handler

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"aarambha-agency/internal/domain"
	"aarambha-agency/internal/service"
)

type SEOHandler struct {
	propertyService *service.PropertyService
	catalogService  *service.CatalogService
	baseURL         string
}

func NewSEOHandler(pService *service.PropertyService, cService *service.CatalogService, baseURL string) *SEOHandler {
	if baseURL == "" {
		baseURL = "https://aarambha.in"
	}
	return &SEOHandler{
		propertyService: pService,
		catalogService:  cService,
		baseURL:         baseURL,
	}
}

type urlEntry struct {
	XMLName    xml.Name `xml:"url"`
	Loc        string   `xml:"loc"`
	LastMod    string   `xml:"lastmod,omitempty"`
	ChangeFreq string   `xml:"changefreq,omitempty"`
	Priority   string   `xml:"priority,omitempty"`
}

type urlSet struct {
	XMLName xml.Name   `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
	URLs    []urlEntry `xml:"url"`
}

func (h *SEOHandler) SitemapXML(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format("2006-01-02")

	staticPages := []struct {
		path     string
		freq     string
		priority string
	}{
		{"/", "daily", "1.0"},
		{"/properties", "hourly", "0.9"},
		{"/locations", "weekly", "0.8"},
		{"/commercial", "daily", "0.8"},
		{"/rent", "daily", "0.8"},
		{"/sell", "monthly", "0.7"},
		{"/new-projects", "daily", "0.8"},
		{"/agents", "weekly", "0.7"},
		{"/insights", "weekly", "0.7"},
		{"/about", "monthly", "0.5"},
		{"/contact", "monthly", "0.6"},
	}

	urls := make([]urlEntry, 0, len(staticPages)+20)
	for _, p := range staticPages {
		urls = append(urls, urlEntry{
			Loc:        h.baseURL + p.path,
			LastMod:    now,
			ChangeFreq: p.freq,
			Priority:   p.priority,
		})
	}

	properties := h.propertyService.ListProperties(domain.PropertyFilter{})
	for _, prop := range properties {
		urls = append(urls, urlEntry{
			Loc:        fmt.Sprintf("%s/properties/%s", h.baseURL, prop.Slug),
			LastMod:    now,
			ChangeFreq: "weekly",
			Priority:   "0.9",
		})
	}

	doc := urlSet{URLs: urls}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	_ = enc.Encode(doc)
}

func (h *SEOHandler) RobotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `# Aarambha Agency Robots Rules
User-agent: *
Allow: /
Disallow: /admin/
Disallow: /api/

Sitemap: %s/sitemap.xml
`, h.baseURL)
}
