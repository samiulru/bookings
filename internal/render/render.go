package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/justinas/nosurf"
	"github.com/samiulru/bookings/internal/config"
	"github.com/samiulru/bookings/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
)

var funcMap = template.FuncMap{}

var app *config.AppConfig
var pathToTemplates = "./templates"

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

	if app.Session.Exists(r.Context(), "user_id") {
		data.IsAuthenticated = 1
	}
	return data
}

// TemplatesRenderer renders templates specified by the templateName with the help of html/template package
func TemplatesRenderer(w http.ResponseWriter, r *http.Request, templateName string, data *models.TemplateData) error {
	var tmplCache map[string]*template.Template
	var err error

	if app.UseCache {
		//getting template cache from AppConfig
		tmplCache = app.TemplateCache
	} else {
		//This is for testing, so this will rebuild the template cache on every request
		tmplCache, err = CreateTemplateCache()
		if err != nil {
			return err
		}

	}
	tmpl, ok := tmplCache[templateName]
	if !ok {
		return errors.New("cannot get the template from the template cache")
	}

	buf := new(bytes.Buffer)
	td := AddDefaultData(data, r)
	err = tmpl.Execute(buf, td)
	if err != nil {
		fmt.Println("error while executing templates")
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("error while writing template to the browser")
		return err
	}

	return nil

}

// CreateTemplateCache creates templates cache
func CreateTemplateCache() (map[string]*template.Template, error) {
	tmplCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return tmplCache, errors.New(fmt.Sprint("Error occur within the CreateTemplateCache function:", err))
	}

	for _, pg := range pages {
		name := filepath.Base(pg)
		ts, err := template.New(name).Funcs(funcMap).ParseFiles(pg)
		if err != nil {
			return tmplCache, errors.New(fmt.Sprint("Error occur within the CreateTemplateCache function:", err))
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return tmplCache, errors.New(fmt.Sprint("Error occur within the CreateTemplateCache function:", err))
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return tmplCache, errors.New(fmt.Sprint("Error occur within the CreateTemplateCache function:", err))
			}
		}

		tmplCache[name] = ts
	}

	return tmplCache, nil
}
