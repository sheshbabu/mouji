package users

import (
	"database/sql"
	"errors"
	"fmt"
	"mouji/commons/components"
	"mouji/commons/session"
	"mouji/commons/templates"
	"net/http"
)

func HandleChangePasswordPage(w http.ResponseWriter, r *http.Request) {
	oldPassword := ""
	newPassword := ""
	oldPasswordError := ""
	newPasswordError := ""
	renderChangePasswordPage(w, oldPassword, newPassword, oldPasswordError, newPasswordError)
}

func HandleChangePasswordSubmit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		err = fmt.Errorf("error parsing form: %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("session_token")
	oldPassword := r.Form.Get("old_password")
	newPassword := r.Form.Get("new_password")

	oldPasswordError := ""
	newPasswordError := ""

	if err != nil && errors.Is(err, http.ErrNoCookie) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sessionID := cookie.Value
	userID, err := session.GetUserID(sessionID)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := GetUserByID(userID)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !IsValidPassword(oldPassword, user.Password) {
		oldPasswordError = "Old password is incorrect"
		renderChangePasswordPage(w, oldPassword, newPassword, oldPasswordError, newPasswordError)
		return
	}

	if oldPassword == newPassword {
		newPasswordError = "New password should be different from old password"
		renderChangePasswordPage(w, oldPassword, newPassword, oldPasswordError, newPasswordError)
		return
	}

	passwordHash, err := HashPassword(newPassword)
	if err != nil {
		err = fmt.Errorf("error hashing password: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = UpdatePassword(userID, passwordHash)
	if err != nil {
		err = fmt.Errorf("error updating password: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func renderChangePasswordPage(w http.ResponseWriter, oldPassword, newPassword, oldPasswordError, newPasswordError string) {
	type templateData struct {
		Navbar           components.Navbar
		OldPasswordInput components.Input
		NewPasswordInput components.Input
		SubmitButton     components.Button
	}

	tmplData := templateData{
		Navbar: components.NewNavbar(false),
		OldPasswordInput: components.Input{
			ID:          "old_password",
			Label:       "Old Password",
			Type:        "password",
			Placeholder: "Enter your old password",
			Error:       oldPasswordError,
			Value:       oldPassword,
		},
		NewPasswordInput: components.Input{
			ID:          "new_password",
			Label:       "New Password",
			Type:        "password",
			Placeholder: "Enter your new password",
			Error:       newPasswordError,
			Value:       newPassword,
		},
		SubmitButton: components.Button{
			Text:      "Update",
			IsSubmit:  true,
			IsPrimary: true,
		},
	}

	templates.Render(w, "password_change.html", tmplData)
}
