package pageviews

import (
	"database/sql"
	"fmt"
	"log/slog"
	"mouji/commons/components"
	"mouji/commons/sqlite"
	"slices"
)

type PageViewRecord struct {
	ProjectID   string
	Path        string
	Title       string
	Referrer    string
	VisitorHash string
	UserAgent   string
}

type PaginatedPageViewRecord struct {
	Title        string
	Path         string
	Views        int
	TotalRecords int
}

type PageViewCountRecord struct {
	Interval   string
	Count      int
	TotalCount int
}

func InsertPageView(record PageViewRecord) error {
	query := "INSERT INTO pageviews (project_id, path, title, referrer, visitor_hash, user_agent) VALUES (?, ?, ?, ?, ?, ?);"

	_, err := sqlite.DB.Exec(query, record.ProjectID, record.Path, record.Title, record.Referrer, record.VisitorHash, record.UserAgent)
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

func GetPageViewCountsByInterval(projectID string, daterange components.DataRangeType) ([]PageViewCountRecord, error) {
	var records []PageViewCountRecord
	var rows *sql.Rows
	var err error

	if daterange == "24h" {
		query := `
			SELECT
				-- https://stackoverflow.com/a/33116186
				STRFTIME('%d ', received_at) || SUBSTR('--JanFebMarAprMayJunJulAugSepOctNovDec', STRFTIME('%m', received_at) * 3, 3) || ', ' ||
					CASE 
						WHEN STRFTIME('%H', received_at) = '00' THEN '12 - 01 AM'
						WHEN STRFTIME('%H', received_at) = '12' THEN '12 - 01 PM'
						ELSE
							STRFTIME('%I', received_at) || ' - ' || STRFTIME('%I', DATETIME(received_at, '+1 hour')) || 
								CASE 
									WHEN STRFTIME('%p', received_at) = 'AM' AND STRFTIME('%H', received_at) = '11' THEN ' PM'
									WHEN STRFTIME('%p', received_at) = 'PM' AND STRFTIME('%H', received_at) = '23' THEN ' AM'
									ELSE STRFTIME(' %p', received_at)
								END
					END AS interval,
				COUNT(*) AS count, 
				SUM(COUNT(*)) OVER() AS total_count
			FROM
				pageviews
			WHERE
				project_id = ?
				AND
				received_at >= DATETIME('now', '-24 hours')
			GROUP BY
				STRFTIME('%Y-%m-%d %H', received_at)
			ORDER BY
				STRFTIME('%Y-%m-%d %H', received_at);
		`
		rows, err = sqlite.DB.Query(query, projectID)
	} else if daterange == "1w" || daterange == "1m" || daterange == "3m" {
		query := `
			SELECT
				STRFTIME('%d ', received_at) || SUBSTR('--JanFebMarAprMayJunJulAugSepOctNovDec', STRFTIME('%m', received_at) * 3, 3) AS interval,
				COUNT(*) AS count, 
				SUM(COUNT(*)) OVER() AS total_count
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
		rows, err = sqlite.DB.Query(query, projectID, getDateRangeFilter(daterange))
	} else {
		query := `
			SELECT
				STRFTIME('%Y ', received_at) || SUBSTR('--JanFebMarAprMayJunJulAugSepOctNovDec', STRFTIME('%m', received_at) * 3, 3) AS interval,
				COUNT(*) AS count, 
				SUM(COUNT(*)) OVER() AS total_count
			FROM
				pageviews
			WHERE
				project_id = ?
				AND
				received_at >= DATETIME('now', '-1 years')
			GROUP BY
				interval
			ORDER BY
				received_at
		`
		rows, err = sqlite.DB.Query(query, projectID, getDateRangeFilter(daterange))
	}
	if err != nil {
		err = fmt.Errorf("error retrieving pageview counts: %w", err)
		slog.Error(err.Error())
		return records, err
	}
	defer rows.Close()

	for rows.Next() {
		var record PageViewCountRecord
		err = rows.Scan(&record.Interval, &record.Count, &record.TotalCount)
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
