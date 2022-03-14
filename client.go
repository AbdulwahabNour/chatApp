package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 120 * time.Second
	pongWait       = 120 * time.Second
	pingPeriod     = (pongWait * 8) / 10
	maxMessageSize = 25000
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewClient(conn *websocket.Conn, ws *WsServer, name string) *Client {
	return &Client{
		Name:     name,
		conn:     conn,
		wsServer: ws,
		send:     make(chan []byte),
		rooms:    make(map[*Room]bool)}
}

type Client struct {
	Name     string `json:"name"`
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	rooms    map[*Room]bool
}

func (c *Client) read() {
	defer c.conn.Close()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	c.conn.SetPongHandler(func(s string) error {

		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, jsonMsg, err := c.conn.ReadMessage()

		if err != nil {
			log.Println("Read :", err)
			return
		}
      
		c.handleNewMessage(jsonMsg)
	}

}

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod) 
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	c.conn.PingHandler()

	for {
		select {
		case msg, ok := <-c.send:

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println("write :", err)
				return
			}
			w.Write(msg)
			n := len(c.send)
 
			for i := 0; i < n; i++ {
				w.Write([]byte("/n"))
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		}

	}

}
func (c *Client) getName() string {
	return c.Name
}

func (c *Client) disconnect() {
	c.wsServer.leave <- c
	for r := range c.rooms {
		r.leave <- c
	}
}

func (c *Client) handleNewMessage(jsonMsg []byte) {

	var msg Message
 
	if err := json.Unmarshal(jsonMsg, &msg); err != nil {
		log.Printf("Error on unmarshal Message %s", err)

	}
	msg.Sender = c

	switch msg.Action {

	case SendMessageAction:
		roomName := msg.Target
		if room := c.wsServer.findRoom(roomName); room != nil {
			room.forward <- &msg
		}
	case JoinRoomActon:
		c.JoinRoomMessage(msg)
	case LeaveRoomAction:
		c.LeaveRoomMessage(msg)

	}
}
func (c *Client) JoinRoomMessage(msg Message) {
	roomName := msg.Message
	room := c.wsServer.findRoom(roomName)
	if room == nil {
		room = c.wsServer.createRoom(roomName)
	}
	c.rooms[room] = true
	room.join <- c
}
func (c *Client) LeaveRoomMessage(msg Message) {

	room := c.wsServer.findRoom(msg.Message)

	if room == nil {
		return
	}
	if _, ok := c.rooms[room]; ok {
		delete(c.rooms, room)
	}
	room.leave <- c
}

func Wshandler(w http.ResponseWriter, r *http.Request, ws *WsServer) {
	name, ok := r.URL.Query()["name"]
	if !ok || len(name[0]) < 1 {
		log.Println("Url parameter 'name' is missing")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(conn, ws, name[0])

	go client.read()
	go client.write()

	ws.join <- client

}
