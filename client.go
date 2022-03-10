package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const(
    writeWait = 10 * time.Second
    pongWait =  120 * time.Second
    pingPeriod = (pongWait * 8)/10
    maxMessageSize = 25000
)


var upgrader = websocket.Upgrader{
     ReadBufferSize : 1024,
     WriteBufferSize : 1024,
}


func NewClient(conn *websocket.Conn,ws  *WsServer) *Client{
    return &Client{ conn: conn,
                    wsServer: ws,
                    send: make(chan []byte),
                    rooms: make(map[*Room]bool),}
}

type Client struct {
    conn *websocket.Conn 
    wsServer *WsServer  
    send chan []byte
    rooms map[*Room]bool
}
func(c *Client)read(){
      defer c.conn.Close()
      c.conn.SetReadLimit(maxMessageSize)
      c.conn.SetReadDeadline(time.Now().Add(pongWait))
      c.conn.SetPongHandler(func (s string)error  {
         
          c.conn.SetReadDeadline(time.Now().Add(pongWait))
          return nil
      })

 
      for{
             _, jsonMsg, err := c.conn.ReadMessage()
           
             if err !=nil{
                  break
             }
           
            c.wsServer.forward <- jsonMsg
      }


}

func(c *Client)write(){
    ticker := time.NewTicker(pingPeriod)
    defer func(){
        ticker.Stop()
        c.conn.Close()
    }()
    c.conn.PingHandler()
 
    for{
        select {
                    case msg, ok:= <- c.send :
                                           
                        c.conn.SetWriteDeadline(time.Now().Add(writeWait))
                        if !ok {
                             c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                        }
                      
                        w, err := c.conn.NextWriter(websocket.TextMessage)
                        if err != nil {
                            return
                        }
                        w.Write(msg)
                        n := len(c.send)
                         for i := 0; i < n; i++ {
                            w.Write([]byte("/n"))
                            w.Write(<-c.send)
                        }
                        if err := w.Close(); err!= nil{ return}
                        case <-ticker.C:

                         
                        c.conn.SetWriteDeadline(time.Now().Add(writeWait))
                        if err := c.conn.WriteMessage(websocket.PongMessage, nil);err != nil{
                             return
                        }
                        
        }

    }

}

func(c *Client)disconnect(){
    c.wsServer.leave <- c
    for r := range c.rooms{
          r.leave <- c
    }
}


func Wshandler(w http.ResponseWriter, r *http.Request, ws *WsServer){
      conn, err := upgrader.Upgrade(w, r, nil)
      if err != nil{
           log.Println(err)
           return
      }
      client := NewClient(conn, ws)

                go client.read()
                go client.write()

      ws.join <-client     
  
}

 