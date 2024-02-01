package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"premium", "/premium", "GET", []postData{}, http.StatusOK},
	{"economical", "/economical", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"search-availability-post", "/search-availability", "POST", []postData{
		{"start_date", "2024-01-01"},
		{"end_date", "2024-01-01"},
	}, http.StatusOK},
	{"search-availability-json", "/search-availability-json", "POST", []postData{
		{"start_date", "2024-01-01"},
		{"end_date", "2024-01-01"},
	}, http.StatusOK},
	{"make-reservation", "/make-reservation", "POST", []postData{
		{"first_name", "Samiul"},
		{"last_name", "Bashir"},
		{"email", "samiul@gmail.com"},
		{"mobile_number", "01742135093"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	ts := httptest.NewTLSServer(routes) //ts is the test server
	defer ts.Close()
	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("For %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else if e.method == "POST" {
			values := url.Values{}
			for _, v := range e.params {
				values.Add(v.key, v.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("For %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		}

	}

}
