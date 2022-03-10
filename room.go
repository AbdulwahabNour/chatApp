package main


func NewRoom(name string) *Room{
     return &Room{
          name: name,
          clients: make(map[*Client]bool),
          join: make(chan *Client),
          leave: make(chan *Client),
          forward: make(chan *Message),
     }
}
 
type Room struct {
     name string
     clients map[*Client]bool
     join chan *Client
     leave chan *Client
     forward chan *Message
}

func (r *Room)RunRoom(){
       for{
            select{
                  case c := <- r.join :
                      r.addClientInRoom(c)
            }
       }
}

func (r *Room)addClientInRoom(c *Client){
            r.notifyClientJoined(c)
            r.clients[c]= true
}
func (r *Room)removeClientFromRoom(c *Client){
           if _, ok := r.clients[c];ok{
               delete(r.clients, c)
           }
}
func (r *Room) forwardToClientsInRoom(msg []byte){
    for c := range r.clients{
        c.send <- msg
    }
}

func (r *Room)notifyClientJoined(c * Client){

}

