package projects

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"mouji/commons/sqlite"
)

type ProjectRecord struct {
	ProjectID string
	Name      string
	BaseURL   string
}

func HasProjects() bool {
	query := "SELECT project_id FROM projects LIMIT 1"
	var projectID string

	row := sqlite.DB.QueryRow(query)
	err := row.Scan(&projectID)
	if errors.Is(err, sql.ErrNoRows) {
		return false
	}
	if err != nil {
		err = fmt.Errorf("error retrieving projects: %w", err)
		panic(err)
	}

	return true
}

func GetAllProjects() []ProjectRecord {
	var projects []ProjectRecord
	query := "SELECT project_id, name, base_url FROM projects ORDER BY created_at DESC"

	rows, err := sqlite.DB.Query(query)
	if err != nil {
		err = fmt.Errorf("error retrieving projects: %w", err)
		panic(err)
	}

	for rows.Next() {
		var project ProjectRecord
		err = rows.Scan(&project.ProjectID, &project.Name, &project.BaseURL)
		if err != nil {
			err = fmt.Errorf("error retrieving projects: %w", err)
			panic(err)
		}
		projects = append(projects, project)
	}

	return projects
}

func getProjectByID(projectID string) (ProjectRecord, error) {
	var project ProjectRecord

	query := "SELECT project_id, name, base_url FROM projects where project_id = ?"

	row := sqlite.DB.QueryRow(query, projectID)
	err := row.Scan(&project.ProjectID, &project.Name, &project.BaseURL)
	if err != nil {
		err = fmt.Errorf("error retrieving project: %w", err)
		slog.Error(err.Error())
		return project, err
	}

	return project, nil
}

func insertProject(projectName string, serverBaseURL string) (ProjectRecord, error) {
	var project ProjectRecord

	query := `
		INSERT INTO projects (
			project_id,
			name,
			base_url
		)
		VALUES(
			LOWER(HEX(RANDOMBLOB (16))),
			?,
			?
		)
		RETURNING
			project_id,
			name,
			base_url`

	row := sqlite.DB.QueryRow(query, projectName, serverBaseURL)
	err := row.Scan(&project.ProjectID, &project.Name, &project.BaseURL)
	if err != nil {
		err = fmt.Errorf("error inserting project: %w", err)
		slog.Error(err.Error())
		return project, err
	}

	return project, nil
}

func updateProject(projectID string, projectName string, serverBaseURL string) (ProjectRecord, error) {
	var project ProjectRecord

	query := `
		UPDATE projects 
		SET
			name = ?,
			base_url = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			project_id = ?
		RETURNING
			project_id,
			name,
			base_url
	`

	row := sqlite.DB.QueryRow(query, projectName, serverBaseURL, projectID)
	err := row.Scan(&project.ProjectID, &project.Name, &project.BaseURL)
	if err != nil {
		err = fmt.Errorf("error updating project: %w", err)
		slog.Error(err.Error())
		return project, err
	}

	return project, nil
}
