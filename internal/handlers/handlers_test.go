package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/samiulru/bookings/internal/driver"
	"github.com/samiulru/bookings/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"premium", "/premium", "GET", http.StatusOK},
	{"economical", "/economical", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	//{"search-availability-post", "/search-availability", "POST", []postData{
	//		{"start_date", "2024-01-01"},
	//		{"end_date", "2024-01-01"},
	//	}, http.StatusOK},
	//{"search-availability-json", "/search-availability-json", "POST", []postData{
	//	{"start_date", "2024-01-01"},
	//	{"end_date", "2024-01-01"},
	//}, http.StatusOK},
	//{"make-reservation", "/make-reservation", "POST", []postData{
	//	{"first_name", "Samiul"},
	//	{"last_name", "Bashir"},
	//	{"email", "samiul@gmail.com"},
	//	{"mobile_number", "01742135093"},
	//}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes) //ts is the test server
	defer ts.Close()
	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("For %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}

}

// TestNewRepo tests the NewRepo function returns correct type repository or not
func TestNewRepo(t *testing.T) {
	var db driver.DB
	testRepo := NewRepo(&app, &db)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, wanted *Repository", reflect.TypeOf(testRepo).String())
	}
}

// TestRepository_Reservation tests the make-reservation handler
func TestRepository_Reservation(t *testing.T) {

	//Test Case:00 No error
	Reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			RoomName: "Economical Quarter",
			ID:       1,
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)

	ctx := getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", Reservation)
	handler := http.HandlerFunc(Repo.Reservation)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Test Case:00 No error:::Reservation handler wrong response code %d, expected %d", rr.Code, http.StatusOK)
	}

	//Test Case:01 For Reservation not in the session
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	req = req.WithContext(getCtx(req))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Test Case:01 For Reservation not in the session:::Reservation handler wrong response code %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

// getCtx returns context.context
func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), "X-Session")
	if err != nil {
		log.Println(err)
	}
	return ctx
}

// TestRepository_ChooseRoom tests the ChooseRoom handler
func TestRepository_ChooseRoom(t *testing.T) {
	res := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			RoomName: "Economical Quarter",
			ID:       1,
		},
	}

	//Test Case:00 No error
	req, _ := http.NewRequest("GET", "/choose-room/1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/1"
	session.Put(ctx, "reservation", res)
	handler := http.HandlerFunc(Repo.ChooseRoom)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Test Case:00 No error::: ChooseRoom handler wrong response code %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	//Test Case:02 Reservation not in the session(reset request)
	req, _ = http.NewRequest("GET", "/choose-room/1", nil)
	req = req.WithContext(getCtx(req))
	req.RequestURI = "/choose-room/1"
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Test Case:02 Reservation not in the session::: ChooseRoom handler wrong response code %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}
	//Test Case:03 Non-existent RoomID
	req, _ = http.NewRequest("GET", "/choose-room/noting", nil)
	req = req.WithContext(getCtx(req))
	req.RequestURI = "/choose-room/noting"
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Test Case:03 Non-existent RoomID::: ChooseRoom handler wrong response code %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//Test Case:04 Missing RoomID
	req, _ = http.NewRequest("GET", "/choose-room/10", nil)
	req = req.WithContext(getCtx(req))
	req.RequestURI = "/choose-room/10"
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Test Case:04 Missing RoomID::: ChooseRoom handler wrong response code %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

