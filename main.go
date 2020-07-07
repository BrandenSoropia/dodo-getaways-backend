package main

import (
	"log"
	"net/http"

	"github.com/BrandenSoropia/dodo-getaways-backend/controllers"
	"github.com/BrandenSoropia/dodo-getaways-backend/db"
	"github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {
	db.Connect()

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	// Routes consist of a path and a handler function.
	r.HandleFunc("/get-islands", controllers.GetIslands).Methods("GET")
	r.HandleFunc("/get-island", controllers.GetIsland).Methods("POST")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
