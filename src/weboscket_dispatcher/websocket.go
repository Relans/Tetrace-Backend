package weboscket_dispatcher

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
)

var upgrader = websocket.Upgrader{
	Subprotocols: []string{"Sec-WebSocket-Extensions"},
}

type player struct {
	name         string
	subscription string
	lobbyId      int
	conn         *websocket.Conn
}

type message struct {
	messageType MessageType
	message     string
}

type lobby struct {
	id       int
	capacity int
	players  []*player
}

type MessageType int

const (
	SUBSCRIBE   MessageType = 0
	UNSUBSCRIBE MessageType = 1
	MESSAGE     MessageType = 2
)

var connections = make(map[*websocket.Conn]*player)
var subscriptions = make(map[string][]lobby)

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, w.Header())
	vars := r.URL.Query()["name"][0]
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		fmt.Println(err)
	}
	connections[conn] = &player{name: vars, subscription: "", conn: conn}
	go echo(conn)
	println(len(connections))
}

func echo(conn *websocket.Conn) {
	defer func() {
		conn.Close()
		delete(connections, conn)
		println("Connection closed")
	}()

	for {
		var message message;
		err := conn.ReadJSON(&message)
		if err != nil {
			println("Connection lost")
			break
		}

		p := connections[conn]
		switch message.messageType {
		case SUBSCRIBE:
			if p.subscription != message.message {
				delete(subscriptions, message.message)
				handleNewSubscription(message.message, conn)
			}
		case UNSUBSCRIBE:
			delete(subscriptions, message.message)
		case MESSAGE:
			for _, v := range subscriptions[p.subscription][p.lobbyId].players {
				if p != v {
					if w, err := p.conn.NextWriter(websocket.TextMessage); err != nil {
						return
					} else {
						w.Write([]byte(message.message))
					}
				}
			}
		}

	}
}

func handleNewSubscription(subscription string, conn *websocket.Conn) {
	p := connections[conn]
	avalibleLobbies := subscriptions[p.subscription]
	if avalibleLobbies != nil && len(avalibleLobbies) > 0 {
		avalibleLobbies[0].players = append(avalibleLobbies[0].players, p);
	} else {
		avalibleLobbies = []lobby{lobby{id: 0, players: []*player{p}}}
	}
	p.subscription = subscription
	p.lobbyId = 0
}
