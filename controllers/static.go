package controllers

import (
	"net/http"

	"profile.com/views"
)

// Static defines the shape of the static controller
type Static struct {
	HomeView *views.Views
}

// NewStatic returns the Static struct
func NewStatic() *Static {
	return &Static{
		HomeView: views.NewView("bootstrap", "static/home"),
	}
}

// Home handles the / GET request
func (s Static) Home(w http.ResponseWriter, r *http.Request) {
	s.HomeView.Render(w, nil)
}
