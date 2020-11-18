package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// POST /user/register
// POST /game
// POST /game/<id>/join
// POST /game/<id>/start
// POST /game/<id>/roll_dice
// POST /game/<id>/move
// GET /game/<id>/subscribe
// GET /game/<id>

/* game states:
0 WAITING FOR PLAYERS
1 PLAYING
2 FINISHED
*/

func getGame(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(`{"coucou": "le monde"}`))
}

func main() {
	r := gin.Default()

	r.GET("/game/42", getGame)

	r.Run()
}
