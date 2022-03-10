package main


func NewWsServer() * WsServer{
    return &WsServer{
         forward: make(chan []byte),
         clients: make(map[*Client]bool),
         join: make(chan *Client),
         leave: make(chan *Client),
    }
}


type WsServer struct{
    forward chan []byte
    join  chan *Client
    leave chan *Client
    clients map[*Client] bool
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