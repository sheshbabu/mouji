package login

import (
	"database/sql"
	"errors"
	"fmt"
	"mouji/commons/components"
	"mouji/commons/session"
	"mouji/commons/templates"
	"mouji/features/users"
	"net/http"
)

func HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	email := ""
	emailError := ""
	passwordError := ""

	renderLoginForm(w, email, emailError, passwordError)
}

func HandleLoginSubmit(w http.ResponseWriter, r *http.Request) {
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

	user, err := users.GetUserByEmail(email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			emailError = "Email address does not exist"
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if !users.IsValidPassword(password, user.Password) {
		passwordError = "Password is incorrect"
	}

	if emailError != "" || passwordError != "" {
		renderLoginForm(w, email, emailError, passwordError)
		return
	}

	sess, err := session.NewSession(user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.SetSessionCookie(w, sess)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func renderLoginForm(w http.ResponseWriter, email string, emailError string, passwordError string) {
	type templateData struct {
		Navbar        components.Navbar
		EmailInput    components.Input
		PasswordInput components.Input
		SubmitButton  components.Button
	}

	tmplData := templateData{
		Navbar: components.NewNavbar(false),
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
			Text:      "Login",
			Icon:      "arrow-right",
			IsSubmit:  true,
			IsPrimary: true,
		},
	}

	templates.Render(w, "login.html", tmplData)
}
