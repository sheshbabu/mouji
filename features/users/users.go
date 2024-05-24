package users

import (
	"fmt"
	"mouji/commons/components"
	"mouji/commons/session"
	"mouji/commons/templates"
	"net/http"
	"net/mail"
	"strings"
)

func HandleNewUserPage(w http.ResponseWriter, r *http.Request) {
	hasUsers := HasUsers()

	email := ""
	emailError := ""
	passwordError := ""

	renderNewUserPage(w, hasUsers, email, emailError, passwordError)
}

func HandleNewUserSubmit(w http.ResponseWriter, r *http.Request) {
	hasUsers := HasUsers()

	err := r.ParseForm()
	if err != nil {
		err = fmt.Errorf("error parsing form: %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	emailError := ""
	passwordError := ""

	if !isValidEmail(email) {
		emailError = "Please enter a valid email address"
	}

	if !isValidPassword(password) {
		passwordError = "Password should not be empty"
	}

	if emailError != "" || passwordError != "" {
		renderNewUserPage(w, hasUsers, email, emailError, passwordError)
		return
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		err = fmt.Errorf("error hashing password: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If there are no users, then the first user becomes an admin
	isAdmin := !hasUsers

	err = InsertUser(email, passwordHash, isAdmin)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !hasUsers {
		sess, err := session.NewSession()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.SetSessionCookie(w, sess)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func renderNewUserPage(w http.ResponseWriter, hasUsers bool, email string, emailError string, passwordError string) {
	type templateData struct {
		Navbar        components.Navbar
		HasUsers      bool
		EmailInput    components.Input
		PasswordInput components.Input
		SubmitButton  components.Button
	}

	submitButtonText := "Create"
	submitButtonIcon := ""
	if !hasUsers {
		submitButtonText = "Continue"
		submitButtonIcon = "arrow-right"
	}

	tmplData := templateData{
		Navbar:   components.NewNavbar(false),
		HasUsers: hasUsers,
		EmailInput: components.Input{
			ID:          "email",
			Label:       "Email",
			Type:        "email",
			Placeholder: "Enter your email address",
			Error:       emailError,
			Value:       email,
		},
		PasswordInput: components.Input{
			ID:          "password",
			Label:       "Password",
			Type:        "password",
			Placeholder: "Enter your password",
			Error:       passwordError,
		},
		SubmitButton: components.Button{
			Text:      submitButtonText,
			Icon:      submitButtonIcon,
			IsSubmit:  true,
			IsPrimary: true,
		},
	}

	templates.Render(w, "user_detail.html", tmplData)
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidPassword(password string) bool {
	return len(strings.TrimSpace(password)) > 0
}
