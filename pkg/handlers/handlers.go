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

// Home handles the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.TemplatesRenderer(w, "home.page.tmpl", &models.TemplateData{})
}
// About handles the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.TemplatesRenderer(w, "about.page.tmpl", &models.TemplateData{})
}
// Home handles the about page
func (m *Repository) Economical(w http.ResponseWriter, r *http.Request) {
	render.TemplatesRenderer(w, "economical.page.tmpl", &models.TemplateData{})
}
// Premium handles the room page
func (m *Repository) Premium(w http.ResponseWriter, r *http.Request) {
	render.TemplatesRenderer(w, "premium.page.tmpl", &models.TemplateData{})
}
// SearchAvailability handles search availability page 
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.TemplatesRenderer(w, "search-availability.page.tmpl", &models.TemplateData{})
}
// Reservation handles Make-reservation page 
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.TemplatesRenderer(w, "make-reservation.page.tmpl", &models.TemplateData{})
}
// Contact handles contact page 
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.TemplatesRenderer(w, "contact.page.tmpl", &models.TemplateData{})
}
