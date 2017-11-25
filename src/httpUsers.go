package main

import (
	"encoding/json"
	//"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func httpUserAPI(c *gin.Context) {
	// Local variables
	w := c.Writer
	type Racer struct {
		Username   string `json:"user"`
		TwitchName string `json:"twitchName"`
	}

	// Parse the player name from the URL
	player := c.Params.ByName("racername")
	if player == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	user := []byte(`{"user":"player","twitchName":"player","currentOpponents":[{"name":"Zamiel","twitchName":"Zamiell"},{"name":"Hyphen_ated","twitchname":"Hyphen_ated"}]}`)

	w.Header().Set("Content-Type", "application/json")
	w.Write(user)
}
