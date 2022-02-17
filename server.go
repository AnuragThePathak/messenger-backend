package main

import (
	"fmt"
	"messenger-backend/auth"
	"messenger-backend/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func server(s *models.Service) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "Great Job")
	})

	r.Mount("/signup", auth.Signup(s))
	return r
}
