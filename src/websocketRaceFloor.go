package main

import (
	"strconv"

	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceFloor(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username
	floorNum := d.FloorNum
	stageType := d.StageType

	/*
		Validation
	*/

	// Validate that the race exists
	var race *Race
	if v, ok := races[d.ID]; !ok {
		return
	} else {
		race = v
	}

	// Validate that the race has started
	if race.Status != "in progress" {
		return
	}

	// Validate that they are in the race
	var racer *Racer
	if v, ok := race.Racers[username]; !ok {
		return
	} else {
		racer = v
	}

	// Validate that they are still racing
	if racer.Status != "racing" {
		return
	}

	// Validate that the floor is sane
	if floorNum < 1 || floorNum > 13 {
		// The Void is floor 12, and we use floor 13 to signify Mega Satan
		log.Warning("User \"" + username + "\" attempted to update their floor, but \"" + strconv.Itoa(floorNum) + "\" is a bogus floor number.")
		websocketError(s, d.Command, "That is not a valid floor number.")
		return
	} else if stageType < 0 || stageType > 3 {
		// 3 is Greed Mode
		log.Warning("User \"" + username + "\" attempted to update their floor, but \"" + strconv.Itoa(stageType) + "\" is a bogus stage type.")
		websocketError(s, d.Command, "That is not a valid stage type.")
		return
	}

	/*
		Set the floor
	*/

	racer.FloorNum = floorNum
	racer.StageType = stageType
	racer.DatetimeArrivedFloor = getTimestamp()

	for racerName := range race.Racers {
		// Not all racers may be online during a race
		if s, ok := websocketSessions[racerName]; ok {
			type RacerSetFloorMessage struct {
				ID                   int    `json:"id"`
				Name                 string `json:"name"`
				FloorNum             int    `json:"floorNum"`
				StageType            int    `json:"stageType"`
				DatetimeArrivedFloor int64  `json:"datetimeArrivedFloor"`
			}
			websocketEmit(s, "racerSetFloor", &RacerSetFloorMessage{
				race.ID,
				username,
				floorNum,
				stageType,
				racer.DatetimeArrivedFloor,
			})
		}
	}

	race.SetAllPlaceMid()
}
