package main

import (
	"encoding/json"
	"log"
)


const(
     SendMessageAction = "send-message"
     JoinRoomActon     = "join-room"
     LeaveRoomAction   = "leave-room"
)

type Message struct{
      Action string `json:"action"`
      Message string `json:"message"`
      Target string `json:"target"`
      Sender *Client `jsoin:"sender"`
}

 
 
func(m *Message)encode()[]byte{

    json, err := json.Marshal(m)
    if err != nil{
        log.Println(err)
    }
    return json
}