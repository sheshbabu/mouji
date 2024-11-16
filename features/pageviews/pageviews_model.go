package pageviews

import (
	"fmt"
	"log/slog"
	"mouji/commons/components"
	"mouji/commons/sqlite"
	"slices"
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

func InsertPageView(record PageViewRecord) error {
	query := "INSERT INTO pageviews (project_id, path, title, referrer) VALUES (?, ?, ?, ?);"

	_, err := sqlite.DB.Exec(query, record.ProjectID, record.Path, record.Title, record.Referrer)
	if err != nil {
		err = fmt.Errorf("error inserting pageview: %w", err)
		slog.Error(err.Error())
		return err
	}

	return nil
}

func GetPaginatedPageViews(projectID string, daterange components.DataRangeType, limit int, offset int) ([]PaginatedPageViewRecord, error) {
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
			AND
			received_at >= DATETIME('now', ?)
		GROUP BY
			path
		ORDER BY
			views DESC
		LIMIT
			?
		OFFSET
			?
	`

	rows, err := sqlite.DB.Query(query, projectID, getDateRangeFilter(daterange), limit, offset)
	if err != nil {
		err = fmt.Errorf("error retrieving pageviews: %w", err)
		slog.Error(err.Error())
		return records, err
	}
	defer rows.Close()

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

func GetPageViewCountsByInterval(projectID string, daterange components.DataRangeType) ([]components.BarChartInputDataPoint, error) {
	var records []components.BarChartInputDataPoint

	query := `
		SELECT
			STRFTIME(?, received_at) AS interval,
			COUNT(*) AS count
		FROM
			pageviews
		WHERE
			project_id = ?
			AND
			received_at >= DATETIME('now', ?)
		GROUP BY
			interval
		ORDER BY
			received_at
	`

	rows, err := sqlite.DB.Query(query, getDateRangeIntervalFormat(daterange), projectID, getDateRangeFilter(daterange))
	if err != nil {
		err = fmt.Errorf("error retrieving pageview counts: %w", err)
		slog.Error(err.Error())
		return records, err
	}
	defer rows.Close()

	for rows.Next() {
		var record components.BarChartInputDataPoint
		err = rows.Scan(&record.Label, &record.Data)
		if err != nil {
			return records, err
		}
		records = append(records, record)
	}

	return records, nil
}

func getDateRangeFilter(daterange components.DataRangeType) string {
	if !slices.Contains(components.DateRangeValues, daterange) {
		daterange = components.DateRangeValues[0]
	}

	switch daterange {
	case "24h":
		return "-24 hours"
	case "1w":
		return "-6 days"
	case "1m":
		return "-1 months"
	case "3m":
		return "-3 months"
	case "1y":
		return "-1 years"
	}

	return "-24 hours"
}

func getDateRangeIntervalFormat(daterange components.DataRangeType) string {
	switch daterange {
	case "24h":
		return "%Y-%m-%d %H"
	case "1w":
		return "%Y-%m-%d"
	case "1m":
		return "%Y-%m-%d"
	case "3m":
		return "%Y-%m-%d"
	case "1y":
		return "%Y-%m"
	}

	return "%Y-%m-%d %H"
}
