package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/system/init/{key}", Initialize) // ----> To request all groceries
	r.HandleFunc("/internal/note", Note)

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
		}),
		handlers.MaxAge(3600),
	)

	handler := corsMiddleware(r)

	log.Fatal(http.ListenAndServe(":8088", handler))
	//http.Handle("/", r)
}
