package handlers

import (
	"github.com/samiulru/bookings/pkg/config"
	"github.com/samiulru/bookings/pkg/models"
	"github.com/samiulru/bookings/pkg/render"
	"net/http"
)

// Repo is the Repository used by handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a ner repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandler sets the repository for the handlers
func NewHandler(r *Repository) {
	Repo = r
}

// Home page handlers
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	remoteIP := r.RemoteAddr
	stringMap["author"] = "Samiul Bashir"
	stringMap["contact_email"] = "coding.samiul@gmail.com"
	stringMap["github"] = "https://github.com/samiulru"
	stringMap["remote_ip"] = remoteIP

	m.App.Session.Put(r.Context(), "author", stringMap["author"])
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.TemplatesRenderer(w, "home.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// About page handlers
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["remote_ip"] = m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["author"] = m.App.Session.GetString(r.Context(), "author")
	render.TemplatesRenderer(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})

}
