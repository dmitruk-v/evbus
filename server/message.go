package server

import "encoding/json"

const (
	CmdEmit       = "Emit"
	CmdSubscribe  = "Subscribe"
	CmdDisconnect = "Disconnect"
)

type message struct {
	Client *client         `json:"-"`
	Cmd    string          `json:"cmd"`
	Topic  string          `json:"topic"`
	Data   json.RawMessage `json:"data"`
}
