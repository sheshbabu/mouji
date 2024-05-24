package seed

import (
	"embed"
	"encoding/csv"
	"io"
	"math/rand"
	"mouji/features/pageviews"
	"mouji/features/projects"
	"mouji/features/users"
)

func SeedUsers(resources embed.FS) {
	rows := readCSV(resources, "commons/seed/users.csv")
	for i, row := range rows {
		if i == 0 {
			continue
		}

		email := row[0]
		passwordHash, _ := users.HashPassword(row[1])
		isAdmin := row[2] == "true"

		users.InsertUser(email, passwordHash, isAdmin)
	}
}

func SeedProjects(resources embed.FS) {
	rows := readCSV(resources, "commons/seed/projects.csv")
	for i, row := range rows {
		if i == 0 {
			continue
		}

		name := row[0]
		baseURL := row[1]

		projects.InsertProject(name, baseURL)
	}
}

func SeedPageViews(resources embed.FS) {
	projects := projects.GetAllProjects()
	projectNameToIDMap := make(map[string]string)

	for _, project := range projects {
		projectNameToIDMap[project.Name] = project.ProjectID
	}

	rows := readCSV(resources, "commons/seed/pageviews.csv")
	for i, row := range rows {
		if i == 0 {
			continue
		}

		randomViewCount := rand.Intn(50)
		for i := 0; i <= randomViewCount; i++ {
			pageviews.InsertPageView(pageviews.PageViewRecord{
				ProjectID: projectNameToIDMap[row[0]],
				Path:      row[1],
				Title:     row[2],
				Referrer:  row[3],
			})
		}
	}
}

func readCSV(resources embed.FS, path string) [][]string {
	file, err := embed.FS.Open(resources, path)

	defer file.Close()

	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(file)

	var rows [][]string
	for {
		row, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		rows = append(rows, row)
	}

	return rows
}
