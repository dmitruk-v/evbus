package server

import (
	"encoding/json"
	"time"
)

const (
	CmdEmit       = "Emit"
	CmdSubscribe  = "Subscribe"
	CmdDisconnect = "Disconnect"
	CmdSync       = "Sync"
)

type message struct {
	Client    *client         `json:"-"`
	Cmd       string          `json:"cmd"`
	Topic     string          `json:"topic"`
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"created_at"`
}
