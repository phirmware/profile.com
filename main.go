package main

import (
	"fmt"
	"net/http"

	"profile.com/controllers"

	"github.com/gorilla/mux"
)

const port = ":8080"

func main() {
	staticC := controllers.NewStatic()
	userC := controllers.NewUser()

	r := mux.NewRouter()
	r.HandleFunc("/", staticC.Home).Methods("GET")
	r.HandleFunc("/signup", userC.New).Methods("GET")
	r.HandleFunc("/signup", userC.Register).Methods("POST")

	fmt.Printf("Listening at port %s", port)
	http.ListenAndServe(port, r)
}
