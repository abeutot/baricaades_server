package main

import (
	"errors"
	"math/rand"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/namsral/flag"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		/* TODO */
		return true
	},
}

func CustomContextMiddleware() gin.HandlerFunc {
	/* FIXME concurrency accesses */
	users := make(map[string]string)
	games := make(map[string]*Game)

	return func(c *gin.Context) {
		c.Set("users", users)
		c.Set("games", games)

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, found := c.MustGet("users").(map[string]string)
		if !found {
			panic("users not in context")
		}

		/* we protect every path except for register and game updates */
		matched, err := regexp.MatchString(`^(/user/register)|(/game/[^/]*/updates)$`, c.Request.URL.Path)
		if err != nil {
			panic(err)
		}
		if !matched {
			username, password, found := c.Request.BasicAuth()

			if !found {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":  http.StatusUnauthorized,
					"error": "invalid auth header",
				})
				return
			}

			token, found := users[username]
			if !found {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":  http.StatusUnauthorized,
					"error": "invalid credentials",
				})
				return
			}
			if token != password {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":  http.StatusUnauthorized,
					"error": "invalid credentials",
				})
				return
			}

			c.Set("auth_user", username)
		}

		c.Next()
	}
}

// POST /user/register
type UserBody struct {
	Username string `form:"username"`
}

func userRegister(c *gin.Context) {
	users, found := c.MustGet("users").(map[string]string)
	if !found {
		panic("users not in context")
	}

	var b UserBody
	c.Bind(&b)

	if _, alreadyExists := users[b.Username]; alreadyExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "username already exists",
		})
		return
	}

	/* generate a new token */
	sourceRunes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	token := make([]byte, 64)
	for i := range token {
		token[i] = sourceRunes[rand.Intn(len(sourceRunes))]
	}

	users[b.Username] = string(token)

	c.JSON(http.StatusCreated, gin.H{
		"username": b.Username,
		"token":    string(token),
	})
}

// GET /user/check
func userCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "success"})
}

func utilGetGamesNUser(c *gin.Context) (map[string]*Game, string) {
	games, found := c.MustGet("games").(map[string]*Game)
	if !found {
		panic("games not in context")
	}
	user, found := c.MustGet("auth_user").(string)
	if !found {
		panic("auth_user not in context")
	}

	return games, user
}

func utilGetGameNUser(c *gin.Context) (*Game, string, error) {
	games, user := utilGetGamesNUser(c)

	var gameParams GameParams

	if err := c.ShouldBindUri(&gameParams); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "error": "Invalid ID"})
		return nil, "", errors.New("invalid ID")
	}

	game, found := games[gameParams.ID]
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "error": "Invalid ID"})
		return nil, "", errors.New("invalid ID")
	}

	return game, user, nil
}

// GET /game
func listGames(c *gin.Context) {
	games, user := utilGetGamesNUser(c)

	sortedGames := make([]*Game, 0, len(games))
	mine := make([]*GameMiniJson, 0, 0)
	open := make([]*GameMiniJson, 0, 0)

	for id := range games {
		sortedGames = append(sortedGames, games[id])
	}

	/* sort by creation date DESC */
	sort.Slice(sortedGames, func(i, j int) bool {
		return sortedGames[i].created.After(sortedGames[j].created)
	})

	for _, g := range sortedGames {
		userInGame := false
		for _, p := range g.players {
			if p == user {
				userInGame = true
				break
			}
		}
		if userInGame {
			mine = append(mine, g.MiniJSON())
			continue
		}

		if g.state != STATE_WAITING_FOR_PLAYERS || len(g.players) == 4 {
			continue
		}

		open = append(open, g.MiniJSON())
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"mine": mine,
			"open": open,
		},
	)
}

// POST /game
func createGame(c *gin.Context) {
	/* TODO check that a user can create only one game */

	games, user := utilGetGamesNUser(c)

	/* generate new game id */
	gameID := make([]byte, 12)
	for { /* FIXME potential DOS */
		base32 := "BAFYBEICZSSCDSBS7FFQZ55ASQDF3SMV6KLCW3GOFSZVWLYARCI47BGF354"

		for i := range gameID {
			gameID[i] = base32[rand.Intn(len(base32))]
		}

		_, found := games[string(gameID)]
		if !found {
			break
		}
	}

	ID := string(gameID)
	games[ID] = InitGame(ID, user)

	c.JSON(201, gin.H{
		"id": string(gameID),
	})
}

type GameParams struct {
	ID string `uri:"id" binding:"required"`
}

