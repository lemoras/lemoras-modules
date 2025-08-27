package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	// r.Use(rateMiddleware)

	r.HandleFunc("/system/init/{key}", Initialize) // ----> To request all groceries
	r.HandleFunc("/services/note", Note)
	r.HandleFunc("/services/drive", Drive)

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

// func rateMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println(r.Method, r.URL.Path)

// 		// if isOk := u.RateTokenhandler(w, r); !isOk {
// 		// 	return
// 		// }

// 		w.Header().Set("Content-Type", "application/json")

// 		next.ServeHTTP(w, r)
// 	})
// }
