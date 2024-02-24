package main

import (
	"encoding/gob"
	"flag"
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
	//
	//defer close(app.MailChan)
	//fmt.Println("Mail Server is running on port: 1025")
	//listenForMail()

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
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})
	gob.Register(map[string]int{})

	// read flags
	inProduction := flag.Bool("production", true, "Application is in production mode")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbsSSL := flag.String("dbssl", "disable", "Database SSL settings (disable, prefer, require)")

	flag.Parse()
	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	//mailChan carries all the mail from any part of the app
	//mailChan := make(chan models.MailData)
	//app.MailChan = mailChan

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
	session.Cookie.Secure = app.InProduction //localhost is insecure connection which is used in InProduction mode

	//Setting up the app-config values
	app.UseCache = *useCache //true when In-Production mode and false when in developer mode
	app.TemplateCache = tmplCache
	app.InProduction = *inProduction //change it to true when in developer mode
	app.Session = session

	// connect to Database
	log.Println("Connecting to database...")

	dbConnection := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbsSSL)
	db, err := driver.ConnectSQL(dbConnection)
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
