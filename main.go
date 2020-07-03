package main

import (
	"log"
	"net/http"

	"github.com/BrandenSoropia/dodo-getaways-backend/controllers"
	"github.com/BrandenSoropia/dodo-getaways-backend/db"
	"github.com/gorilla/mux"
)

func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

func main() {
	db.Connect()

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", YourHandler)
	r.HandleFunc("/islands", controllers.GetIslands).Methods("GET")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
