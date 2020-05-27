package controllers

import (
	"fmt"
	"net/http"

	"profile.com/models"

	"profile.com/views"
)

// User defines the shape of the user
type User struct {
	NewView             *views.Views
	CompleteProfileView *views.Views
	DashboardView       *views.Views
	us                  *models.UserService
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

// NewUser returns the user struct
func NewUser(connectionString string) *User {
	us, err := models.NewUserService(connectionString)
	if err != nil {
		panic(err)
	}
	return &User{
		NewView:             views.NewView("bootstrap", "user/new"),
		CompleteProfileView: views.NewView("bootstrap", "user/profile"),
		DashboardView:       views.NewView("bootstrap", "user/dashboard"),
		us:                  us,
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
		return
	}
	fmt.Println(user, "After the create method")
	if err := u.signIn(w, &user); err != nil {
		u.NewView.Render(w, nil)
		return
	}
	uri := fmt.Sprintf("/complete-profile?email=%s", user.Email)
	http.Redirect(w, r, uri, http.StatusFound)
}

// CompleteProfile renders the page to complete profile
func (u *User) CompleteProfile(w http.ResponseWriter, r *http.Request) {
	email := FromQuery(r, "email")
	u.CompleteProfileView.Render(w, email)
}

// Profile completes the user profile
func (u *User) Profile(w http.ResponseWriter, r *http.Request) {
	email := FromQuery(r, "email")
	var form completeForm
	var data views.Data

	ParseForm(r, &form)
	// skills := strings.Split(form.Skills, ",")

	user, err := u.us.ByEmail(email)
	if err != nil {
		data.SetAlert(&data, views.ErrLevelDanger, models.ErrInternalServerError)
		u.CompleteProfileView.Render(w, data)
		return
	}

	user.Skills = form.Skills
	user.Summary = form.Summary
	user.Title = form.Title

	err = u.us.Update(user)
	if err != nil {
		data.SetAlert(&data, views.ErrLevelDanger, models.ErrInternalServerError)
		u.CompleteProfileView.Render(w, data)
		return
	}

	if err := u.signIn(w, user); err != nil {
		u.CompleteProfileView.Render(w, data)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

// Dashboard renders the dashboard page
func (u *User) Dashboard(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Redirect(w, r, "/signup", http.StatusFound)
		return
	}
	user, err := u.us.ByRemember(cookie.Value)
	fmt.Printf("%+v\n", user)
	if err != nil {
		http.Redirect(w, r, "/signup", http.StatusFound)
		return
	}
	u.DashboardView.Render(w, user)
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

// AutoMigrate automigrate creates the table in the db
func (u *User) AutoMigrate() {
	u.us.AutoMigrate()
}
