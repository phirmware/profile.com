package main

import (
	"fmt"
	"net/http"

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

	staticC := controllers.NewStatic()
	userC := controllers.NewUser(psqlInfo)
	userC.AutoMigrate()

	r := mux.NewRouter()
	r.HandleFunc("/", staticC.Home).Methods("GET")
	r.HandleFunc("/signup", userC.New).Methods("GET")
	r.HandleFunc("/signup", userC.Register).Methods("POST")
	r.HandleFunc("/complete-profile", userC.CompleteProfile).Queries("email", "{email}").Methods("GET")
	r.HandleFunc("/complete-profile", userC.Profile).Queries("email", "{email}").Methods("POST")
	r.HandleFunc("/dashboard", userC.Dashboard).Methods("GET")

	fmt.Printf("Listening at port %s", serverPort)
	http.ListenAndServe(serverPort, r)
}
