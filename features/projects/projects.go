package projects

import (
	"fmt"
	"mouji/commons/components"
	"mouji/commons/config"
	"mouji/commons/templates"
	"net/http"
	"net/url"
	"strings"
)

func HandleNewProjectPage(w http.ResponseWriter, r *http.Request) {
	isOnboarding := r.URL.Query().Get("is_onboarding") == "true"
	isNewProject := true

	projectID := ""
	projectName := ""
	siteBaseURL := ""
	projectNameError := ""
	siteBaseURLError := ""

	serverURL, err := config.GetConfig("server_url")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderProjectDetailPage(w, isOnboarding, isNewProject, projectID, projectName, siteBaseURL, serverURL, projectNameError, siteBaseURLError)
}

func HandleEditProjectPage(w http.ResponseWriter, r *http.Request) {
	isOnboarding := false
	isNewProject := false

	projectID := r.PathValue("project_id")
	project, err := getProjectByID(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serverURL, err := config.GetConfig("server_url")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projectNameError := ""
	siteBaseURLError := ""

	renderProjectDetailPage(w, isOnboarding, isNewProject, project.ProjectID, project.Name, project.BaseURL, serverURL, projectNameError, siteBaseURLError)
}

func HandleProjectDetailSubmit(w http.ResponseWriter, r *http.Request) {
	isOnboarding := r.URL.Query().Get("is_onboarding") == "true"
	projectID := r.PathValue("project_id")
	isNewProject := projectID == "new"

	err := r.ParseForm()
	if err != nil {
		err = fmt.Errorf("error parsing form: %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	projectName := r.Form.Get("name")
	siteBaseURL := r.Form.Get("base_url")
	projectNameError := ""
	siteBaseURLError := ""

	serverURL, err := config.GetConfig("server_url")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isValidProjectName(projectName) {
		projectNameError = "Project name should not be empty"
	}

	if !isValidURL(siteBaseURL) {
		siteBaseURLError = "Please enter a valid URL"
	}

	if projectNameError != "" || siteBaseURLError != "" {
		renderProjectDetailPage(w, isOnboarding, isNewProject, projectID, projectName, siteBaseURL, serverURL, projectNameError, siteBaseURLError)
		return
	}

	var project ProjectRecord
	if isNewProject {
		project, err = InsertProject(projectName, siteBaseURL)
	} else {
		project, err = updateProject(projectID, projectName, siteBaseURL)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isOnboarding {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	projectDetailURL := fmt.Sprintf("/projects/%s", project.ProjectID)
	http.Redirect(w, r, projectDetailURL, http.StatusSeeOther)
}

func renderProjectDetailPage(w http.ResponseWriter, isOnboarding bool, isNewProject bool, projectID string, projectName string, siteBaseURL string, serverURL string, projectNameError string, siteBaseURLError string) {
	type templateData struct {
		Navbar               components.Navbar
		IsOnboarding         bool
		IsNewProject         bool
		ProjectID            string
		ProjectNameInput     components.Input
		SiteURLInput         components.Input
		TrackingSnippetInput components.TextArea
		SubmitButton         components.Button
	}

	trackingSnippet := ""
	submitButtonText := "Create"
	if !isNewProject {
		submitButtonText = "Update"
		trackingSnippet = getTrackingSnippet(serverURL, projectID)
	}

	tmplData := templateData{
		Navbar:       components.NewNavbar(false),
		IsOnboarding: isOnboarding,
		IsNewProject: isNewProject,
		ProjectID:    projectID,
		ProjectNameInput: components.Input{
			ID:          "name",
			Label:       "Name",
			Type:        "text",
			Placeholder: "Enter your project name",
			Error:       projectNameError,
			Value:       projectName,
		},
		SiteURLInput: components.Input{
			ID:          "base_url",
			Label:       "Site URL",
			Type:        "url",
			Placeholder: "Example: https://www.blogpost.com",
			Error:       siteBaseURLError,
			Value:       siteBaseURL,
			Hint:        "Enter the base URL of the site associated with this project",
		},
		TrackingSnippetInput: components.TextArea{
			ID:         "tracking_snippet",
			Label:      "Tracking Snippet",
			Content:    trackingSnippet,
			Hint:       "Copy paste this tracking snippet in your site's HTML file at the end of head tag",
			IsDisabled: true,
		},
		SubmitButton: components.Button{
			Text:      submitButtonText,
			IsSubmit:  true,
			IsPrimary: true,
		},
	}

	templates.Render(w, "project_detail.html", tmplData)
}

func isValidProjectName(projectName string) bool {
	return len(strings.TrimSpace(projectName)) > 0
}

func isValidURL(siteBaseURL string) bool {
	_, err := url.ParseRequestURI(siteBaseURL)
	return err == nil
}

func getTrackingSnippet(serverURL string, projectID string) string {
	var snippet = `
<!-- mouji snippet -->
<script>
	(function() {
		var COLLECT_URL = "%s/collect";
		var PROJECT_ID = "%s";
		var GLOBAL_VAR_NAME = "__mouji__";

		window[GLOBAL_VAR_NAME] = {};

		window[GLOBAL_VAR_NAME].sendPageView = function() {
			var path = location.pathname;
			var title = document.title;
			var referrer = document.referrer;

			var url =
				COLLECT_URL +
				"?project_id=" +
				PROJECT_ID +
				"&title=" +
				encodeURIComponent(title) +
				"&path=" +
				encodeURIComponent(path) +
				"&referrer=" +
				encodeURIComponent(referrer);

			var xhr = new XMLHttpRequest();
			xhr.open("GET", url);
			xhr.send();
		};

		window[GLOBAL_VAR_NAME].sendPageView();
	})();
</script>
`
	snippet = strings.TrimSpace(snippet)
	return fmt.Sprintf(snippet, serverURL, projectID)
}
