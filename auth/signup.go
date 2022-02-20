package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"unicode"

	"messenger-backend/models"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(s *models.Service) http.Handler {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "This is signup page.")
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var user models.User
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

func checkForDuplicate(s *models.Service, r *http.Request, w *http.ResponseWriter,
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

func dataValidation(user models.User) []string {
	var issues []string

	if ok, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`,
		user.Email); !ok {
		issues = append(issues, "Invalid email.")
		return issues
	}

	if len(user.Name) < 4 {
		issues = append(issues, "Name must be of at least 4 characters.")
		return issues
	}
	if ok, _ := regexp.MatchString("^[a-zA-Z ]+$", user.Name); !ok {
		issues = append(issues, "Invalid name.")
		return issues
	}

	if len(user.Password) < 6 {
		issues = append(issues, "Password must be of at least 6 characters")
		return issues
	}
	if !isASCII(user.Password) {
		issues = append(issues, "Password contains invalid characters.")
		return issues
	}

	if len(user.Username) < 4 {
		issues = append(issues, "Username must be minimum 4 characters.")
		return issues
	}
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", user.Username); !ok {
		issues = append(issues,
			"Username can consist of a-z, A-Z, _ only.")
		return issues
	}

	return issues
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
