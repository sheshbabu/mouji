package pageviews

import (
	"net/http"
	"net/url"
	"strings"
)

type PageViewPayload struct {
	ProjectID string `json:"project_id"`
	Url       string `json:"url"`
	Title     string `json:"title"`
	Referrer  string `json:"referrer"`
}

func HandleCollect(w http.ResponseWriter, r *http.Request) {
	var payload PageViewPayload
	var record PageViewRecord

	payload.ProjectID = r.URL.Query().Get("project_id")
	payload.Url = r.URL.Query().Get("path")
	payload.Title = r.URL.Query().Get("title")
	payload.Referrer = r.URL.Query().Get("referrer")

	normalizedPath, err := normalizePath(payload.Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record = PageViewRecord{
		ProjectID: payload.ProjectID,
		Path:      normalizedPath,
		Title:     payload.Title,
		Referrer:  payload.Referrer,
	}

	err = insertPageView(record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func normalizePath(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	defaultPath := "/"

	if err != nil {
		return defaultPath, err
	}

	var normalizedPath string

	if strings.TrimSpace(parsedURL.Path) == "" {
		normalizedPath = defaultPath
	} else {
		normalizedPath = parsedURL.Path
	}

	return normalizedPath, nil
}
