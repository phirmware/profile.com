package main

import (
	"fmt"
	"net/http"

	"profile.com/middleware"

	"profile.com/models"

	"profile.com/controllers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const serverPort = ":8080"

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "profile_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		" dbname=%s sslmode=disable",
		host, port, user, dbname)

	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	services.AutoMigrate()

	staticC := controllers.NewStatic()
	userC := controllers.NewUser(services.User)

	requireUserMW := middleware.NewRequireUserMiddleWare(services.User)
	userMW := middleware.NewUserMiddleWare(services.User)
	dashboard := requireUserMW.ApplyFn(userC.Dashboard)
	completeProfile := requireUserMW.ApplyFn(userC.CompleteProfile)
	profile := requireUserMW.ApplyFn(userC.Profile)

	r := mux.NewRouter()
	r.HandleFunc("/", staticC.Home).Methods("GET")
	r.HandleFunc("/signup", userC.New).Methods("GET")
	r.HandleFunc("/signup", userC.Register).Methods("POST")
	r.HandleFunc("/login", userC.Login).Methods("GET")
	r.HandleFunc("/login", userC.HandleLogin).Methods("POST")
	r.HandleFunc("/complete-profile", completeProfile).Queries("email", "{email}").Methods("GET")
	r.HandleFunc("/complete-profile", profile).Queries("email", "{email}").Methods("POST")
	r.HandleFunc("/dashboard", dashboard).Methods("GET")
	r.HandleFunc("/users", userC.Users).Methods("GET")

	fmt.Printf("Listening at port %s", serverPort)
	http.ListenAndServe(serverPort, userMW.Apply(r))
}
