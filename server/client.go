package server

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
)

// type clientError struct {
// 	client *client
// 	err    error
// }

type client struct {
	id          uuid.UUID
	conn        net.Conn
	log         *log.Logger
	broadcastCh chan *message
	messageCh   chan<- *message
	topics      map[string]bool
	squitCh     chan struct{}
	cquitCh     chan struct{}
	wg          sync.WaitGroup
}

// Spawn two goroutines. One for reading from connection and another
// for writing to connection. Both parts will exit when we close client quit channel.
func (c *client) handle() {
	c.wg.Add(1)
	go func() {
		c.read()
		c.wg.Done()
	}()
	c.wg.Add(1)
	go func() {
		c.write()
		c.wg.Done()
	}()
	c.wg.Wait()
}

func (c *client) read() {
	decoder := json.NewDecoder(c.conn)
	for {
		select {
		case <-c.cquitCh:
			return
		default:
			msg := &message{}
			if err := decoder.Decode(msg); err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				c.log.Println(err)
				return
			}
			msg.Client = c
			c.messageCh <- msg
		}
	}
}

func (c *client) write() {
	encoder := json.NewEncoder(c.conn)
	for {
		select {
		case <-c.cquitCh:
			return
		case msg := <-c.broadcastCh:
			if err := encoder.Encode(msg); err != nil {
				c.log.Println(err)
				return
			}
		}
	}
}
