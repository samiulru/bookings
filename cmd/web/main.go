package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/samiulru/bookings/pkg/config"
	"github.com/samiulru/bookings/pkg/handlers"
	"github.com/samiulru/bookings/pkg/render"
	"github.com/samiulru/bookings/pkg/test"
	"log"
	"net/http"
	"time"
)

// specified port that is listen to serve web request
const portNumber = ":10526"

var app config.AppConfig
var session *scs.SessionManager

// The webapp entry point
func main() {
	//Creating template cache
	tmplCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot get the template files")
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

	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)
	render.NewTemplates(&app)
	test.Main(&app)
	//Http server for our web app
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	fmt.Println("WebApp run on port:", portNumber)
	err = srv.ListenAndServe()
	log.Fatal(err)
}
