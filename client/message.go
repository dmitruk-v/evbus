package client

import "encoding/json"

const (
	CmdEmit       = "Emit"
	CmdSubscribe  = "Subscribe"
	CmdDisconnect = "Disconnect"
	CmdSync       = "Sync"
)

type message struct {
	Cmd   string `json:"cmd"`
	Topic string `json:"topic"`
}

type outMessage struct {
	message
	Data interface{} `json:"data"`
}

type inMessage struct {
	message
	Data json.RawMessage `json:"data"`
}
