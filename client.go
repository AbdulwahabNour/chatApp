package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)



var upgrader = websocket.Upgrader{
     ReadBufferSize : 1024,
     WriteBufferSize : 1024,
}

type Client struct{
    conn *websocket.Conn    
}

func newClient(conn *websocket.Conn) *Client{
    return &Client{conn: conn}
}

func Wshandler(w http.ResponseWriter, r *http.Request){
      conn, err := upgrader.Upgrade(w, r, nil)
      if err != nil{
           log.Println(err)
           return
      }

      client := newClient(conn)
  
}