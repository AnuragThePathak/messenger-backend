package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"unicode"

	"github.com/go-chi/chi/v5"
)

func Signup() http.Handler {
	r := chi.NewRouter()

	r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "This is signup page.")
	})

	r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		var user user
		json.NewDecoder(r.Body).Decode(&user)

		issues := dataValidation(user)
		if len(issues) != 0 {
			errorResponse := struct {
				Error []string `json:"error"`
			}{
				Error: issues,
			}

			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(&errorResponse)
			return
		}
	})

	return r
}

func dataValidation(user user) []string {
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