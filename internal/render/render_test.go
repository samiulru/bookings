package render

import (
	"github.com/samiulru/bookings/internal/models"
	"net/http"
	"testing"
)

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()

	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "77105")
	result := AddDefaultData(&td, r)
	if result.Flash != "77105" {
		t.Error("flash value of 77105 not found in session")
	}
}

var checkList = map[string]string{
	"Home Page":       "home.page.tmpl",
	"Economical Page": "economical.page.tmpl",
	"Premium Page":    "premium.page.tmpl",
}

func TestTemplatesRenderer(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	var ww respWriter

	for i, v := range checkList {
		err = TemplatesRenderer(&ww, r, v, &models.TemplateData{})
		if err != nil {
			t.Errorf("error while writing %s to the browser", i)
		}

	}

}
func TestCreateTemplateCache(t *testing.T) {
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, err

}
