package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type EventBusClient struct {
	log         *log.Logger
	saddr       string
	conn        net.Conn
	subs        map[string]SubscribeHandler
	msgCh       chan *inMessage
	execCh      chan executionData
	herrCh      chan error // handlers errors
	isConnected bool
}

type SubscribeHandler func(data []byte) error

type executionData struct {
	msg     *inMessage
	handler SubscribeHandler
}

func New(saddr string) *EventBusClient {
	conn, err := net.Dial("tcp", saddr)
	if err != nil {
		panic(fmt.Sprintf("[EventBus]: connecting to event-bus server at %v, with error: %v", saddr, err))
	}
	c := &EventBusClient{
		log:         log.Default(),
		saddr:       saddr,
		conn:        conn,
		subs:        make(map[string]SubscribeHandler),
		msgCh:       make(chan *inMessage),
		execCh:      make(chan executionData),
		herrCh:      make(chan error),
		isConnected: true,
	}
	c.log.Printf("[EventBus]: client connected to the server at %v\n", c.conn.RemoteAddr())
	return c
}

func (c *EventBusClient) Run() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT)

	// Read from server
	go func() {
		defer c.disconnect()
		decoder := json.NewDecoder(c.conn)
		for {
			msg := &inMessage{}
			if err := decoder.Decode(msg); err != nil {
				// If there is an error or server closes connection - exit client
				c.log.Printf("[EventBus]: %v\n", err)
				return
			}
			c.msgCh <- msg
		}
	}()

	// Sequential handlers execution
	go func() {
		for edata := range c.execCh {
			if err := edata.handler([]byte(edata.msg.Data)); err != nil {
				c.herrCh <- err
			}
		}
	}()

	// Main processing
	for {
		select {
		case msg := <-c.msgCh:
			c.execute(msg)
		case <-sigCh:
			c.Disconnect()
		case herr, ok := <-c.herrCh:
			if !ok {
				return
			}
			c.log.Printf("[EventBus]: [Handler error]: %v\n", herr)
		}
	}
}

func (c *EventBusClient) execute(msg *inMessage) {
	handler, ok := c.subs[msg.Topic]
	if !ok {
		return
	}
	c.execCh <- executionData{
		msg:     msg,
		handler: handler,
	}
}

func (c *EventBusClient) SetLogger(log *log.Logger) {
	c.log = log
}

func (c *EventBusClient) Emit(topic string, data interface{}) {
	if !c.isConnected {
		c.log.Println("[EventBus]: trying to emmit, but client is not connected to the server")
		return
	}
	msg := outMessage{
		message: message{Cmd: CmdEmit, Topic: topic},
		Data:    data,
	}
	if err := json.NewEncoder(c.conn).Encode(msg); err != nil {
		c.log.Printf("[EventBus]: send emit message %v to the server, with error: %v\n", msg, err)
		return
	}
}

func (c *EventBusClient) Subscribe(topic string, handler SubscribeHandler) {
	if !c.isConnected {
		c.log.Println("[EventBus]: trying to subscribe, but client is not connected to the server")
		return
	}
	msg := outMessage{
		message: message{Cmd: CmdSubscribe, Topic: topic},
	}
	if err := json.NewEncoder(c.conn).Encode(msg); err != nil {
		c.log.Printf("[EventBus]: send subscribe message %v to the server, with error: %v\n", msg, err)
		return
	}
	c.subs[msg.Topic] = handler
}

func (c *EventBusClient) Sync() {
	if !c.isConnected {
		c.log.Println("[EventBus]: trying to sync, but client is not connected to the server")
		return
	}
	msg := outMessage{
		message: message{Cmd: CmdSync},
	}
	if err := json.NewEncoder(c.conn).Encode(msg); err != nil {
		c.log.Printf("[EventBus]: send sync message %v to the server, with error: %v\n", msg, err)
		return
	}
}

// Client initiated disconnect
func (c *EventBusClient) Disconnect() {
	defer c.conn.Close()
	if !c.isConnected {
		c.log.Println("[EventBus]: trying to disconnect, but client not connected to the server")
		return
	}
	msg := outMessage{
		message: message{Cmd: CmdDisconnect},
	}
	if err := json.NewEncoder(c.conn).Encode(msg); err != nil {
		c.log.Printf("[EventBus]: send disconnect message %v to the server, with error: %v\n", msg, err)
		return
	}
}

// Server initialted disconnect
func (c *EventBusClient) disconnect() {
	close(c.execCh)
	close(c.herrCh)
	c.isConnected = false
	c.log.Printf("[EventBus]: client disconnected from server at %v\n", c.conn.RemoteAddr())
}
