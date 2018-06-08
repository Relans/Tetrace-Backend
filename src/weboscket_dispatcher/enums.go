package weboscket_dispatcher

type MessageType int

const (
	SUBSCRIBE   MessageType = 0
	UNSUBSCRIBE MessageType = 1
	MESSAGE     MessageType = 2
)
