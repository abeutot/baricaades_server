package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/namsral/flag"
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
	var corsAllowOrigins string

	flag.StringVar(&corsAllowOrigins, "cors_allow_origins", "http://localhost:3000", "Origins to allow with CORS (comma separated)")
	flag.Parse()

	r := gin.Default()

	/* CORS config */
	r.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(corsAllowOrigins, ","),
		AllowMethods: []string{"GET", "POST"},
		MaxAge:       24 * time.Hour,
	}))

	r.GET("/game/42", getGame)

	r.Run()
}
