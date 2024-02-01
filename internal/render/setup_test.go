package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/samiulru/bookings/internal/config"
	"github.com/samiulru/bookings/internal/models"
	"net/http"
	"os"
	"testing"
	"time"
)

var testApp config.AppConfig
var session *scs.SessionManager

func TestMain(m *testing.M) {
	//What I am going to put in the session
	gob.Register(models.Reservation{})
	//Sessions for users
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true //change it to false if it needs to delete cookie at the closing of the browser
	session.Cookie.Secure = false //local host is insecure connection which is used in InProduction mode
	session.Cookie.SameSite = http.SameSiteLaxMode
	//Setting up the app-config values
	testApp.InProduction = false //change it to true when in developer mode
	testApp.Session = session
	app = &testApp
	os.Exit(m.Run())
}

type respWriter struct{}

func (myWriter *respWriter) Header() http.Header {
	var header http.Header
	return header
}

func (myWriter *respWriter) Write(b []byte) (int, error) {
	n := len(b)
	return n, nil
}

func (myWriter *respWriter) WriteHeader(statusCode int) {

}
