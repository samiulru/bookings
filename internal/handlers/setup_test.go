package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/samiulru/bookings/internal/config"
	"github.com/samiulru/bookings/internal/models"
	"github.com/samiulru/bookings/internal/render"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var funcMap = template.FuncMap{}
var pathToTemplates = "./../../templates"
var app config.AppConfig
var session *scs.SessionManager

func TestMain(m *testing.M) {
	//What I am going to put in the session
	gob.Register(models.Reservation{})
	//Creating template cache
	tmplCache, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Cannot get the template files")
	}
	//Sessions for users
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true            //change it to false if it needs to delete cookie at the closing of the browser
	session.Cookie.Secure = app.InProduction //local host is insecure connection which is used in InProduction mode

	//Setting up the app-config values
	app.UseCache = true //false when in developer mode
	app.TemplateCache = tmplCache
	app.InProduction = false //change it to true when in developer mode
	app.Session = session
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	repo := NewTestRepo(&app)
	NewHandler(repo)
	render.NewTemplates(&app)

	os.Exit(m.Run())
}

func getRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/economical", Repo.Economical)
	mux.Get("/premium", Repo.Premium)
	mux.Get("/contact", Repo.Contact)

	mux.Get("/search-availability", Repo.SearchAvailability)
	mux.Post("/search-availability", Repo.PostSearchAvailability)
	mux.Post("/search-availability-json", Repo.SearchAvailabilityJSON)

	mux.Get("/choose-room/{id}", Repo.ChooseRoom)
	mux.Get("/book-now", Repo.BookNow)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	//FileServer for serving files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux

}

// NoSurf adds CSRF protection to all post request
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session in every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache creates templates cache
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	tmplCache := map[string]*template.Template{}

	//printErr checks and print these errors, if there is any error
	printErr := func(err error) {
		if err != nil {
			fmt.Println("Error occur within the CreateTemplateCache function:", err)
		}
	}
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	printErr(err)

	for _, pg := range pages {
		name := filepath.Base(pg)
		ts, err := template.New(name).Funcs(funcMap).ParseFiles(pg)
		printErr(err)

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		printErr(err)
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			printErr(err)
		}

		tmplCache[name] = ts
	}

	return tmplCache, nil
}
