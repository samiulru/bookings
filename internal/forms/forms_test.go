package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var f *Form

func TestNew(t *testing.T) {
	//do nothing
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required value missing")
	}

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "b")
	postData.Add("c", "c")
	r, _ = http.NewRequest("POST", "/some-url", nil)
	r.PostForm = postData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows does not have required fields when it does")
	}
}

func TestForm_MinLength(t *testing.T) {
	postData := url.Values{}
	form := New(postData)
	form.MinLength("x", 3)
	if form.Valid() {
		t.Error("form shows value for a non-existent value")
	}

	postData = url.Values{}
	postData.Add("a", "samiul")
	form = New(postData)
	form.MinLength("a", 10)
	if form.Valid() {
		t.Error("shows minlength criterion is fulfilled while it doesn't")
	}
	postData = url.Values{}
	postData.Add("a", "samiul")
	form = New(postData)
	form.MinLength("a", 3)
	if !form.Valid() {
		t.Error("shows minlength criterion is fulfilled while it doesn't")
	}

}

func TestForm_Has(t *testing.T) {
	postData := url.Values{}
	form := New(postData)
	has := form.Has("whatever")
	if has {
		t.Error("shows field exist while it doesn't")
	}

	postData = url.Values{}
	postData.Add("a", "a")

	form = New(postData)
	has = form.Has("a")
	if !has {
		t.Error("shows field does not exist while it does")
	}

}
func TestForm_Valid(t *testing.T) {
	postData := url.Values{}
	form := New(postData)
	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should be valid")
	}
}
func TestForm_IsEmail(t *testing.T) {
	postData := url.Values{}
	form := New(postData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("form shows valid email for non-existent email")
	}
	isError := form.Err.Get("email")
	if isError == "" {
		t.Error("Get func gives no error when should have an error")
	}
	postData = url.Values{}
	postData.Add("email", "coding.samiul@gma")
	form = New(postData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("form shows valid email for invalid email")
	}
	postData = url.Values{}
	postData.Add("email", "coding.samiul@gmail.com")
	form = New(postData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("form shows invalid email for valid email")
	}
	isError = form.Err.Get("email")
	if isError != "" {
		t.Error("Get func gives error when shouldn't have any error")
	}
}
