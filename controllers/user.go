package controllers

import (
	"fmt"
	"net/http"

	"profile.com/context"
	"profile.com/models"

	"profile.com/views"
)

// User defines the shape of the user
type User struct {
	NewView             *views.Views
	LoginView           *views.Views
	CompleteProfileView *views.Views
	DashboardView       *views.Views
	us                  models.UserService
}

// UserForm defines the shape of the signup form
type UserForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type completeForm struct {
	Title   string `schema:"title"`
	Summary string `schema:"summary"`
	Skills  string `schema:"skills"`
}

type loginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// NewUser returns the user struct
func NewUser(us models.UserService) *User {
	return &User{
		NewView:             views.NewView("bootstrap", "user/new"),
		LoginView:           views.NewView("bootstrap", "user/login"),
		CompleteProfileView: views.NewView("bootstrap", "user/profile"),
		DashboardView:       views.NewView("bootstrap", "user/dashboard"),
		us:                  us,
	}
}

// New handles route /signup
func (u *User) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, r, nil)
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
	if err := u.us.Create(&user); err != nil {
		data.SetAlert(views.ErrLevelDanger, err)
		u.NewView.Render(w, r, data)
		return
	}
	if err := u.signIn(w, &user); err != nil {
		u.NewView.Render(w, r, nil)
		return
	}
	uri := fmt.Sprintf("/complete-profile?email=%s", user.Email)
	http.Redirect(w, r, uri, http.StatusFound)
}

// CompleteProfile renders the page to complete profile
func (u *User) CompleteProfile(w http.ResponseWriter, r *http.Request) {
	email := FromQuery(r, "email")
	u.CompleteProfileView.Render(w, r, email)
}

// Profile completes the user profile
func (u *User) Profile(w http.ResponseWriter, r *http.Request) {
	var form completeForm
	var data views.Data

	ParseForm(r, &form)
	// skills := strings.Split(form.Skills, ",")
	user := context.GetUserFromContext(r.Context())

	user.Skills = form.Skills
	user.Summary = form.Summary
	user.Title = form.Title

	if err := u.us.Update(user); err != nil {
		data.SetAlert(views.ErrLevelDanger, models.ErrInternalServerError)
		u.CompleteProfileView.Render(w, r, data)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

// Login renders the login view
func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	u.LoginView.Render(w, r, nil)
}

// HandleLogin logs in the user
func (u *User) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var form loginForm
	var data views.Data
	ParseForm(r, &form)

	user := &models.User{
		Email:    form.Email,
		Password: form.Password,
	}
	foundUser, err := u.us.Authenticate(user)
	if err != nil {
		data.SetAlert(views.ErrLevelDanger, err)
		u.LoginView.Render(w, r, data)
		return
	}
	if err := u.signIn(w, foundUser); err != nil {
		data.SetAlert(views.ErrLevelDanger, err)
		u.LoginView.Render(w, r, data)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

// Dashboard renders the dashboard page
func (u *User) Dashboard(w http.ResponseWriter, r *http.Request) {
	user := context.GetUserFromContext(r.Context())
	u.DashboardView.Render(w, r, user)
}

// Users gets all users
func (u *User) Users(w http.ResponseWriter, r *http.Request) {
	users, err := u.us.All()
	// user, err := u.us.ByEmail("chibuzor.ojukwu@gmail.com")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "%+v\n", users)
}

func (u *User) signIn(w http.ResponseWriter, user *models.User) error {
	cookie := &http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return nil
}
