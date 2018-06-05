package connector

import (
	"config"
	"gopkg.in/mgo.v2/bson"
)

type GameStatus int

const (
	WAITING     GameStatus = 0
	IN_PROGRESS GameStatus = 1
)

type Game struct {
	Id         string     `bson:"_id"`
	MaxPlayers int        `bson:"max_players"`
	Players    []Player   `bson:"players"`
	GameStatus GameStatus `bson:"status"`
	Full       bool       `bson:"full"`
}

func (gameStatus GameStatus) String() string {
	names := [...]string{
		"WAITING",
		"IN_PROGRESS",
		"IN_LOBBY",
		"PLAYING"}

	if gameStatus < WAITING || gameStatus > IN_PROGRESS {
		return "Unknown"
	}

	return names[gameStatus]
}

var gameCollection = GetCollectionFromSession(config.GAME_COLLECTION, GetSession())

func GetOpenGame(maxPlayers int) Game {
	var games []Game
	if err := gameCollection.Find(bson.M{"full": false, "status": WAITING}).All(&games); err != nil || len(games) == 0 {
		return NewGame(maxPlayers)
	} else {
		return games[0]
	}
}

func NewGame(maxPlayers int) Game {
	var game = Game{Id: bson.NewObjectId().String(), Full: false, MaxPlayers: maxPlayers, GameStatus: WAITING, Players: nil}
	gameCollection.Insert(game)
	return game;
}
