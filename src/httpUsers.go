package main

import (
	"encoding/json"
	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/gin-gonic/gin"
	//"net/http"
)

// TODO
// Add stream URL for each opponents

func httpUserAPI(c *gin.Context) {
	// Local variables
	w := c.Writer
	// Start serving the header so that we can decide which data to send later
	w.Header().Set("Content-Type", "application/json")

	// Opponent struct to give information on each
	type Opponent struct {
		Username    string `json:"user"`
		TwitchName  string `json:"twitchName"`
		ReadyStatus string `json:"readyStatus"`
		FloorNum    int    `json:"floorNum"`
		StageType   int    `json:"stageType"`
	}

	// CurrentOpponents holds all the race info and the map of opponents
	type CurrentOpponents struct {
		RaceID     int    `json:"raceid"`
		RaceName   string `json:"raceName"`
		RaceStatus string `json:"raceStatus"`
		Opponents  []Opponent
	}

	// Needs to be init'd after we define the local structs
	var opponents []Opponent
	var opponent Opponent
	var currentOpponents CurrentOpponents

	// Parse the player name from the URL
	player := c.Params.ByName("racername")
	if player == "" {
		// http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		w.Write([]byte("{\"Error\": \"Please use a racer account\"}"))
		return
	}

	// Run through all the races
	for _, raceID := range races {
		// If we see that the name in the URL is a player in the map, collect info
		if _, found := raceID.Racers[player]; found {
			currentOpponents.RaceID = raceID.ID
			currentOpponents.RaceName = raceID.Name
			currentOpponents.RaceStatus = raceID.Status
			for _, racer := range raceID.Racers {
				// Remove self from the list of opponents
				if racer.Name != player {
					opponent.Username = racer.Name
					opponent.TwitchName = racer.Name
					opponent.ReadyStatus = racer.Status
					opponent.FloorNum = racer.FloorNum
					opponent.StageType = racer.StageType
					opponents = append(opponents, opponent)
				}
			}
			currentOpponents.Opponents = opponents
			// We break here because we want to stop looking
			break
		}
	}
	// Once we've gotten the race and opponent info convert it to JSON
	jsonData, err := json.MarshalIndent(currentOpponents, "", "\t")
	if err != nil {
		log.Error("Couldn't generate JSON")
		w.Write([]byte("Please search for a user"))
		return
	}
	w.Write(jsonData)

}