// TestRepository_PostReservation tests the form in post request
func TestRepository_PostReservation(t *testing.T) {

	//Test Case:00 No error in PostReservation handler
	validReqBody := "start_date=01-01-2050"
	validReqBody = fmt.Sprintf("%s&%s", validReqBody, "end_date=11-01-2050")
	validReqBody = fmt.Sprintf("%s&%s", validReqBody, "first_name=Samiul")
	validReqBody = fmt.Sprintf("%s&%s", validReqBody, "last_name=Bashir")
	validReqBody = fmt.Sprintf("%s&%s", validReqBody, "email=samiulru@outlook.com")
	validReqBody = fmt.Sprintf("%s&%s", validReqBody, "mobile_number=01797331020")
	validReqBody = fmt.Sprintf("%s&%s", validReqBody, "room_id=1")

	layout := "02-01-2006"
	sd, _ := time.Parse(layout, "01-01-2050")
	ed, _ := time.Parse(layout, "11-01-2050")
	res := models.Reservation{
		StartDate: sd,
		EndDate:   ed,
		RoomID:    1,
		Room: models.Room{
			RoomName: "Economical Quarter",
			ID:       1,
		},
	}
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(validReqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	session.Put(req.Context(), "reservation", res)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Test Case:00 No error in PostReservation handler::::wrong response code: %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	//Test Case:01 no form info
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Test Case:01 no form info::::wrong response code for no form info: %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//Test Case:02 reservation info exist in the session
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(validReqBody))
	req = req.WithContext(getCtx(req))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Test Case:02 reservation info exist in the session:::wrong response code %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//Test Case:03 invalid form
	invalidReqBody := "start_date=01-01-2050"
	invalidReqBody = fmt.Sprintf("%s&%s", invalidReqBody, "end_date=11-01-2050")
	invalidReqBody = fmt.Sprintf("%s&%s", invalidReqBody, "first_name=S")
	invalidReqBody = fmt.Sprintf("%s&%s", invalidReqBody, "last_name=Bashir")
	invalidReqBody = fmt.Sprintf("%s&%s", invalidReqBody, "email=samiul@gmail.com")
	invalidReqBody = fmt.Sprintf("%s&%s", invalidReqBody, "mobile_number=01797331020")
	invalidReqBody = fmt.Sprintf("%s&%s", invalidReqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(invalidReqBody))
	req = req.WithContext(getCtx(req))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	session.Put(req.Context(), "reservation", res)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Test Case:03 invalid form:::wrong response code for invalid form %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	//Test Case:04 Unable to insertReservation
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(validReqBody))
	req = req.WithContext(getCtx(req))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res.RoomID = 100
	session.Put(req.Context(), "reservation", res)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Test Case:04 Unable to insertReservation:::wrong response code for invalid form %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}
	//Test Case:05 Unable to insertReservation
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(validReqBody))
	req = req.WithContext(getCtx(req))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res.RoomID = 200
	session.Put(req.Context(), "reservation", res)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Test Case:05 Unable to insert Room Restriction:::wrong response code for invalid form %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

// TestRepository_PostReservation tests the form in post request
func TestRepository_ReservationSummary(t *testing.T) {

	//Test Case:-=01 No reservation info doesn't exist in the session
	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservaionSummary handler wrong response code %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//Test Case:-=02 Reservation info exist in the session
	res := models.Reservation{
		ID:           1,
		FirstName:    "Samiul",
		LastName:     "Bashir",
		Email:        "samiul@gmail.com",
		MobileNumber: "01797331020",
		StartDate:    time.Time{},
		EndDate:      time.Time{},
		RoomID:       0,
		Room: models.Room{
			RoomName: "Economical Quarter",
			ID:       1,
		},
		CreatedAt: time.Time{},
		UpdateAt:  time.Time{},
	}
	req, _ = http.NewRequest("GET", "/reservation-summary", nil)
	req = req.WithContext(getCtx(req))
	session.Put(req.Context(), "reservation", res)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)
}