// POST /game/<id>/start
func gameStart(c *gin.Context) {
	game, user, err := utilGetGameNUser(c)
	if err != nil {
		return
	}

	if err := game.Start(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{}) /* TODO body? */
}

// POST /game/<id>/join
func gameJoin(c *gin.Context) {
	game, user, err := utilGetGameNUser(c)
	if err != nil {
		return
	}

	err = game.Join(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{}) /* TODO body? */
}

// POST /game/<id>/roll_dice
func gameRollDice(c *gin.Context) {
	game, user, err := utilGetGameNUser(c)
	if err != nil {
		return
	}

	err = game.RollDice(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"dice": game.dice})
}

// POST /game/<id>/move
type MoveBody struct {
	From     string `form:"from"`
	To       string `to:"to"`
	Baricade string `baricade:"baricade"`
}

func gameMove(c *gin.Context) {
	game, user, err := utilGetGameNUser(c)
	if err != nil {
		return
	}

	var move MoveBody

	c.Bind(&move)

	if err := game.Move(user, strings.ToUpper(move.From), strings.ToUpper(move.To), strings.ToUpper(move.Baricade)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{}) /* TODO body? */
}

// GET /game/<id>
func getGame(c *gin.Context) {
	game, user, err := utilGetGameNUser(c)
	if err != nil {
		return
	}

	/* limit only to subscribed players */
	userInGame := false
	for _, u := range game.players {
		if u == user {
			userInGame = true
			break
		}
	}

	if !userInGame {
		c.JSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "error": "not in the game"})
		return
	}

	c.JSON(http.StatusOK, game.JSON())
}

// GET /game/<id>/updates
func gameUpdates(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "cannot upgrade connection",
		})
		return
	}
	defer ws.Close()

	var authData struct {
		Username string `json:"username"`
		Token    string `json:"token"`
		GameID   string `json:"game_id"`
	}
	err = ws.ReadJSON(&authData)
	if err != nil {
		panic("error: " + err.Error())
	}

	/* check auth */
	users, found := c.MustGet("users").(map[string]string)
	if !found {
		panic("users not in context")
	}

	token, found := users[authData.Username]
	if !found || token != authData.Token {
		err = ws.WriteJSON(struct {
			Code  int    `json:"code"`
			Error string `json:"error"`
		}{
			Code:  http.StatusUnauthorized,
			Error: "invalid auth credentials",
		})
		if err != nil {
			panic("error: " + err.Error())
		}
		return
	}

	games, found := c.MustGet("games").(map[string]*Game)
	if !found {
		panic("games not in context")
	}
	var gameParams GameParams
	if err := c.ShouldBindUri(&gameParams); err != nil {
		panic("ID not in URL")
	}
	game, found := games[gameParams.ID]
	if !found {
		err = ws.WriteJSON(struct {
			Code  int    `json:"code"`
			Error string `json:"error"`
		}{
			Code:  http.StatusNotFound,
			Error: "game not found",
		})
		if err != nil {
			panic("error: " + err.Error())
		}
		return
	}

	isPartOfGame := false
	for i := range game.players {
		if game.players[i] == authData.Username {
			isPartOfGame = true
			break
		}
	}

	if !isPartOfGame {
		err = ws.WriteJSON(struct {
			Code  int    `json:"code"`
			Error string `json:"error"`
		}{
			Code:  http.StatusForbidden,
			Error: "you are not part of this game",
		})
		if err != nil {
			panic("error: " + err.Error())
		}
		return
	}

	/* setup a channel on the game */
	updates := game.Subscribe()
	defer game.Unsubscribe(updates)

	/* get updates from game channel */
	for u := range updates {
		err = ws.WriteJSON(u)
		if err != nil {
			panic("error: " + err.Error())
		}
	}
}

func main() {
	var corsAllowOrigins string

	flag.StringVar(&corsAllowOrigins, "cors_allow_origins", "http://localhost:3000", "Origins to allow with CORS (comma separated)")
	flag.Parse()

	r := gin.Default()

	/* CORS config */
	r.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(corsAllowOrigins, ","),
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "authorization"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))
	r.Use(CustomContextMiddleware())
	r.Use(AuthMiddleware())

	r.POST("/user/register", userRegister)
	r.GET("/user/check", userCheck)

	r.GET("/game", listGames)
	r.POST("/game", createGame)
	r.POST("/game/:id/start", gameStart)
	r.POST("/game/:id/join", gameJoin)
	r.POST("/game/:id/roll_dice", gameRollDice)
	r.POST("/game/:id/move", gameMove)
	r.GET("/game/:id", getGame)
	r.GET("/game/:id/updates", gameUpdates)

	r.Run()
}
