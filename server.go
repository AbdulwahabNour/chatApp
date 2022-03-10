package main


func NewWsServer() * WsServer{
    return &WsServer{
         forward: make(chan []byte),
         clients: make(map[*Client]bool),
         join: make(chan *Client),
         leave: make(chan *Client),
         rooms: map[*Room]bool{},
    }
}


type WsServer struct{
    forward chan []byte
    join  chan *Client
    leave chan *Client
    clients map[*Client] bool
    rooms map[*Room] bool
}

func(s *WsServer)Run(){
      for{
           select{
           case c := <- s.join:
              s.clients[c] = true
           case c := <-s.leave:
             if _, ok := s.clients[c]; ok{
                 delete(s.clients, c)
             }
            case msg := <- s.forward:
                  
                    s.forwardMsg(msg)
           }
      }
}
func (s *WsServer) forwardMsg(msg []byte){
   
     for c := range s.clients{
       
          c.send <- msg
     }
}
func (s *WsServer)findRoom(name string) *Room{
     var room *Room

     for r := range s.rooms{
         if r.name == name{
             room = r
             break
         }
     }
     return room
}
func (s *WsServer)createRoom(name string) *Room{
    room := NewRoom(name)

    go room.RunRoom()

    s.rooms[room]=true
    return room
}