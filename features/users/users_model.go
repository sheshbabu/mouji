package users

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"mouji/commons/sqlite"
)

type UserRecord struct {
	Email    string
	Password string
}

func HasUsers() bool {
	query := "SELECT user_id FROM users LIMIT 1"
	var userID string

	row := sqlite.DB.QueryRow(query)
	err := row.Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return false
	}
	if err != nil {
		err = fmt.Errorf("error retrieving users: %w", err)
		panic(err)
	}

	return true
}

func GetUserByEmail(email string) (UserRecord, error) {
	var user UserRecord

	query := "SELECT email, password_hash FROM users where email = ?"

	row := sqlite.DB.QueryRow(query, email)
	err := row.Scan(&user.Email, &user.Password)
	if err != nil {
		err = fmt.Errorf("error retrieving user: %w", err)
		slog.Error(err.Error())
		return user, err
	}

	return user, nil
}

func insertUser(email string, passwordHash string, isAdmin bool) error {
	query := "INSERT INTO users (email, password_hash, is_admin) VALUES (?, ?, ?);"

	_, err := sqlite.DB.Exec(query, email, passwordHash, isAdmin)
	if err != nil {
		err = fmt.Errorf("error inserting user: %w", err)
		slog.Error(err.Error())
		return err
	}

	return nil
}
