package connector

import (
	"gopkg.in/mgo.v2/bson"
	"config"
)

type PlayerStatus int

const (
	OFFLINE   PlayerStatus = 0
	MAIN_MENU PlayerStatus = 1
	IN_LOBBY  PlayerStatus = 2
	PLAYING   PlayerStatus = 3
)

type Player struct {
	Id           string       `bson:"_id"`
	Name         string       `bson:"name"`
	PlayerStatus PlayerStatus `bson:"status"`
}

func (playerStatus PlayerStatus) String() string {
	names := [...]string{
		"OFFLINE",
		"MAIN_MENU",
		"IN_LOBBY",
		"PLAYING"}

	if playerStatus < OFFLINE || playerStatus > PLAYING {
		return "Unknown"
	}

	return names[playerStatus]
}

var playerCollection = GetCollectionFromSession(config.PLAYER_COLLECTION, GetSession())

func GetPlayers() []Player {
	var players []Player
	var err = playerCollection.Find(nil).All(&players)
	if err != nil {
		panic(err)
	}
	return players
}

func AddPlayer(name string) Player {
	var player = Player{Name: name}
	playerCollection.Insert(player);
	playerCollection.Find(bson.M{"name": name, "status": OFFLINE}).One(&player)
	return player
}

func GetPlayerByName(name string) Player {
	var player Player
	if err := playerCollection.Find(bson.M{"name": name}).One(&player); err != nil {
		player = AddPlayer(name);
	}
	return player;
}

func LogoutPlayer(name string) {
	playerCollection.Update(bson.M{"name": name}, bson.M{"$set": bson.M{"status": OFFLINE}});
}
