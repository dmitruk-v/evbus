package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/google/uuid"
)

type server struct {
	log       *log.Logger
	addr      string
	listener  net.Listener
	clients   map[uuid.UUID]*client
	messageCh chan *message
	squitCh   chan struct{}
	wg        *sync.WaitGroup
}

func New(addr string, log *log.Logger) *server {
	return &server{
		log:       log,
		addr:      addr,
		clients:   make(map[uuid.UUID]*client),
		messageCh: make(chan *message),
		squitCh:   make(chan struct{}),
		wg:        &sync.WaitGroup{},
	}
}

func (s *server) Listen() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = listener
	s.log.Printf("Event-bus server listens at %v\n", s.addr)

	s.wg.Add(1)
	go func() {
		s.run()
		s.wg.Done()
	}()

	s.wg.Add(1)
	go func() {
		s.serve()
		s.wg.Done()
	}()

	s.wg.Wait()
	return nil
}

func (s *server) serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.squitCh:
				return
			default:
				s.log.Println(err)
			}
		} else {
			s.log.Printf("%v: Connected\n", conn.RemoteAddr())
			client := s.newClient(conn)
			s.wg.Add(1)
			go func() {
				client.handle()
				s.wg.Done()
			}()
		}
	}
}

func (s *server) run() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT)
	for {
		select {
		case <-sigCh:
			s.close()
			return
		case msg := <-s.messageCh:
			if err := s.processMessage(msg); err != nil {
				s.log.Println(err)
			}
		}
	}
}

func (s *server) processMessage(msg *message) error {
	switch msg.Cmd {
	case CmdEmit:
		s.emit(msg)
	case CmdSubscribe:
		s.subscribe(msg)
	case CmdDisconnect:
		s.disconnectClient(msg.Client)
	default:
		return ErrUnknownCommand{cmd: msg.Cmd}
	}
	return nil
}

func (s *server) subscribe(msg *message) {
	msg.Client.topics[msg.Topic] = true
	s.log.Printf("%v: Subscribed to: %v\n", msg.Client.conn.RemoteAddr(), msg.Topic)
}

func (s *server) emit(msg *message) {
	var clients []string
	for _, cl := range s.clients {
		if _, ok := cl.topics[msg.Topic]; !ok {
			continue
		}
		cl.broadcastCh <- msg
		clients = append(clients, cl.conn.RemoteAddr().String())
	}
	s.log.Printf("%v: Emitted: %v -> %v\n", msg.Client.conn.RemoteAddr(), msg.Topic, strings.Join(clients, ", "))
}

func (s *server) disconnectClient(c *client) {
	delete(s.clients, c.id)
	close(c.cquitCh)
	if err := c.conn.Close(); err != nil {
		s.log.Println(err)
	}
	s.log.Printf("%v: Disconnected\n", c.conn.RemoteAddr())
}

func (s *server) close() {
	close(s.squitCh)
	if err := s.listener.Close(); err != nil {
		s.log.Println(err)
	}
	for _, c := range s.clients {
		s.disconnectClient(c)
	}
}

func (s *server) newClient(conn net.Conn) *client {
	id := uuid.New()
	client := &client{
		id:          id,
		conn:        conn,
		log:         s.log,
		broadcastCh: make(chan *message),
		messageCh:   s.messageCh,
		topics:      make(map[string]bool),
		squitCh:     s.squitCh,
		cquitCh:     make(chan struct{}),
		wg:          sync.WaitGroup{},
	}
	s.clients[id] = client
	return client
}
