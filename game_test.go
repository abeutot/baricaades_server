package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame(t *testing.T) {
	g := InitGame("test", "bob")

	assert.Equal(t, 1, len(g.players))
	assert.Equal(t, STATE_WAITING_FOR_PLAYERS, g.state)

	err := g.Start("bob")
	assert.Equal(t, "not enough players to start", err.Error())

	err = g.Join("alice")
	assert.Equal(t, nil, err)

	assert.Equal(t, 2, len(g.players))

	/* check cannot join twice */
	err = g.Join("alice")
	assert.Equal(t, "you are already part of the game", err.Error())

	err = g.RollDice("bob")
	assert.Equal(t, "game has not started yet", err.Error())

	err = g.Move("bob", "A3", "A5", "")
	assert.Equal(t, "game has not started yet", err.Error())

	err = g.Start("alice")
	assert.Equal(t, "only game creator can start the game", err.Error())

	err = g.Start("bob")
	assert.Equal(t, nil, err)

	err = g.Join("player")
	assert.Equal(t, "game has already started", err.Error())

	err = g.RollDice("alice")
	assert.Equal(t, "not your turn to roll", err.Error())

	err = g.RollDice("bob")
	assert.Equal(t, nil, err)

	assert.NotEqual(t, 0, g.dice)

	g.dice = 3

	err = g.Move("bob", "B1", "A3", "")
	assert.Equal(t, nil, err)

	assert.Equal(t, uint8(0), g.dice)

	err = g.Move("alice", "F2", "E3", "")
	assert.Equal(t, "roll the dice first", err.Error())

	g.dice = 3

	err = g.Move("alice", "A3", "B5", "")
	assert.Equal(t, "start position is not your own", err.Error())

	err = g.Move("alice", "F2", "E3", "")
	assert.Equal(t, nil, err)

	/* baricade in the way */
	g.dice = 3
	g.state = STATE_RED_PLAYING
	err = g.Move("bob", "A3", "B5", "")
	assert.Equal(t, "no available path between source and destination (check baricades)", err.Error())

	err = g.Move("bob", "A3", "D3", "")
	assert.Equal(t, nil, err)

	g.dice = 4
	g.state = STATE_RED_PLAYING
	err = g.Move("bob", "D3", "E5", "F3")
	assert.Equal(t, "not allowed baricade destination", err.Error())

	g.dice = 3
	err = g.Move("bob", "D3", "E5", "E4")
	assert.Equal(t, nil, err)

	/* cannot place baricade on to */
	g.state = STATE_RED_PLAYING
	g.dice = 1
	err = g.Move("bob", "E5", "E4", "E4")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "not allowed baricade destination", err.Error())

	/* two player mode, can select any color */
	/* check path m5 -> i4 with dice 5 is valid */
	g.state = STATE_BLUE_PLAYING
	g.dice = 5
	err = g.Move("alice", "N2", "M5", "H5")
	assert.Equal(t, nil, err)
	g.state = STATE_YELLOW_PLAYING
	g.dice = 5
	err = g.Move("bob", "L2", "I5", "I9")
	assert.Equal(t, nil, err)
	g.state = STATE_BLUE_PLAYING
	g.dice = 5
	err = g.Move("alice", "M5", "I4", "")
	assert.Equal(t, nil, err)

	/* TODO check if a move is possible, add option to pass if no path
	* available? */
}

func TestGameCannotJoinMoreThan4(t *testing.T) {
	g := InitGame("test", "bob")
	err := g.Join("alice")
	assert.Equal(t, nil, err)
	err = g.Join("player3")
	assert.Equal(t, nil, err)
	err = g.Join("player4")
	assert.Equal(t, nil, err)
	err = g.Join("player5")
	assert.Equal(t, "game is full", err.Error())
}