// TestRepository_BookNow tests BookNow Handler
func TestRepository_BookNow(t *testing.T) {

	//Test Case:00 No error
	req, _ := http.NewRequest("GET", "/book-room?s=01-01-2050&e=03-01-2050&id=1", nil)
	req = req.WithContext(getCtx(req))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.BookNow)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Test Case:00 No error::: BookNow Handler wrong response code %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	//Test Case:01 Missing RoomID
	req, _ = http.NewRequest("GET", "/book-room?s=01-01-2050&e=03-01-2050&id=wrong", nil)
	req = req.WithContext(getCtx(req))
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.BookNow)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Test Case:01 Missing RoomID::: BookNow Handler wrong response code %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//Test Case:02 Invalid RoomID or Database is not working
	req, _ = http.NewRequest("GET", "/book-room?s=01-01-2050&e=03-01-2050&id=30", nil)
	req = req.WithContext(getCtx(req))
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.BookNow)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Test Case:02 Invalid RoomID or Database is not working::: BookNow Handler wrong response code %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//Test Case:03 Invalid Start Date
	req, _ = http.NewRequest("GET", "/book-room?s=start_date_missing&e=03-01-2050&id=1", nil)
	req = req.WithContext(getCtx(req))
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.BookNow)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Test Case:03 Invalid Start Date::: BookNow Handler wrong response code %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//Test Case:04 Invalid End Date
	req, _ = http.NewRequest("GET", "/book-room?s=01-01-2050&e=end_date_missing&id=1", nil)
	req = req.WithContext(getCtx(req))
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.BookNow)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Test Case:04 Invalid End Date::: BookNow Handler wrong response code %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

// TestRepository_PostSearchAvailability tests PostSearchAvailability Handler
func TestRepository_PostSearchAvailability(t *testing.T) {

	//Test Case:00 No Error
	postData := url.Values{}
	postData.Add("start_date", "03-01-2050")
	postData.Add("end_date", "05-01-2050")
	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(postData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostSearchAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Test Case:00 No Error::: BookNow Handler wrong response code %d, expected %d", rr.Code, http.StatusOK)
	}

	//Test Case:01 Missing start date
	postData = url.Values{}
	postData.Add("end_date", "05-01-2050")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostSearchAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Test Case:01 Missing start date::: BookNow Handler wrong response code %d, expected %d", rr.Code, http.StatusSeeOther)
	}
	//Test Case:02 Missing end date
	postData = url.Values{}
	postData.Add("start_date", "03-01-2050")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostSearchAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Test Case:03 Missing end date::: BookNow Handler wrong response code %d, expected %d", rr.Code, http.StatusSeeOther)
	}
	//Test Case:03 No Room Available
	postData = url.Values{}
	postData.Add("start_date", "01-01-2050")
	postData.Add("end_date", "02-01-2050")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostSearchAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Test Case:03 Database error || No Room Available::: BookNow Handler wrong response code %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	//Test Case:04 Invalid Date
	postData = url.Values{}
	postData.Add("start_date", "02-01-2050")
	postData.Add("end_date", "01-01-2050")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostSearchAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Test Case:03 Database error || No Room Available::: BookNow Handler wrong response code %d, expected %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_SearchAvailabilityJSON(t *testing.T) {

	//Test Case:00 No error || Room Avaliable
	postData := url.Values{}
	postData.Add("start_date", "03-01-2050")
	postData.Add("end_date", "05-01-2050")
	postData.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(postData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.SearchAvailabilityJSON)
	handler.ServeHTTP(rr, req)

	var resp jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &resp)
	if err != nil {
		t.Error("Faild to parse json")
	}
	if !resp.Ok {
		t.Error("Test Case:00 No error || Room Avaliable:::Got availability when none was expected in AvailabilityJSON")
	}

	//Test Case:01 No error || No Room Avaliable
	postData = url.Values{}
	postData.Add("start_date", "03-01-2050")
	postData.Add("end_date", "05-01-2050")
	postData.Add("room_id", "2")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.SearchAvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &resp)
	if err != nil {
		t.Error("Faild to parse json")
	}
	if resp.Ok {
		t.Error("Test Case:01 No error || No Room Avaliable:::Got availability when none was expected in AvailabilityJSON")
	}

	//Test Case:02 Cannot Parse form for unknown Content-Type
	postData = url.Values{}
	postData.Add("start_date", "03-01-2050")
	postData.Add("end_date", "05-01-2050")
	postData.Add("room_id", "2")

	req, _ = http.NewRequest("POST", "/search-availability-json", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.SearchAvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &resp)
	if err != nil {
		t.Error("Faild to parse json")
	}
	if resp.Ok {
		t.Error("Test Case:02 Cannot Parse form for unknown Content-Type:::Got availability when none was expected in AvailabilityJSON")
	}
}
