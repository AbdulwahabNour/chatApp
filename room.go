package main

import (
	"fmt"

	"github.com/google/uuid"
)

func NewRoom(name string, private bool) *Room {
	return &Room{
		ID:      uuid.New(),
		Name:    name,
		Private: private,
		clients: make(map[*Client]bool),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		forward: make(chan *Message),
	}
}

type Room struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	clients map[*Client]bool
	join    chan *Client
	leave   chan *Client
	forward chan *Message
	Private bool `json:"private"`
}

func (r *Room) RunRoom() {

	for {
		select {
		case c := <-r.join:

			r.addClientInRoom(c)
		case c := <-r.leave:
			r.removeClientFromRoom(c)
		case msg := <-r.forward:
			{
				r.forwardToClientsInRoom(msg.encode())
			}
		}
	}
}

func (r *Room) addClientInRoom(c *Client) {
	r.clients[c] = true
	r.notifyClientJoined(c)
	

}
func (r *Room) removeClientFromRoom(c *Client) {
	if _, ok := r.clients[c]; ok {
		delete(r.clients, c)
	}
}
func (r *Room) forwardToClientsInRoom(msg []byte) {
 
	for c := range r.clients {
	 
		c.send <- msg
	}
}

func (r *Room) notifyClientJoined(c *Client) {
	msg := &Message{
		Action:  JoinRoomActon,
		Target:  r,
		Message: fmt.Sprintf("%s joind the room", c.Name),
		Sender:  []*Client{c},
	}
  
	r.forwardToClientsInRoom(msg.encode())
}
func (r *Room) GetId() string {
	return r.ID.String()
}
func (r *Room) GetName() string {
	return r.Name
}
func (r *Room) registerClientInRoom(c *Client) {
	if !r.Private {
		r.notifyClientJoined(c)
	}
	r.clients[c] = true
}
