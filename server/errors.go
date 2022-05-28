package server

import (
	"fmt"
)

type ErrUnknownCommand struct {
	cmd string
}

func (e ErrUnknownCommand) Error() string {
	return fmt.Sprintf("unknown command: %v", e.cmd)
}
