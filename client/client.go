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
	log   *log.Logger
	saddr string
	conn  net.Conn
	subs  map[string]SubscribeHandler
}

type SubscribeHandler func(data []byte) error

func newClient(saddr string, conn net.Conn) *EventBusClient {
	return &EventBusClient{
		log:   log.Default(),
		saddr: saddr,
		conn:  conn,
		subs:  make(map[string]SubscribeHandler),
	}
}

func Connect(saddr string) (*EventBusClient, error) {
	conn, err := net.Dial("tcp", saddr)
	if err != nil {
		return nil, err
	}
	c := newClient(saddr, conn)
	c.log.Printf("Connected to %v\n", c.saddr)
	return c, nil
}

func (ebc *EventBusClient) Emit(topic string, data interface{}) error {
	msg := outMessage{
		message: message{Cmd: CmdEmit, Topic: topic},
		Data:    data,
	}
	if err := json.NewEncoder(ebc.conn).Encode(msg); err != nil {
		return fmt.Errorf("send emit message %v to the server, with error: %v", msg, err)
	}
	return nil
}

func (ebc *EventBusClient) Subscribe(topic string, handler SubscribeHandler) error {
	msg := outMessage{
		message: message{Cmd: CmdSubscribe, Topic: topic},
	}
	if err := json.NewEncoder(ebc.conn).Encode(msg); err != nil {
		return fmt.Errorf("send subscribe message %v to the server, with error: %v", msg, err)
	}
	ebc.subs[msg.Topic] = handler
	return nil
}

func (ebc *EventBusClient) Disconnect() error {
	defer ebc.conn.Close()
	msg := outMessage{
		message: message{Cmd: CmdDisconnect},
	}
	if err := json.NewEncoder(ebc.conn).Encode(msg); err != nil {
		return fmt.Errorf("send disconnect message %v to the server, with error: %v", msg, err)
	}
	fmt.Println("--- done! ---")
	return nil
}

func (ebc *EventBusClient) Run() {
	errCh := make(chan error)
	go func() {
		decoder := json.NewDecoder(ebc.conn)
		for {
			msg := &inMessage{}
			if err := decoder.Decode(msg); err != nil {
				// If there is an error or server closes connection - exit from goroutine
				errCh <- fmt.Errorf("[ERROR]: %v", err)
				return
			}
			handler, ok := ebc.subs[msg.Topic]
			if !ok {
				continue
			}
			go func() {
				if err := handler([]byte(msg.Data)); err != nil {
					errCh <- fmt.Errorf("[HANDLER ERROR]: %v", err)
				}
			}()
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT)
	for {
		select {
		case err := <-errCh:
			ebc.log.Println(err)
		case <-sigCh:
			if err := ebc.Disconnect(); err != nil {
				ebc.log.Println(err)
			}
			ebc.log.Printf("Disconnected from %v\n", ebc.saddr)
			return
		}
	}
}

func (ebc *EventBusClient) SetLogger(log *log.Logger) {
	ebc.log = log
}
