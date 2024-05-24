package seed

import (
	"encoding/csv"
	"io"
	"math/rand"
	"mouji/features/pageviews"
	"mouji/features/projects"
	"mouji/features/users"
	"os"
)

func SeedUsers() {
	rows := readCSV("./commons/seed/users.csv")
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

func SeedProjects() {
	rows := readCSV("./commons/seed/projects.csv")
	for i, row := range rows {
		if i == 0 {
			continue
		}

		name := row[0]
		baseURL := row[1]

		projects.InsertProject(name, baseURL)
	}
}

func SeedPageViews() {
	projects := projects.GetAllProjects()
	projectNameToIDMap := make(map[string]string)

	for _, project := range projects {
		projectNameToIDMap[project.Name] = project.ProjectID
	}

	rows := readCSV("./commons/seed/pageviews.csv")
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

func readCSV(path string) [][]string {
	file, err := os.Open(path)
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
