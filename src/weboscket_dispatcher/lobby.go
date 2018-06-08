package weboscket_dispatcher

type lobby struct {
	id       string
	capacity int
	clients  []*Client
}

func (l *lobby) addClient(client *Client) {
	l.clients = append(l.clients, client)
}

func (l *lobby) removeClient(client *Client) {
	for i, element := range l.clients {
		if element == client {
			copy(l.clients[i:], l.clients[i+1:])
			l.clients[len(l.clients)-1] = nil
			l.clients = l.clients[:len(l.clients)-1]
		}
	}
}

func (l *lobby) canClientJoin() bool {
	return l.capacity > len(l.clients)
}
