package pageviews

import (
	"fmt"
	"log/slog"
	"mouji/commons/sqlite"
)

type PageViewRecord struct {
	ProjectID string `json:"project_id"`
	Path      string `json:"url"`
	Title     string `json:"title"`
	Referrer  string `json:"referrer"`
}

type PaginatedPageViewRecord struct {
	Title        string
	Path         string
	Views        int
	TotalRecords int
}

func insertPageView(record PageViewRecord) error {
	query := "INSERT INTO pageviews (project_id, path, title, referrer) VALUES (?, ?, ?, ?);"

	_, err := sqlite.DB.Exec(query, record.ProjectID, record.Path, record.Title, record.Referrer)
	if err != nil {
		err = fmt.Errorf("error inserting pageview: %w", err)
		slog.Error(err.Error())
		return err
	}

	return nil
}

func GetPaginatedPageViews(projectID string, limit int, offset int) ([]PaginatedPageViewRecord, error) {
	var records []PaginatedPageViewRecord

	query := `
		SELECT
			title,
			path,
			COUNT(*) AS views,
			COUNT(*) OVER() AS total_rows
		FROM
			pageviews 
		WHERE
			project_id = ?
		GROUP BY
			path
		ORDER BY
			views DESC
		LIMIT
			?
		OFFSET
			?
	`

	rows, err := sqlite.DB.Query(query, projectID, limit, offset)
	if err != nil {
		err = fmt.Errorf("error retrieving pageviews: %w", err)
		slog.Error(err.Error())
		return records, err
	}

	for rows.Next() {
		var record PaginatedPageViewRecord
		err = rows.Scan(&record.Title, &record.Path, &record.Views, &record.TotalRecords)
		if err != nil {
			return records, err
		}
		records = append(records, record)
	}

	return records, nil
}
