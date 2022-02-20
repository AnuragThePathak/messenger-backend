package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"messenger-backend/data"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(s *data.Service) http.Handler {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is signup page."))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var user data.User
		json.NewDecoder(r.Body).Decode(&user)

		if issues := dataValidation(user); len(issues) != 0 {
			errorResponse := struct {
				Error []string `json:"error"`
			}{
				Error: issues,
			}

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&errorResponse)
			return
		}

		if checkForDuplicate(s, r, &w, "email", user.Email) {
			return
		}
		if checkForDuplicate(s, r, &w, "username", user.Username) {
			return
		}

		switch hashedBytes, err :=
			bcrypt.GenerateFromPassword([]byte(user.Password), 10); {
		case errors.Is(err, bcrypt.ErrHashTooShort):
			w.Write([]byte("Password too short."))
			return
		case err != nil:
			w.WriteHeader(http.StatusInternalServerError)
			return
		default:
			user.Password = string(hashedBytes)
		}

		switch err := s.CreateAccount(r.Context(), user); {
		case errors.Is(err, context.Canceled), errors.Is(err,
			context.DeadlineExceeded):
			return
		case err != nil:
			w.WriteHeader(http.StatusInternalServerError)
			return
		default:
			w.Write([]byte("Succesfully created account."))
		}
	})

	return r
}

func checkForDuplicate(s *data.Service, r *http.Request, w *http.ResponseWriter,
	credentialType string, credential string) bool {
	switch exists, err := s.IfEmailOrUsernameExists(r.Context(), credentialType,
		credential); {
	case errors.Is(err, context.Canceled), errors.Is(err,
		context.DeadlineExceeded):
		return true
	case err != nil:
		(*w).WriteHeader(http.StatusInternalServerError)
		return true
	case exists:
		message := fmt.Sprintf("%s not available", credentialType)
		(*w).Write([]byte(message))
		return true
	default:
		return false
	}
}
