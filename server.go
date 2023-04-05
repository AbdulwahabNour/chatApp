package main

func NewWsServer() *WsServer {
	return &WsServer{
		forward: make(chan []byte),
		clients: make(map[*Client]bool),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		rooms:   map[*Room]bool{},
	}
}

type WsServer struct {
	forward chan []byte
	join    chan *Client
	leave   chan *Client
	clients map[*Client]bool
	rooms   map[*Room]bool
}

func (s *WsServer) Run() {
	for {
		select {
		case c := <-s.join:
			s.clients[c] = true
			s.listOnlineClients(c)

		case c := <-s.leave:
			if _, ok := s.clients[c]; ok {
				delete(s.clients, c)
			}
			s.notifyClientLeft(c)
		case msg := <-s.forward:

			s.forwardMsg(msg)
		}
	}
}
func (s *WsServer) forwardMsg(msg []byte) {

	for c := range s.clients {

		c.send <- msg
	}
}
func (s *WsServer) findRoomByName(name string) *Room {
	var room *Room

	for r := range s.rooms {
		if r.Name == name {
			room = r
			break
		}
	}
	return room
}
func (s *WsServer) createRoom(name string, private bool) *Room {
	room := NewRoom(name, private)

	go room.RunRoom()

	s.rooms[room] = true

	return room
}
func (s *WsServer) notifyClientJoined(c *Client) {
	msg := &Message{
		Action: UserJoinedAction,
		Sender: []*Client{c},
	}
	s.forwardMsg(msg.encode())
}

func (s *WsServer) notifyClientLeft(c *Client) {
	msg := &Message{
		Action: UserLeftAction,
		Sender: []*Client{c},
	}
	s.forwardMsg(msg.encode())
}
func (s *WsServer) listOnlineClients(c *Client) {
	msg := &Message{
		Action: UserJoinedAction,
	}

	for client := range s.clients {
		msg.Sender = append(msg.Sender, client)

	}
	c.send <- msg.encode()
	s.forwardMsgClients(c)
}
func (s *WsServer) registerClient(c *Client) {
	s.notifyClientJoined(c)
	s.listOnlineClients(c)
	s.clients[c] = true
}
func (s *WsServer) unRegisterClient(c *Client) {
	if _, ok := s.clients[c]; ok {
		delete(s.clients, c)
		s.notifyClientLeft(c)
	}
}
func (s *WsServer) forwardMsgClients(c *Client) {
	msg := &Message{Action: UserJoinedAction, Sender: []*Client{c}}
	for client := range s.clients {
		if client.ID == c.ID {
			continue
		}
		client.send <- msg.encode()

	}
}

func (s *WsServer) findRoomByID(id string) *Room {
	var room *Room
	for r := range s.rooms {
		if r.GetId() == id {
			room = r
			break
		}
	}
	return room
}
func (s *WsServer) findClientByID(id string) *Client {
	var client *Client
    for c := range s.clients{
		if c.getUserId() == id{
			client = c
			break
		}
	}
	return client
}
