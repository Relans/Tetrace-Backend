package connector

import (
	"gopkg.in/mgo.v2"
	"config"
)

var session *mgo.Session = nil

func CreateSession() *mgo.Session {
	Host := []string{
		config.HOST,
	}
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs: Host,
	})

	if err != nil {
		panic(err);
	}

	return session
}

func GetSession() *mgo.Session {
	if session == nil {
		session = CreateSession()
	}
	return session
}

func GetCollectionFromSession(collection string, session *mgo.Session) *mgo.Collection {
	return session.DB(config.DATABASE).C(collection)
}
