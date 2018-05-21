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
	w.Header().Set("Access-Control-Allow-Origin", "http://192.168.1.247:81")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	// Opponent struct to give information on each
	type Opponent struct {
		Username     string  `json:"user"`
		TwitchName   string  `json:"twitchName"`
		ReadyStatus  string  `json:"readyStatus"`
		Place        int     `json:"place"`
		PlaceMid     int     `json:"placeMid"`
		FloorNum     int     `json:"floorNum"`
		StageType    int     `json:"stageType"`
		Seed         string  `json:"seed"`
		StartingItem int     `json:"startingItem"`
		RunTime      int64   `json:"runTime"`
		Items        []*Item `json:"items"`
	}

	// CurrentRace holds all the race info and the map of opponents
	type CurrentRace struct {
		RaceID          int        `json:"raceID"`
		RaceName        string     `json:"raceName"`
		RaceStatus      string     `json:"raceStatus"`
		Opponents       []Opponent `json:"raceOpponents"`
		RaceStartedTime int64      `json:"raceStartedTime"`
	}

	// Needs to be init'd after we define the local structs
	var opponents []Opponent
	var opponent Opponent
	var currentRace CurrentRace

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
			currentRace.RaceID = raceID.ID
			currentRace.RaceName = raceID.Name
			currentRace.RaceStatus = raceID.Status
			currentRace.RaceStartedTime = raceID.DatetimeStarted
			// List everyone so that we can set their information
			for _, racer := range raceID.Racers {
				streamURL, err := db.Users.GetStreamURL(racer.Name)
				if err != nil {
					log.Info("Failed to get streamURL for, '" + player + "', from the database: " + err.Error())
				}
				opponent.Username = racer.Name
				opponent.TwitchName = streamURL
				opponent.ReadyStatus = racer.Status
				opponent.FloorNum = racer.FloorNum
				opponent.StageType = racer.StageType
				opponent.Place = racer.Place
				opponent.PlaceMid = racer.PlaceMid
				opponent.Seed = racer.Seed
				opponent.StartingItem = racer.StartingItem
				opponent.RunTime = racer.RunTime
				opponent.Items = racer.Items
				opponents = append(opponents, opponent)
			}
			currentRace.Opponents = opponents
			// We break here because we want to stop looking
			break
		}
	}
	// Once we've gotten the race and opponent info convert it to JSON
	jsonData, err := json.MarshalIndent(currentRace, "", "\t")
	if err != nil {
		log.Error("Couldn't generate JSON")
		w.Write([]byte("Please search for a user"))
		return
	}
	w.WriteHeaderNow()
	w.Write(jsonData)

}
