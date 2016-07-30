package model

/*
 *  Data types
 */

type RaceParticipantItems struct {
	db *Model
}

/*
 *  race_participant_items table functions
 */

func (self *RaceParticipantItems) GetItemList(username string, raceID int) ([]Item, error) {
	// Get the items that this user got so far in the race
	rows, err := db.Query(
		"SELECT item_id, floor FROM race_participant_items "+
			"WHERE race_participant_id = (SELECT id FROM race_participants WHERE user_id = (SELECT id FROM users WHERE username = ?) AND race_id = ?)",
		username,
		raceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// We have to initialize this way to avoid sending a null on an empty array: https://danott.co/posts/json-marshalling-empty-slices-to-empty-arrays-in-go.html
	itemList := make([]Item, 0)
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Floor)
		if err != nil {
			return nil, err
		}
		itemList = append(itemList, item)
	}

	return itemList, nil
}

func (self *RaceParticipantItems) Insert(userID int, raceID int, itemID int, floor int) error {
	// Add the user to the participants list for that race
	stmt, err := db.Prepare(
		"INSERT INTO race_participant_items (race_participant_id, item_id, floor) " +
			"VALUES ((SELECT id FROM race_participants WHERE user_id = ? AND race_id = ?), ?, ?)",
	)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID, raceID, itemID, floor)
	if err != nil {
		return err
	}

	return nil
}
