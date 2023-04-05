package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/google/uuid"
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
		ID:       uuid.New(),
		Name:     name,
		conn:     conn,
		wsServer: ws,
		send:     make(chan []byte),
		rooms:    make(map[*Room]bool)}
}

type Client struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	rooms    map[*Room]bool
}

func (c *Client) read() {

	defer func() {
		c.conn.Close()
		c.disconnect()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	c.conn.SetPongHandler(func(s string) error {
		fmt.Println("handle pong => ", s)
		c.conn.SetReadDeadline(time.Now().Add(pongWait))

		return nil
	})

	for {
		_, jsonMsg, err := c.conn.ReadMessage()

		if err != nil {
			log.Println("Read :", err)
			break
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

				return
			}
			w.Write(msg)
			 
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("ticker")
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
	
	
	if err := msg.decode(jsonMsg); err != nil {
		log.Printf("Error on unmarshal Message %s", err)
		return
	}

	msg.Sender = []*Client{c}
 
	switch msg.Action {

	case SendMessageAction:
		
		roomID := msg.Target.GetId()
		room := c.wsServer.findRoomByID(roomID)
		if room != nil {
			room.forward <- &msg
		}

	case JoinRoomActon:
		c.JoinRoomMessage(msg)

	case LeaveRoomAction:
		c.LeaveRoomMessage(msg)

	case JoinRoomPrivateAction:

		c.handleJoinRoomPrivate(msg)

	}
}
func (c *Client) handleJoinRoomPrivate(m Message) {

         target := c.wsServer.findClientByID(m.Message)
		 if target == nil{
			 return
		 }

		 roomName := m.Message + c.getUserId()
	 
		 c.joinRoom(roomName, target)
		 c.joinRoom(roomName, c)

}
func (c *Client)joinRoom(roomName string, sender *Client){
	  room := c.wsServer.findRoomByName(roomName)
 
	if room == nil {
		
		room = c.wsServer.createRoom(roomName, sender != nil)
	}
	if sender == nil && room.Private {
              return
	}
	
	if !c.isInRoom(room){
		c.rooms[room] = true 
		room.join <- c 
	}


}
func(c *Client)isInRoom(r *Room)bool{
	 if _,ok := c.rooms[r]; ok{
		 return true
	 }

	 return false
}
func (c *Client) JoinRoomMessage(msg Message) {
	roomName := msg.Message
	c.joinRoom(roomName, nil)
 
}
func (c *Client) LeaveRoomMessage(msg Message) {

	roomId := msg.Target.GetId()
	room := c.wsServer.findRoomByID(roomId)

	if room == nil {
		return
	}

	if _, ok := c.rooms[room]; ok {
		delete(c.rooms, room)
	}

	room.leave <- c
}
func (c *Client)getUserId()string{
	return c.ID.String()
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

	fmt.Printf("%v \n %v \n", ws.rooms, runtime.NumGoroutine())

}
