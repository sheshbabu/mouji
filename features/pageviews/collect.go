package pageviews

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func HandleCollect(w http.ResponseWriter, r *http.Request) {
	var record PageViewRecord

	projectID := r.URL.Query().Get("project_id")
	path := r.URL.Query().Get("path")
	title := r.URL.Query().Get("title")
	referrer := r.URL.Query().Get("referrer")
	userAgent := r.Header.Get("User-Agent")
	ipAddress := r.RemoteAddr

	visitorHash := generateVisitorHash(projectID, ipAddress, userAgent)

	normalizedPath, err := normalizePath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record = PageViewRecord{
		ProjectID:   projectID,
		Path:        normalizedPath,
		Title:       title,
		Referrer:    referrer,
		VisitorHash: visitorHash,
		UserAgent:   userAgent,
	}

	err = InsertPageView(record)
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

// Generates a transient visitor hash that rotates daily
// hash(daily_salt + website_domain + ip_address + user_agent)
// https://news.ycombinator.com/item?id=24696768
func generateVisitorHash(projectID string, ipAddress string, userAgent string) string {

	dailySalt := time.Now().Format("2006-01-02")

	// https://gobyexample.com/sha256-hashes
	hash := sha256.New()
	hash.Write([]byte(dailySalt + projectID + ipAddress + userAgent))
	visitorHash := fmt.Sprintf("%x", hash.Sum(nil))

	return visitorHash
}
