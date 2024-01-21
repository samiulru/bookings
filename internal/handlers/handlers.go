package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/samiulru/bookings/internal/config"
	"github.com/samiulru/bookings/internal/models"
	"github.com/samiulru/bookings/internal/render"
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
	render.TemplatesRenderer(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About handles the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.TemplatesRenderer(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Economical handles the room page
func (m *Repository) Economical(w http.ResponseWriter, r *http.Request) {
	render.TemplatesRenderer(w, r, "economical.page.tmpl", &models.TemplateData{})
}

// Premium handles the room page
func (m *Repository) Premium(w http.ResponseWriter, r *http.Request) {
	render.TemplatesRenderer(w, r, "premium.page.tmpl", &models.TemplateData{})
}

// SearchAvailability handles search availability page for GET request
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.TemplatesRenderer(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostSearchAvailability handles search availability page for POST request
func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.FormValue("arrival")
	end := r.FormValue("deperture")
	w.Write([]byte(fmt.Sprintf("Starting date: %s and Ending date: %s", start, end)))
}

type jsonResponse struct{
	Ok bool `json:"ok"`
	Message string `json:"message"`
}

// PostSearchAvailability handles search availability page for POST request
func (m *Repository) SearchAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		Ok: true,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Reservation handles Make-reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.TemplatesRenderer(w, r, "make-reservation.page.tmpl", &models.TemplateData{})
}

// Contact handles contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.TemplatesRenderer(w, r, "contact.page.tmpl", &models.TemplateData{})
}
