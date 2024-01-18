package test

import (
	"fmt"
	"github.com/samiulru/bookings/pkg/config"
)

func Main(app *config.AppConfig) {

	fmt.Println("-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_")
	fmt.Println("App Configuration: ")
	fmt.Println("-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_")
	fmt.Println("UseCache: ", app.UseCache)
	fmt.Println("TemplateCache: ", app.TemplateCache)
	fmt.Println("InProduction: ", app.InProduction)
	fmt.Println("Session: ", app.Session)
	fmt.Println("-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_")
	fmt.Println("-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_")
}

//UseCache      bool
//TemplateCache map[string]*template.Template
//InProduction  bool
//Session       *scs.SessionManager
