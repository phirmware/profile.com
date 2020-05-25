package controllers

import (
	"net/http"

	"github.com/gorilla/schema"
)

// ParseForm maps the form input to the userform struct
func ParseForm(r *http.Request, user *UserForm) {
	dec := schema.NewDecoder()
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	dec.Decode(user, r.PostForm)
}
