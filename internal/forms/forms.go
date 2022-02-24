package forms

import (
	"net/http"
	"net/mail"
	"net/url"
	"strings"
)

// Form creates a custom form struct embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}

func (f *Form) Required(fields ...string) bool {
	for _, field := range fields {
		x := f.Get(field)
		if strings.TrimSpace(x) == "" {
			f.Errors.Add(field, "This field cannot be blank")
			return false
		}
	}
	return true
}

func (f *Form) EmailLength(field string) bool {
	_, err := mail.ParseAddress(f.Get(field))
	if err != nil {
		f.Errors.Add(field, "Please add valid email address")
		return false
	}
	return true
}
