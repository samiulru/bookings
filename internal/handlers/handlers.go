package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/samiulru/bookings/internal/config"
	"github.com/samiulru/bookings/internal/driver"
	"github.com/samiulru/bookings/internal/forms"
	"github.com/samiulru/bookings/internal/helpers"
	"github.com/samiulru/bookings/internal/models"
	"github.com/samiulru/bookings/internal/render"
	"github.com/samiulru/bookings/internal/repository"
	"github.com/samiulru/bookings/internal/repository/dbrepo"
)

// Repo is the Repository used by handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a testing repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
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
	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Invalid start date format")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Invalid end date format")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Database Error while searching for all rooms! Try Again")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
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
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)

	_ = render.TemplatesRenderer(w, r, "choose-rooms.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// SearchAvailabilityJSON handles search availability page for JSON
func (m *Repository) SearchAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		//can't parse form
		resp := jsonResponse{
			Ok:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	id := r.Form.Get("room_id")

	start_date, _ := time.Parse("02-01-2006", sd)
	end_date, _ := time.Parse("02-01-2006", ed)
	roomID, err := strconv.Atoi(id)
	if err != nil {
		resp := jsonResponse{
			Ok:      false,
			Message: "Error connecting to database",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(start_date, end_date, roomID)
	if err != nil {
		resp := jsonResponse{
			Ok:      false,
			Message: "Error connecting to database",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	resp := jsonResponse{
		Ok:        available,
		Message:   "",
		RoomID:    strconv.Itoa(roomID),
		StartDate: sd,
		EndDate:   ed,
	}

	out, _ := json.MarshalIndent(resp, "", "     ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// ChooseRoom handles room selection for the users
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// split the URL up by /, and grab the 3rd element
	URIPartition := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(URIPartition[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get room info! Please try again")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get room info! Please try again")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.Room = room
	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookNow handles room booking for the users
func (m *Repository) BookNow(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get roomID from url! Please try again")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get room info! Please try again")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")
	layout := "02-01-2006"
	start_date, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date! Please try again")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	end_date, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date! Please try again")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res := models.Reservation{
		StartDate: start_date,
		EndDate:   end_date,
		RoomID:    roomID,
		Room:      room,
	}
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// Reservation handles reservation form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation info from the session!Please try agian")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	strMap := map[string]string{}
	strMap["start_date"] = res.StartDate.Format("02-01-2006")
	strMap["end_date"] = res.EndDate.Format("02-01-2006")
	strMap["room_name"] = res.Room.RoomName

	data := make(map[string]interface{})
	data["reservation"] = res

	render.TemplatesRenderer(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: strMap,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form info! Please try agian")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation info from the session! Please try again")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.MobileNumber = r.Form.Get("mobile_number")

	form := forms.New(r.PostForm)
	form.Required("first_name", "email", "mobile_number")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() { //Invalid user input
		data := make(map[string]interface{})
		data["reservation"] = res
		strMap := map[string]string{}
		strMap["start_date"] = res.StartDate.Format("02-01-2006")
		strMap["end_date"] = res.EndDate.Format("02-01-2006")
		strMap["room_name"] = res.Room.RoomName

		render.TemplatesRenderer(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: strMap,
		})
		http.Error(w, "", http.StatusSeeOther)
		return
	}

	newReservationID, err := m.DB.InsertReservations(res) //Updating reservation table for successful reservation
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation info! Please try agian")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	roomRestrictions := models.RoomRestriction{
		StartDate:     res.StartDate,
		EndDate:       res.EndDate,
		RoomID:        res.RoomID,
		ReservationID: newReservationID,
		RestrictionId: 1,
	}
	err = m.DB.InsertRoomRestriction(roomRestrictions) //Updating room_restrictions databases since reservation is successful
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction info!Please try agian")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//Send confirmation notification to the client
	//
	//mailData := models.MailData{
	//	From:     "samiulprogramming@gmal.com",
	//	To:       res.Email,
	//	Subject:  "Reservation Confirmation",
	//	Content:  "clientMailContent.html",
	//	Template: "basic.html",
	//}
	//m.App.MailChan <- mailData

	//Send confirmation notification to the Owner
	//
	//mailData = models.MailData{
	//	From:     "samiulprogramming@gmal.com",
	//	To:       "samiul@gmail.com",
	//	Subject:  "Reservation Confirmation",
	//	Content:  "ownerMailContent.html",
	//	Template: "basic.html",
	//}
	//m.App.MailChan <- mailData

	m.App.Session.Put(r.Context(), "reservation", res) //Put user input to the session manager
	m.App.Session.Put(r.Context(), "flash", "Congratulations! Room Reservation Succeed")

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther) //redirecting to the path

}

// ReservationSummary renders the reservation information
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		//m.App.ErrorLog.Println("Can't get information from session")
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

/*......................................................................
....................Admin Tools Handler Functions....................
......................................................................*/

// UserLogin handles UserLogin page
func (m *Repository) UserLogin(w http.ResponseWriter, r *http.Request) {
	_ = render.TemplatesRenderer(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostUserLogin handles authentication and Login of users
func (m *Repository) PostUserLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		log.Println(err)
		render.TemplatesRenderer(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		http.Error(w, "", http.StatusSeeOther)
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "You are logged in successfully")
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

// AdminLogout logs an Admin out
func (m *Repository) AdminLogout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

// AdminDashboard handles Admins dashboard
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.TemplatesRenderer(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

// AdminNewReservations shows lists of new reservations to the admin panel
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.ViewNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.TemplatesRenderer(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminAllReservations shows list of all reservations to the admin panel
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.ViewALlReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.TemplatesRenderer(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminShowReservation shows the reservation calendar to the admin panel
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	// split the URL up by /, and grab the 4th element as src and 5th element as id
	URIPartition := strings.Split(r.RequestURI, "/")
	src := URIPartition[3]
	id, err := strconv.Atoi(URIPartition[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap := make(map[string]string)
	stringMap["src"] = src
	stringMap["year"] = year
	stringMap["month"] = month
	data := make(map[string]interface{})
	data["reservation"] = res
	render.TemplatesRenderer(w, r, "admin-show-reservation.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
		Form:      forms.New(nil),
	})
}

// AdminPostShowReservation updates the reservation to the database
func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	// split the URL up by /, and grab the 4th element as src and 5th element as id
	URIPartition := strings.Split(r.RequestURI, "/")
	src := URIPartition[3]
	id, err := strconv.Atoi(URIPartition[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.MobileNumber = r.Form.Get("mobile_number")

	form := forms.New(r.PostForm)
	form.Required("first_name", "email", "mobile_number")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	//Invalid user input
	if !form.Valid() {
		stringMap := make(map[string]string)
		stringMap["src"] = src
		data := make(map[string]interface{})
		data["reservation"] = res

		render.TemplatesRenderer(w, r, "admin-show-reservation.page.tmpl", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		})
		return
	}

	stringMap := make(map[string]string)
	stringMap["src"] = src
	data := make(map[string]interface{})
	data["reservation"] = res

	err = m.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year := r.Form.Get("year")
	month := r.Form.Get("month")
	m.App.Session.Put(r.Context(), "flash", "Change Saved")

	if year == ""{
		http.Redirect(w, r, fmt.Sprintf("/admin/%s-reservations", src), http.StatusSeeOther)
	} else{
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calender?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
	
}

// AdminProcessReservation mark the reservation as processed
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	// split the URL up by /, and grab the 4th element as src and 5th element as id
	URIPartition := strings.Split(r.RequestURI, "/")
	src := URIPartition[3]
	id, err := strconv.Atoi(URIPartition[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	err = m.DB.UpdateProcessedForReservation(id, 1)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")
	if year == ""{
		http.Redirect(w, r, fmt.Sprintf("/admin/%s-reservations", src), http.StatusSeeOther)
	} else{
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calender?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminDeleteReservation deletes reservation from the database
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	// split the URL up by /, and grab the 4th element as src and 5th element as id
	URIPartition := strings.Split(r.RequestURI, "/")
	src := URIPartition[3]
	id, err := strconv.Atoi(URIPartition[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	err = m.DB.DeleteReservation(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation Deleted Successfully")

	if year == ""{
		http.Redirect(w, r, fmt.Sprintf("/admin/%s-reservations", src), http.StatusSeeOther)
	} else{
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calender?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminReservationsCalender shows the reservation calendar to the admin panel
func (m *Repository) AdminReservationsCalender(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}
	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear
	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	//get the first and last days of the month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
	}

	data := make(map[string]interface{})
	data["now"] = now
	data["rooms"] = rooms

	for _, x := range rooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; !d.After(lastOfMonth); d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2-01-2006")] = 0
			blockMap[d.Format("2-01-2006")] = 0
		}

		// get all restrictions for the current room
		roomRestrictions, err := m.DB.GetRestrictionsForRoomByDate(x.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		for _, y := range roomRestrictions {
			if y.ReservationID > 0 {
				//It's a reservations
				for d := y.StartDate; !d.After(y.EndDate); d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2-01-2006")] = y.ReservationID
				}

			} else {
				//It's a block
				blockMap[y.StartDate.Format("2-01-2006")] = y.ID
			}
		}

		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap
		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)
	}

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	render.TemplatesRenderer(w, r, "admin-reservations-calender.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}

// AdminPostReservationsCalender handles the post request for the Reservation calender page
func (m *Repository) AdminPostReservationsCalender(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	form := forms.New(r.PostForm)

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	for _, x := range rooms {
		blockMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.ID)).(map[string]int)
		for name, value := range blockMap {
			if value > 0{
				if !form.Has(fmt.Sprintf("remove_block_%d_%s", x.ID, name)){
					err := m.DB.DeleteBlockForRoom(value)
					if err != nil {
						helpers.ServerError(w, err)
						return
					}
				}
			}
		}
	}

	for name, _ := range r.PostForm{
		if strings.HasPrefix(name, "add_block_"){
			URLPartition := strings.Split(name, "_")
			roomID, _ := strconv.Atoi(URLPartition[2])
			date, _ := time.Parse("2-01-2006", URLPartition[3])

			err := m.DB.InsertBlockForRoom(roomID, date, date.AddDate(0,0,1,))
			if err != nil {
				helpers.ServerError(w, err)
				return
			}
		}
	}

	m.App.Session.Put(r.Context(), "flash", "Changes Saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calender?y=%d&m=%d", year, month), http.StatusSeeOther)
}
