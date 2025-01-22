package config

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"mouji/commons/sqlite"
)

func SetConfig(key string, value string) error {
	query := "INSERT OR REPLACE INTO config (key, value) VALUES (?, ?)"

	_, err := sqlite.DB.Exec(query, key, value)
	if err != nil {
		err = fmt.Errorf("error upserting app setting: %w", err)
		slog.Error(err.Error())
		return err
	}

	return nil
}

func GetConfig(key string) (string, error) {
	value := ""

	query := "SELECT value FROM config where key = ?"

	row := sqlite.DB.QueryRow(query, key)
	err := row.Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return value, nil
		}

		err = fmt.Errorf("error retrieving app setting: %w", err)
		slog.Error(err.Error())
		return value, err
	}

	return value, nil
}
