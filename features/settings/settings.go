package settings

import (
	"fmt"
	"mouji/commons/components"
	"mouji/commons/config"
	"mouji/commons/templates"
	"mouji/features/projects"
	"net/http"
	"net/url"
)

func HandleSettingsPage(w http.ResponseWriter, r *http.Request) {
	allProjects := projects.GetAllProjects()

	renderSettingsPage(w, allProjects)
}

func HandleServerURLPage(w http.ResponseWriter, r *http.Request) {
	isOnboarding := r.URL.Query().Get("is_onboarding") == "true"

	serverURL, err := config.GetConfig("server_url")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serverURLError := ""

	renderServerURLPage(w, isOnboarding, serverURL, serverURLError)
}

func HandleServerURLSubmit(w http.ResponseWriter, r *http.Request) {
	isOnboarding := r.URL.Query().Get("is_onboarding") == "true"

	err := r.ParseForm()
	if err != nil {
		err = fmt.Errorf("error parsing form: %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	serverURL := r.Form.Get("server_url")
	serverURLError := ""

	if !isValidURL(serverURL) {
		serverURLError = "Please enter a valid URL"
		renderServerURLPage(w, isOnboarding, serverURL, serverURLError)
		return
	}

	err = config.SetConfig("server_url", serverURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isOnboarding {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func renderSettingsPage(w http.ResponseWriter, allProjects []projects.ProjectRecord) {
	type templateData struct {
		Navbar               components.Navbar
		Projects             []projects.ProjectRecord
		NewProjectButton     components.Button
		ChangePasswordButton components.Button
		ServerURLButton      components.Button
	}

	tmplData := templateData{
		Navbar:   components.NewNavbar(false),
		Projects: allProjects,
		NewProjectButton: components.Button{
			Text: "New Project",
			Icon: "plus",
			Link: "/projects/new",
		},
		ChangePasswordButton: components.Button{
			Text: "Change Password",
			Icon: "key",
			Link: "/users/me/password",
		},
		ServerURLButton: components.Button{
			Text: "Change Server URL",
			Icon: "server-stack",
			Link: "/settings/server_url",
		},
	}

	templates.Render(w, "settings.html", tmplData)
}

func renderServerURLPage(w http.ResponseWriter, isOnboarding bool, serverURL string, serverURLError string) {
	type templateData struct {
		Navbar         components.Navbar
		IsOnboarding   bool
		ServerURLInput components.Input
		SubmitButton   components.Button
	}

	tmplData := templateData{
		Navbar:       components.NewNavbar(false),
		IsOnboarding: isOnboarding,
		ServerURLInput: components.Input{
			ID:          "server_url",
			Label:       "Server URL",
			Type:        "url",
			Placeholder: "Example: https://www.myserver.com",
			Error:       serverURLError,
			Value:       serverURL,
			Hint:        "Enter the Base URL of the server to correctly generate your tracking snippet",
		},
		SubmitButton: components.Button{
			Text:      "Update",
			IsSubmit:  true,
			IsPrimary: true,
		},
	}

	templates.Render(w, "server_url.html", tmplData)
}

func isValidURL(serverURL string) bool {
	_, err := url.ParseRequestURI(serverURL)
	return err == nil
}
