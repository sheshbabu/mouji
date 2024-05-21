package settings

import (
	"mouji/commons/components"
	"mouji/commons/templates"
	"mouji/features/projects"
	"net/http"
)

func HandleSettingsPage(w http.ResponseWriter, r *http.Request) {
	allProjects := projects.GetAllProjects()

	renderSettingsPage(w, allProjects)
}

func renderSettingsPage(w http.ResponseWriter, allProjects []projects.ProjectRecord) {
	type templateData struct {
		Navbar           components.Navbar
		Projects         []projects.ProjectRecord
		NewProjectButton components.Button
	}

	tmplData := templateData{
		Navbar:   components.NewNavbar(false),
		Projects: allProjects,
		NewProjectButton: components.Button{
			Text: "New Project",
			Icon: "plus",
			Link: "/projects/new",
		},
	}

	templates.Render(w, "settings.html", tmplData)
}
