package main

import (
	"encoding/gob"
	"fmt"
	"github.com/samiulru/bookings/internal/driver"
	"github.com/samiulru/bookings/internal/helpers"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/samiulru/bookings/internal/config"
	"github.com/samiulru/bookings/internal/handlers"
	"github.com/samiulru/bookings/internal/models"
	"github.com/samiulru/bookings/internal/render"
)

// specified port that is listen to serve web request
const portNumber = ":10526"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// The webapp entry point
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	//Http server for our web app
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	fmt.Println("WebApp run on port:", portNumber)
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	//What I am going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	//Creating template cache
	tmplCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot get the template files")
		return nil, err
	}
	//Sessions for users
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true            //change it to false if it needs to delete cookie at the closing of the browser
	session.Cookie.Secure = app.InProduction //local host is insecure connection which is used in InProduction mode

	//Setting up the app-config values
	app.UseCache = false //false when in developer mode
	app.TemplateCache = tmplCache
	app.InProduction = false //change it to true when in developer mode
	app.Session = session

	// connect to Database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=samiul@10526")
	if err != nil {
		log.Fatal("Cannot connect to database")
		return nil, err
	}
	log.Println("Success! Connected to database!!")
	infoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandler(repo)
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)
	return db, nil
}
