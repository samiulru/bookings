package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
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

	layout := "02-01-2006"
	start_date, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	end_date, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(start_date, end_date)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No room available")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: start_date,
		EndDate:   end_date,
	}
	m.App.Session.Put(r.Context(), "reservation", res)
	m.App.Session.Put(r.Context(), "rooms", rooms)

	err = render.TemplatesRenderer(w, r, "choose-rooms.page.tmpl", &models.TemplateData{
		Data: data,
	})
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
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

// ChooseRoom handles room selection for ther users
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation information from the form, please retry"))
		return
	}
	AvailableRooms, ok := m.App.Session.Get(r.Context(), "rooms").([]models.Room)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation information from the form, please retry"))
		return
	}
	for _, value := range AvailableRooms {
		if value.ID == roomID {
			res.Room.RoomName = value.RoomName
			break
		}
	}
	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

// Reservation handles reservation form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation information from the form, please retry"))
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = res

	strMap := map[string]string{}
	strMap["start_date"] = res.StartDate.Format("02-01-2006")
	strMap["end_date"] = res.EndDate.Format("02-01-2006")
	strMap["room_name"] = res.Room.RoomName

	err := render.TemplatesRenderer(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: strMap,
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
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation information from the form, please retry"))
		return
	}
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.MobileNumber = r.Form.Get("mobile_number")

	form := forms.New(r.PostForm)
	form.Required("first_name", "email", "mobile_number")
	form.MinLength("first_name", 3)
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
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
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
