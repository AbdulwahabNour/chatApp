package main

import (
	"encoding/json"
	"log"
)

const (
	SendMessageAction     = "send-message"
	JoinRoomActon         = "join-room"
	LeaveRoomAction       = "leave-room"
	UserJoinedAction      = "user-join"
	UserLeftAction        = "user-left"
	JoinRoomPrivateAction = "join-room-private"
 
)

type Message struct {
	Action  string    `json:"action"`
	Message string    `json:"message"`
	Target  *Room     `json:"target"`
	Sender  []*Client `jsoin:"sender"`

}
 

func (m *Message) encode() []byte {

	json, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}
	return json
}
func (m *Message) decode(data []byte) error {
	err := json.Unmarshal(data, m)
	if err != nil {
		return err
	}
	return nil 
}