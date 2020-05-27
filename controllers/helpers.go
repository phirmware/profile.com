package controllers

import (
	"net/http"

	"github.com/gorilla/schema"
)

// ParseForm maps the form input to the userform struct
func ParseForm(r *http.Request, form interface{}) {
	dec := schema.NewDecoder()
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	dec.Decode(form, r.PostForm)
}

// FromQuery gets the value of the query from key
func FromQuery(r *http.Request, key string) string {
	return r.FormValue(key)
}

// // GetCookies gets the cookies from the request
// func GetCookies(r *http.Request, name string) (*http.Cookie, error) {
// 	cookie, err := r.Cookie(name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return cookie, nil
// }
