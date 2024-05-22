package projects

import (
	"fmt"
	"mouji/commons/components"
	"mouji/commons/templates"
	"net/http"
	"net/url"
	"strings"
)

func HandleNewProjectPage(w http.ResponseWriter, r *http.Request) {
	isNewProject := true

	projectID := ""
	projectName := ""
	serverBaseURL := ""
	projectNameError := ""
	serverBaseURLError := ""

	renderProjectDetailPage(w, isNewProject, projectID, projectName, serverBaseURL, projectNameError, serverBaseURLError)
}

func HandleEditProjectPage(w http.ResponseWriter, r *http.Request) {
	isNewProject := false

	projectID := r.PathValue("project_id")
	project, err := getProjectByID(projectID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projectNameError := ""
	serverBaseURLError := ""

	renderProjectDetailPage(w, isNewProject, project.ProjectID, project.Name, project.BaseURL, projectNameError, serverBaseURLError)
}

func HandleProjectDetailSubmit(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("project_id")
	isNewProject := projectID == "new"

	err := r.ParseForm()
	if err != nil {
		err = fmt.Errorf("error parsing form: %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	projectName := r.Form.Get("name")
	serverBaseURL := r.Form.Get("base_url")
	projectNameError := ""
	serverBaseURLError := ""

	if !isValidProjectName(projectName) {
		projectNameError = "Project name should not be empty"
	}

	if !isValidServerBaseURL(serverBaseURL) {
		serverBaseURLError = "Please enter a valid Base URL"
	}

	if projectNameError != "" || serverBaseURLError != "" {
		renderProjectDetailPage(w, isNewProject, projectID, projectName, serverBaseURL, projectNameError, serverBaseURLError)
		return
	}

	var project ProjectRecord
	if isNewProject {
		project, err = insertProject(projectName, serverBaseURL)
	} else {
		project, err = updateProject(projectID, projectName, serverBaseURL)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projectDetailURL := fmt.Sprintf("/projects/%s", project.ProjectID)
	http.Redirect(w, r, projectDetailURL, http.StatusSeeOther)
}

func renderProjectDetailPage(w http.ResponseWriter, isNewProject bool, projectID string, projectName string, serverBaseURL string, projectNameError string, serverBaseURLError string) {
	type templateData struct {
		Navbar               components.Navbar
		IsNewProject         bool
		ProjectID            string
		ProjectNameInput     components.Input
		BaseURLInput         components.Input
		TrackingSnippetInput components.TextArea
		SubmitButton         components.Button
	}

	trackingSnippet := ""
	submitButtonText := "Create"
	if !isNewProject {
		submitButtonText = "Update"
		trackingSnippet = getTrackingSnippet(serverBaseURL, projectID)
	}

	tmplData := templateData{
		Navbar:       components.NewNavbar(false),
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
		BaseURLInput: components.Input{
			ID:          "base_url",
			Label:       "Base URL",
			Type:        "url",
			Placeholder: "Example: https://www.myserver.com",
			Error:       serverBaseURLError,
			Value:       serverBaseURL,
			Hint:        "Enter the Base URL of the server to correctly generate your tracking snippet",
		},
		TrackingSnippetInput: components.TextArea{
			ID:         "tracking_snippet",
			Label:      "Tracking Snippet",
			Content:    trackingSnippet,
			Hint:       "Copy paste this tracking snippet in your HTML file at the end of head tag",
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

func isValidServerBaseURL(serverBaseURL string) bool {
	_, err := url.ParseRequestURI(serverBaseURL)
	return err == nil
}

func getTrackingSnippet(serverBaseURL string, projectID string) string {
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
	return fmt.Sprintf(snippet, serverBaseURL, projectID)
}
