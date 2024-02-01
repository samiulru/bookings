package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds an url.Values object
type Form struct {
	url.Values
	Err errors
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		map[string][]string{},
	}
}

// Required checks required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Err.Add(field, "This field cannot be blank")
		}
	}
}

// MinLength checks the min length of the given field of the form
func (f *Form) MinLength(field string, length int) bool {
	x := strings.TrimSpace(f.Get(field))
	if len(x) < length {
		f.Err.Add(field, fmt.Sprintf("This filed must have at least %d characters", length))
		return false
	}
	return true

}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	if f.Get(field) == "" {
		return false
	}

	return true
}

// Valid returns true if there is no errors, otherwise returns false
func (f *Form) Valid() bool {
	return len(f.Err) == 0
}

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Err.Add(field, "Invalid Email Address")
	}
}
