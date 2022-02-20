package auth

import (
	"messenger-backend/data"
	"regexp"
	"unicode"
)

func dataValidation(user data.User) []string {
	var issues []string

	if ok, emailIssues := isEmailFormat(user.Email); !ok {
		issues = append(issues, emailIssues...)
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

	if ok, passwordIssue := isPasswordFormat(user.Password); !ok {
		issues = append(issues, passwordIssue...)
		return issues
	}

	if ok, usernameIssues := isUsernameFormat(user.Username); !ok {
		issues = append(issues, usernameIssues...)
		return issues
	}

	return issues
}

func isPasswordFormat(password string) (bool, []string) {
	var issues []string
	
	if len(password) < 6 {
		issues = append(issues, "Password must be of at least 6 characters")
		return false, issues
	}
	if !isASCII(password) {
		issues = append(issues, "Password contains invalid characters.")
		return false, issues
	}
	return true, issues
}

func isEmailFormat(email string) (bool, []string)  {
	var issues []string

	if ok, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`,
		email); !ok {
		issues = append(issues, "Invalid email.")
		return false, issues
	}
	return true, issues
}

func isUsernameFormat(username string) (bool, []string) {
	var issues []string

	if len(username) < 4 {
		issues = append(issues, "Username must be minimum 4 characters.")
		return false ,issues
	}
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username); !ok {
		issues = append(issues,
			"Username can consist of a-z, A-Z, _ only.")
		return false ,issues
	}
	return true, issues
}


func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
