package weboscket_dispatcher

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
	"strconv"
	"sync"
)

type Player struct {
	name string
}

type Client struct {
	conn         *websocket.Conn
	player       *Player
	subscription string
	lobby        string
	mu           sync.Mutex
}

type message struct {
	MessageType MessageType `json:"messageType"`
	Message     string      `json:"message"`
}

var upgrader = websocket.Upgrader{
	Subprotocols: []string{"Sec-WebSocket-Extensions"},
}

var rooms = map[string]map[string]*lobby{"main": {"0": &lobby{id: "0", capacity: -1}}}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		fmt.Println(err)
	}
	//vars := r.URL.Query()["name"]

	player := &Player{name: "test"}
	client := &Client{player: player, conn: conn, subscription: "main", lobby: "0"}
	rooms["main"]["0"].addClient(client)
	go echo(client)
}

func echo(client *Client) {
	defer func() {
		client.conn.Close()
		rooms[client.subscription][client.lobby].removeClient(client)
		client.player = nil
		client = nil
		println("Connection closed")
	}()

	for {
		message := message{}
		err := client.conn.ReadJSON(&message)
		if err != nil {
			println("Error while reading message", err.Error())
			break
		}

		switch message.MessageType {
		case SUBSCRIBE:
			unsubscribeClien(client)
			handleNewSubscription(message.Message, client)
		case UNSUBSCRIBE:
			unsubscribeClien(client)
		case MESSAGE:
			println("Message from client: ", message.Message)
			for _, c := range rooms[client.subscription][client.lobby].clients {
				if c != client && c != nil && c.conn != nil {
					if c.send(message) != nil {
						println("Client left!")
						c.conn.Close();
						rooms[client.subscription][client.lobby].removeClient(c)
					}
				}
			}
		}
	}
}

func (c *Client) send(v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(v)
}

func handleNewSubscription(subscription string, client *Client) {
	var newLobby *lobby
	availableLobbies, ok := rooms[subscription]
	if availableLobbies != nil && len(availableLobbies) > 0 && ok {
		newLobby = getFreeLobbyForSubscription(subscription)
	} else {
		if availableLobbies == nil || !ok {
			availableLobbies = make(map[string]*lobby)
			rooms[subscription] = availableLobbies
		}
		newLobby = &lobby{id: "ROOM-0", capacity: 16}
		availableLobbies["ROOM-0"] = newLobby
	}
	newLobby.addClient(client)
	client.lobby = newLobby.id
	client.subscription = subscription
}

func getFreeLobbyForSubscription(subscription string) *lobby {
	for _, v := range rooms[subscription] {
		if v.canClientJoin() {
			return v
		}
	}
	newLobby := lobby{id: "ROOM-0" + strconv.Itoa(len(rooms[subscription])), capacity: 16}
	rooms[subscription][newLobby.id] = &newLobby

	return &newLobby
}

func unsubscribeClien(client *Client) {
	lobbies, ok := rooms[client.subscription]
	if (ok) {
		l, ok := lobbies[client.lobby]
		if (ok) {
			l.removeClient(client)
			client.subscription = "main"
			client.lobby = "0"
		}
	}

}
