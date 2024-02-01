package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/samiulru/bookings/internal/config"
	"github.com/samiulru/bookings/internal/driver"
	"github.com/samiulru/bookings/internal/forms"
	"github.com/samiulru/bookings/internal/helpers"
	"github.com/samiulru/bookings/internal/models"
	"github.com/samiulru/bookings/internal/render"
	"github.com/samiulru/bookings/internal/repository"
	"github.com/samiulru/bookings/internal/repository/dbrepo"
	"net/http"
	"strconv"
	"time"
)

// Repo is the Repository used by handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a ner repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandler sets the repository for the handlers
func NewHandler(r *Repository) {
	Repo = r
}

// Home handles the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	_ = render.TemplatesRenderer(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About handles the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	_ = render.TemplatesRenderer(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Economical handles the room page
func (m *Repository) Economical(w http.ResponseWriter, r *http.Request) {
	_ = render.TemplatesRenderer(w, r, "economical.page.tmpl", &models.TemplateData{})
}

// Premium handles the room page
func (m *Repository) Premium(w http.ResponseWriter, r *http.Request) {
	_ = render.TemplatesRenderer(w, r, "premium.page.tmpl", &models.TemplateData{})
}

// Contact handles contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	_ = render.TemplatesRenderer(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// SearchAvailability handles search availability page for GET request
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	_ = render.TemplatesRenderer(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostSearchAvailability handles search availability page for POST request
func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.FormValue("start_date")
	end := r.FormValue("end_date")
	_, _ = w.Write([]byte(fmt.Sprintf("Starting date: %s and Ending date: %s", start, end)))
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// SearchAvailabilityJSON handles search availability page for JSON
func (m *Repository) SearchAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		Ok:      true,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(out)
}

// Reservation handles reservation form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	err := render.TemplatesRenderer(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	//Form date layout
	layout := "02-01-2006"
	start_date, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	end_date, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	room_id, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	reservation := models.Reservation{
		FirstName:    r.Form.Get("first_name"),
		LastName:     r.Form.Get("last_name"),
		Email:        r.Form.Get("email"),
		MobileNumber: r.Form.Get("mobile_number"),
		StartDate:    start_date,
		EndDate:      end_date,
		RoomID:       room_id,
	}
	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "mobile_number")
	form.MinLength("first_name", 3)
	form.MinLength("last_name", 3)
	form.IsEmail("email")

	if !form.Valid() { //Invalid user input
		data := make(map[string]interface{})
		data["reservation"] = reservation
		err = render.TemplatesRenderer(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservations(reservation) //Updating reservation table for successful reservation
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	roomRestrictions := models.RoomRestriction{
		StartDate:     start_date,
		EndDate:       end_date,
		RoomID:        room_id,
		ReservationID: newReservationID,
		RestrictionId: 1,
	}
	err = m.DB.InsertRoomRestriction(roomRestrictions) //Updating room_restrictions databases since reservation is successful
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "reservation", reservation) //Put user input to the session manager

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther) //redirecting to the path

}

// ReservationSummary renders the reservation information
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get information from session")
		m.App.Session.Put(r.Context(), "error", "Internal Error! Can't get reservation information from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	_ = render.TemplatesRenderer(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
