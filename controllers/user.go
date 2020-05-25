package controllers

import (
	"net/http"

	"profile.com/models"

	"profile.com/views"
)

// User defines the shape of the user
type User struct {
	NewView *views.Views
	us      *models.UserService
}

// UserForm defines the shape of the signup form
type UserForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// NewUser returns the user struct
func NewUser() *User {
	return &User{
		NewView: views.NewView("bootstrap", "user/new"),
		us:      models.NewUserService(),
	}
}

// New handles route /signup
func (u *User) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

// Register creates a new user in the database
func (u *User) Register(w http.ResponseWriter, r *http.Request) {
	var form UserForm
	var data views.Data
	ParseForm(r, &form)
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.UserDB.Create(&user); err != nil {
		data.SetAlert(&data, views.ErrLevelDanger, err)
		u.NewView.Render(w, data)
	}

}
