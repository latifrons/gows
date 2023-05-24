package gows

import "github.com/sirupsen/logrus"

type UnicastMessage struct {
	Client  *Client
	Message []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
	unicast    chan *UnicastMessage
	Status     int
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		unicast:    make(chan *UnicastMessage),
		Status:     0,
	}
}

func (h *Hub) Broadcast(bytes []byte) {
	h.broadcast <- bytes
}

func (h *Hub) Unicast(msg *UnicastMessage) {
	h.unicast <- msg
}

func (h *Hub) Run(closeOnNobody bool) {
RUN:
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.Status = 1
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.sendChan)
			}
			if len(h.clients) == 0 && h.Status == 1 && closeOnNobody {
				logrus.Info("hub stopped")
				break RUN
			}
		case message := <-h.broadcast:
			logrus.WithField("count", len(h.clients)).Info("Broadcasting: ", string(message))
			if len(h.clients) == 0 && h.Status == 1 && closeOnNobody {
				logrus.Info("hub stopped")
				break RUN
			}
			for client := range h.clients {
				client.sendChan <- message
			}
		case umessage := <-h.unicast:
			if _, ok := h.clients[umessage.Client]; ok {
				umessage.Client.sendChan <- umessage.Message
			}
			logrus.Info("Unicasted: ", string(umessage.Message))
		}
	}
	logrus.Info("RUN ends")
	h.Status = 2
}
