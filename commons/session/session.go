package session

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"mouji/commons/sqlite"
	"net/http"
	"time"
)

var sessionLength = time.Hour * 24 * 7 // 7 days
var expiresAtFormat = "2006-01-02T15:04:05-07:00"

type SessionRecord struct {
	SessionID string
	ExpiresAt time.Time
}

func NewSession() (SessionRecord, error) {
	expiresAt := time.Now().Add(sessionLength)

	var session SessionRecord
	var expiresAtString string

	query := "INSERT INTO sessions (session_id, expires_at) VALUES (LOWER(HEX(RANDOMBLOB (16))), ?) RETURNING session_id, expires_at"

	row := sqlite.DB.QueryRow(query, expiresAt.Format(expiresAtFormat))
	err := row.Scan(&session.SessionID, &expiresAtString)

	if err != nil {
		err = fmt.Errorf("error retrieving sessions: %w", err)
		slog.Error(err.Error())
		return session, err
	}

	session.ExpiresAt, err = time.Parse(expiresAtFormat, expiresAtString)

	if err != nil {
		err = fmt.Errorf("error parsing session expiry: %w", err)
		slog.Error(err.Error())
		return session, err
	}

	return session, nil
}

func SetSessionCookie(w http.ResponseWriter, session SessionRecord) {
	cookies := http.Cookie{Name: "session_token", Value: session.SessionID, Expires: session.ExpiresAt, Path: "/"}
	http.SetCookie(w, &cookies)
}

func IsValidSession(sessionID string) bool {
	query := "SELECT expires_at FROM sessions WHERE session_id = ?"
	var expiresAtString string

	row := sqlite.DB.QueryRow(query, sessionID)
	err := row.Scan(&expiresAtString)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	if err != nil {
		slog.Error("error retrieving sessions", "error", err)
		return false
	}

	expiresAt, err := time.Parse(expiresAtFormat, expiresAtString)

	if err != nil {
		slog.Error("error parsing session expiry", "error", err)
		return false
	}

	if expiresAt.Before(time.Now()) {
		return false
	}

	return true
}
