package main

import (
	"github.com/gin-gonic/gin"
)

// TODO
// Add stream URL for each opponents

func httpDocAPI(c *gin.Context) {
	// Local variables
	w := c.Writer

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"Documentation\": \"API Documentation\" }"))

}
