package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MinLength checks the min length of the given field of the form
func (f *Form) MinLenght(field string, length int, r *http.Request) bool {
	x := strings.TrimSpace(r.Form.Get(field))
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This filed must have at least %d characters", length))
		return false
	}
	return true

}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	return true
}

// Valid returns true if ther is no errors, otherwise returns false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

//IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid Email Address")
	}
}
