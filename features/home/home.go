package home

import (
	"fmt"
	"mouji/commons/components"
	"mouji/commons/templates"
	"mouji/features/pageviews"
	"mouji/features/projects"
	"mouji/features/users"
	"net/http"
	"strconv"
)

type urlState struct {
	selectedProjectID          string
	currentPageViewTableOffset string
}

type pageViewsTable struct {
	Records              []pageviews.PaginatedPageViewRecord
	ShouldShowPagination bool
	Pagination           components.Pagination
}

func HandleHomePage(w http.ResponseWriter, r *http.Request) {
	hasUsers := users.HasUsers()
	projects := projects.GetAllProjects()

	if !hasUsers {
		http.Redirect(w, r, "/users/new", http.StatusSeeOther)
		return
	}

	if len(projects) == 0 {
		http.Redirect(w, r, "/projects/new", http.StatusSeeOther)
		return
	}

	var state urlState
	state.selectedProjectID = r.URL.Query().Get("project_id")
	state.currentPageViewTableOffset = r.URL.Query().Get("current_pageview_table_offset")

	if state.selectedProjectID == "" {
		newURL := fmt.Sprintf("/?project_id=%s&current_pageview_table_offset=%d", projects[0].ProjectID, 0)
		http.Redirect(w, r, newURL, http.StatusSeeOther)
		return
	}

	renderHomePage(w, state, projects)
}

func renderHomePage(w http.ResponseWriter, state urlState, projects []projects.ProjectRecord) {
	type templateData struct {
		Navbar         components.Navbar
		PageViewsTable pageViewsTable
	}

	navbar := getNavbar(state, projects)
	pageviews, err := getPageViewsTable(state)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmplData := templateData{
		Navbar:         navbar,
		PageViewsTable: pageviews,
	}

	templates.Render(w, "home.html", tmplData)
}

func getNavbar(state urlState, projects []projects.ProjectRecord) components.Navbar {
	navbar := components.NewNavbar(true)
	var allOptions []components.DropdownOption
	var selectedOption components.DropdownOption
	for _, project := range projects {
		var option components.DropdownOption
		option.Name = project.Name
		option.Link = fmt.Sprintf("/?project_id=%s", project.ProjectID)
		allOptions = append(allOptions, option)
		if project.ProjectID == state.selectedProjectID {
			selectedOption = option
		}
	}
	if state.selectedProjectID == "" {
		selectedOption = allOptions[0]
	}
	navbar.ProjectsDropdown.SelectedOption = selectedOption
	navbar.ProjectsDropdown.AllOptions = allOptions
	return navbar
}

func getPageViewsTable(state urlState) (pageViewsTable, error) {
	var records []pageviews.PaginatedPageViewRecord
	limit := 10

	pageViewTableOffset, err := strconv.Atoi(state.currentPageViewTableOffset)
	if err != nil {
		pageViewTableOffset = 0
	}

	table := pageViewsTable{
		Records:              records,
		ShouldShowPagination: false,
		Pagination: components.Pagination{
			PageStartRecord: pageViewTableOffset + 1,
			PageEndRecord:   0,
			TotalRecords:    0,
			PrevLink:        "",
			NextLink:        "",
		},
	}

	records, err = pageviews.GetPaginatedPageViews(state.selectedProjectID, limit, pageViewTableOffset)
	if err != nil {
		return table, err
	}

	if len(records) > 0 {
		table.Records = records
		table.Pagination.TotalRecords = records[0].TotalRecords
		table.Pagination.PageStartRecord = pageViewTableOffset + 1
		table.Pagination.PageEndRecord = pageViewTableOffset + len(records)
		table.ShouldShowPagination = records[0].TotalRecords > limit
	}

	if table.ShouldShowPagination && pageViewTableOffset != 0 {
		table.Pagination.PrevLink = fmt.Sprintf("/?project_id=%s&current_pageview_table_offset=%d", state.selectedProjectID, pageViewTableOffset-limit)
	}

	if table.ShouldShowPagination && pageViewTableOffset+limit < table.Pagination.TotalRecords {
		table.Pagination.NextLink = fmt.Sprintf("/?project_id=%s&current_pageview_table_offset=%d", state.selectedProjectID, pageViewTableOffset+limit)
	}

	return table, nil
}
