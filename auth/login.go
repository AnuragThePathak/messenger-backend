package auth

import (
	"context"
	"encoding/json"
	"errors"
	"messenger-backend/data"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(s *data.Service) http.Handler {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is login page"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var user struct {
			Credential string `json:"credential"`
			Password   string `json:"password"`
		}
		var credentialType string

		json.NewDecoder(r.Body).Decode(&user)

		if ok, _ := isUsernameFormat(user.Credential); ok {
			credentialType = "username"
		} else if ok, _ := isEmailFormat(user.Credential); ok {
			credentialType = "email"
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if ok, _ := isPasswordFormat(user.Password); !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		hash, err := s.GetHashByCredential(r.Context(), credentialType,
			user.Credential)
		switch {
		case errors.Is(err, context.Canceled), errors.Is(err,
			context.DeadlineExceeded):
			return
		case errors.Is(err, pgx.ErrNoRows):
			w.Write([]byte("Couldn't find your account."))
			return
		case err != nil:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if verifyPassword(user.Password, hash) {
			w.WriteHeader(http.StatusAccepted)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})

	return r
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
