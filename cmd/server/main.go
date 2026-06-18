package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"staus": "ok"}`)
	})
	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", r)
}
