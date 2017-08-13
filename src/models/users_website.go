package models

import (
	"database/sql"
	"strconv"
)

/*
	These are more functions for querying the "users" table,
	but these functions are only used for the website
*/

type StatsSeeded struct {
	ELO            int
	LastELOChange  int
	NumSeededRaces int
	LastSeededRace int
}

type StatsUnseeded struct {
	UnseededAdjustedAverage int
	UnseededRealAverage     int
	NumUnseededRaces        int
	NumForfeits             int
	ForfeitPenalty          int
	LowestUnseededTime      int
	LastUnseededRace        int
}

// Used in the leaderboards
type LeaderboardRowSeeded struct {
	Name           string
	ELO            int
	LastELOChange  int
	NumSeededRaces int
	LastSeededRace int
	Verified       int
}
type LeaderboardRowUnseeded struct {
	Name                    string
	UnseededAdjustedAverage int
	UnseededRealAverage     int
	NumUnseededRaces        int
	NumForfeits             int
	ForfeitPenalty          int
	LowestUnseededTime      int
	LastUnseededRace        int
	Verified                int
}
type LeaderboardRowTopTimes struct {
	Name     string
	Time     int
	Date     int
	Verified int
}
type LeaderboardRowMostPlayed struct {
	Name     string
	Total    int
	Verified int
}

type ProfilesRow struct {
	Username     string
	DateCreated  int
	StreamUrl    string
	Achievements int
}
type ProfileData struct {
	Username      string
	DateCreated   int
	Verified      int
	ELO           int
	LastELOChange int
	SeededRaces   int
	UnseededRaces int
	StreamUrl     string
}

func (*Users) GetStatsSeeded(username string) (StatsSeeded, error) {
	var stats StatsSeeded
	if err := db.QueryRow(`
		SELECT
			elo,
			last_elo_change,
			num_seeded_races,
			last_seeded_race
		FROM
			users
		WHERE
			username = ?
	`, username).Scan(
		&stats.ELO,
		&stats.LastELOChange,
		&stats.NumSeededRaces,
		&stats.LastSeededRace,
	); err != nil {
		return stats, err
	} else {
		return stats, nil
	}
}

func (*Users) GetStatsUnseeded(username string) (StatsUnseeded, error) {
	var stats StatsUnseeded
	if err := db.QueryRow(`
		SELECT
			unseeded_adjusted_average,
			unseeded_real_average,
			num_unseeded_races,
			num_forfeits,
			forfeit_penalty,
			lowest_unseeded_time,
			last_unseeded_race
		FROM
			users
		WHERE
			username = ?
	`, username).Scan(
		&stats.UnseededAdjustedAverage,
		&stats.UnseededRealAverage,
		&stats.NumUnseededRaces,
		&stats.NumForfeits,
		&stats.ForfeitPenalty,
		&stats.LowestUnseededTime,
		&stats.LastUnseededRace,
	); err != nil {
		return stats, err
	} else {
		return stats, nil
	}
}

// Make a leaderboard for the seeded format based on all of the users
func (*Users) GetLeaderboardSeeded() ([]LeaderboardRowSeeded, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			username,
			elo,
			last_elo_change,
			num_seeded_races,
			last_seeded_race
		FROM
			users
		WHERE
			num_seeded_races > 1
	`); err != nil {
		return nil, err
	} else {
		rows = v
	}
	defer rows.Close()

	// Iterate over the users
	leaderboard := make([]LeaderboardRowSeeded, 0)
	for rows.Next() {
		var row LeaderboardRowSeeded
		if err := rows.Scan(
			&row.Name,
			&row.ELO,
			&row.LastELOChange,
			&row.NumSeededRaces,
			&row.LastSeededRace,
		); err != nil {
			return nil, err
		}

		// Append this row to the leaderboard
		leaderboard = append(leaderboard, row)
	}

	return leaderboard, nil
}

// Make a leaderboard for the unseeded format based on all of the users
func (*Users) GetLeaderboardUnseeded() ([]LeaderboardRowUnseeded, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			username,
			unseeded_adjusted_average,
			unseeded_real_average,
			num_unseeded_races,
			num_forfeits,
			forfeit_penalty,
			lowest_unseeded_time,
			last_unseeded_race
		FROM
			users
		WHERE
			num_unseeded_races > 15
	`); err != nil {
		return nil, err
	} else {
		rows = v
	}
	defer rows.Close()

	// Iterate over the users
	leaderboard := make([]LeaderboardRowUnseeded, 0)
	for rows.Next() {
		var row LeaderboardRowUnseeded
		if err := rows.Scan(
			&row.Name,
			&row.UnseededAdjustedAverage,
			&row.UnseededRealAverage,
			&row.NumUnseededRaces,
			&row.NumForfeits,
			&row.ForfeitPenalty,
			&row.LowestUnseededTime,
			&row.LastUnseededRace,
		); err != nil {
			return nil, err
		}

		// Append this row to the leaderboard
		leaderboard = append(leaderboard, row)
	}

	return leaderboard, nil
}

// Get player data to populate the player's profile page
func (*Users) GetProfileData(username string) (ProfileData, error) {
	var profileData ProfileData
	if err := db.QueryRow(`
		SELECT
			username,
			datetime_created,
			verified,
			elo,
			last_elo_change,
			num_seeded_races,
			num_unseeded_races,
			stream_url
		FROM
			users
		WHERE
			steam_id > 0 and
			username = ?
	`, username).Scan(
		&profileData.Username,
		&profileData.DateCreated,
		&profileData.Verified,
		&profileData.ELO,
		&profileData.LastELOChange,
		&profileData.SeededRaces,
		&profileData.UnseededRaces,
		&profileData.StreamUrl,
	); err == sql.ErrNoRows {
		return profileData, nil
	} else if err != nil {
		return profileData, err
	} else {
		return profileData, nil
	}
}

// Get players data to populate the profiles page
func (*Users) GetUserProfiles(currentPage int, usersPerPage int) ([]ProfilesRow, int, error) {
	usersOffset := (currentPage - 1) * usersPerPage
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			u.username,
			u.datetime_created,
			u.stream_url,
			count(ua.achievement_id)
		FROM
			users u
		LEFT JOIN
			user_achievements ua
			ON
				u.id = ua.user_id
		WHERE
			u.steam_id > 0
		GROUP BY
			u.username
		ORDER BY
			u.username ASC
		LIMIT
			?
		OFFSET
			?
	`, strconv.Itoa(usersPerPage), strconv.Itoa(usersOffset)); err == sql.ErrNoRows {
		return nil, 0, nil
	} else if err != nil {
		return nil, 0, err
	} else {
		rows = v
	}
	defer rows.Close()

	// Iterate over the user profile results
	profiles := make([]ProfilesRow, 0)
	for rows.Next() {
		var row ProfilesRow
		if err := rows.Scan(
			&row.Username,
			&row.DateCreated,
			&row.StreamUrl,
			&row.Achievements,
		); err != nil {
			return nil, 0, err
		}

		// Append this row to the leaderboard
		profiles = append(profiles, row)
	}

	// Find total amount of users
	var allProfilesCount int
	if err := db.QueryRow(`
		SELECT count(id)
		FROM users
		WHERE steam_id > 0
	`).Scan(&allProfilesCount); err != nil {
		return nil, 0, err
	}

	return profiles, allProfilesCount, nil
}