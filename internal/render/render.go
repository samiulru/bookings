package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/samiulru/bookings/internal/config"
	"github.com/samiulru/bookings/internal/models"
)

var funcMap = template.FuncMap{
	"PadLeft":  PadLeft,
	"PadRight": PadRight,
}
// PadLeft pads the input string with spaces on the left to reach the specified width.
func PadLeft(s string, width int) string {
	return fmt.Sprintf("%-*s", width, s)
}

// PadRight pads the input string with spaces on the right to reach the specified width.
func PadRight(s string, width int) string {
	return fmt.Sprintf("%*s", width, s)
}

var app *config.AppConfig

// NewTemplates sets the config for template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData sets the template data for each handler
func AddDefaultData(data *models.TemplateData, r *http.Request) *models.TemplateData {
	data.Flash = app.Session.PopString(r.Context(), "flash")
	data.Error = app.Session.PopString(r.Context(), "error")
	data.Warning = app.Session.PopString(r.Context(), "warning")
	data.CSRFToken = nosurf.Token(r)
	return data
}

// TemplatesRenderer renders templates specified by the templateName with the help of html/template package
func TemplatesRenderer(w http.ResponseWriter, r *http.Request, templateName string, data *models.TemplateData) {
	////checksErr checks if there is any error and stops the app immediately after printing the error logs
	var err error
	checksErr := func(msg string) {
		if err != nil {
			log.Fatal(msg)
		}
	}
	//getting template cache from AppConfig
	tmplCache := app.TemplateCache
	if !app.UseCache {
		tmplCache, err = CreateTemplateCache()
		checksErr("Error occur while creating new template cache")
	}
	tmpl, ok := tmplCache[templateName]
	if !ok {
		log.Fatal("Could not get the template from the template cache")
	}

	buf := new(bytes.Buffer)
	td := AddDefaultData(data, r)
	err = tmpl.Execute(buf, td)
	checksErr("Error occur while executing template")

	_, err = buf.WriteTo(w)
	checksErr("Error occur while writing to the response writer")
}

// CreateTemplateCache creates templates cache
func CreateTemplateCache() (map[string]*template.Template, error) {
	tmplCache := map[string]*template.Template{}

	//printErr checks and print these errors, if there is any error 
	printErr := func(err error) {
		if err != nil {
			fmt.Println("Error occur within the CreateTemplateCache function:", err)
		}
	}
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	printErr(err)

	for _, pg := range pages {
		name := filepath.Base(pg)
		ts, err := template.New(name).Funcs(funcMap).ParseFiles(pg)
		printErr(err)

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		printErr(err)
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			printErr(err)
		}

		tmplCache[name] = ts
	}

	return tmplCache, nil
}

