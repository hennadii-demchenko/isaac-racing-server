package models

import (
	"database/sql"
)

/*
	These are more functions for querying the "users" table,
	but these functions are only used for the website
*/

/*
	Data structures
*/

// StreamURL gets each row for all profiles
type StreamURL struct {
	RacerURL sql.NullString
}

// GetStreamURL gets the stream url
func (*Users) GetStreamURL(username string) (string, error) {
	// Find total amount of users
	var streamURL string
	if err := db.QueryRow(`
		SELECT stream_url
		FROM users
		WHERE username = ?
	`, username).Scan(&streamURL); err != nil {
		return "", err
	}
	return streamURL, nil
}
