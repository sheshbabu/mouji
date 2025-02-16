package home

import (
	"fmt"
	"mouji/commons/components"
	"mouji/commons/config"
	"mouji/commons/templates"
	"mouji/features/pageviews"
	"mouji/features/projects"
	"mouji/features/users"
	"net/http"
	"strconv"
)

type urlState struct {
	selectedProjectID          string
	selectedDateRange          components.DataRangeType
	currentPageViewTableOffset string
}

type pageViewsTable struct {
	Records              []pageviews.PaginatedPageViewRecord
	ShouldShowPagination bool
	Pagination           components.Pagination
}

type pageViewsChart struct {
	TotalCount int
	BarChart   components.BarChart
}

func HandleHomePage(w http.ResponseWriter, r *http.Request) {
	hasUsers := users.HasUsers()
	if !hasUsers {
		http.Redirect(w, r, "/users/new?is_onboarding=true", http.StatusSeeOther)
		return
	}

	server_url, err := config.GetConfig("server_url")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if server_url == "" {
		http.Redirect(w, r, "/settings/server_url?is_onboarding=true", http.StatusSeeOther)
		return
	}

	projects := projects.GetAllProjects()
	if len(projects) == 0 {
		http.Redirect(w, r, "/projects/new?is_onboarding=true", http.StatusSeeOther)
		return
	}

	var state urlState
	state.selectedProjectID = r.URL.Query().Get("project_id")
	state.selectedDateRange = components.DataRangeType(r.URL.Query().Get("daterange"))
	state.currentPageViewTableOffset = r.URL.Query().Get("current_pageview_table_offset")

	if state.selectedProjectID == "" {
		newURL := fmt.Sprintf("/?project_id=%s&daterange=%s&current_pageview_table_offset=%d", projects[0].ProjectID, components.DateRangeValues[0], 0)
		http.Redirect(w, r, newURL, http.StatusSeeOther)
		return
	}

	renderHomePage(w, state, projects)
}

func renderHomePage(w http.ResponseWriter, state urlState, projects []projects.ProjectRecord) {
	type templateData struct {
		Navbar         components.Navbar
		PageViewsChart pageViewsChart
		PageViewsTable pageViewsTable
	}

	navbar := getNavbar(state, projects)

	pageViewsCount, err := pageviews.GetPageViewCountsByInterval(state.selectedProjectID, state.selectedDateRange)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalCount := 0
	barChartInputDataPoints := []components.BarChartInputDataPoint{}
	for _, record := range pageViewsCount {
		barChartInputDataPoint := components.BarChartInputDataPoint{
			Label: record.Interval,
			Data:  record.Count,
		}
		barChartInputDataPoints = append(barChartInputDataPoints, barChartInputDataPoint)
		totalCount = record.TotalCount
	}
	barChart := components.NewBarChart(barChartInputDataPoints)
	chart := pageViewsChart{
		TotalCount: totalCount,
		BarChart:   barChart,
	}

	table, err := getPageViewsTable(state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmplData := templateData{
		Navbar:         navbar,
		PageViewsChart: chart,
		PageViewsTable: table,
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
		option.Link = fmt.Sprintf("/?project_id=%s&daterange=%s", project.ProjectID, state.selectedDateRange)
		option.Value = ""
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
	navbar.ProjectsDropdown.InputName = ""

	navbar.DateRange = getDateRange(state)

	return navbar
}

func getDateRange(state urlState) components.DateRange {
	var daterange components.DateRange

	for _, value := range components.DateRangeValues {
		var option components.DateRangeOption
		option.Name = value
		option.Link = fmt.Sprintf("/?project_id=%s&daterange=%s", state.selectedProjectID, value)
		if value == state.selectedDateRange {
			option.IsSelected = true
		}
		daterange.Options = append(daterange.Options, option)
	}

	return daterange
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

	records, err = pageviews.GetPaginatedPageViews(state.selectedProjectID, state.selectedDateRange, limit, pageViewTableOffset)
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
		table.Pagination.PrevLink = fmt.Sprintf("/?project_id=%s&daterange=%s&current_pageview_table_offset=%d", state.selectedProjectID, state.selectedDateRange, pageViewTableOffset-limit)
	}

	if table.ShouldShowPagination && pageViewTableOffset+limit < table.Pagination.TotalRecords {
		table.Pagination.NextLink = fmt.Sprintf("/?project_id=%s&daterange=%s&current_pageview_table_offset=%d", state.selectedProjectID, state.selectedDateRange, pageViewTableOffset+limit)
	}

	return table, nil
}
