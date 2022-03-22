package server

import (
	"fmt"
	"goChallenge/chat/config"
	"log"
	"net/http"
	"time"
)

var RoomsMessages map[string][]string

type Server struct {
	hubs    map[string]*Hub
	options chan Option
}

func NewServer() *Server {
	return &Server{
		hubs:    make(map[string]*Hub),
		options: make(chan Option),
	}
}

func (s *Server) GetHubs() *map[string]*Hub {
	return &s.hubs
}

func (s *Server) Run() {
	for cmd := range s.options {
		switch cmd.ID {
		case OPT_JOIN:
			s.Join(cmd.Client, cmd.Argument)
		case OPT_QUIT:
			s.QuitCurrentRoom(cmd.Client, cmd.Argument)
		}
	}
}

// serveWs handles websocket requests from the peer.
func (s *Server) ServeWs(w http.ResponseWriter, r *http.Request) {

	// claims, _ := controllers.ExtractTokenMetadata(r.Cookies()[0].Value)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	if len(s.hubs) == 0 {
		initializeHubs(s)
	}

	userName, err := r.Cookie("userName")
	if err != nil {
		userName.Value = "guest"
	}

	client := &Client{hub: s.hubs["#Public"],
		nick:    userName.Value,
		conn:    conn,
		send:    make(chan []byte, 256),
		options: s.options,
	}

	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func initializeHubs(s *Server) {
	hub := newHub()
	hub.name = "#Public"
	s.hubs[hub.name] = hub
	go s.hubs[hub.name].run()

	hub1 := newHub()
	hub1.name = "#Teams"
	s.hubs[hub1.name] = hub1
	go s.hubs[hub1.name].run()

	hub2 := newHub()
	hub2.name = "#Zoom"
	s.hubs[hub2.name] = hub2
	go s.hubs[hub2.name].run()
}

func (s *Server) Join(c *Client, argument string) {
	h, ok := s.hubs[argument]
	if !ok {
		hub := newHub()
		hub.name = argument
		s.hubs[hub.name] = hub
		c.hub = hub
		h = hub

		go s.hubs[hub.name].run()
	}

	h.clients[c] = true
	c.hub = h
	c.hub.register <- c

	plainMsg := fmt.Sprintf("%v joined the room %s", c.nick, argument)

	c.hub.broadcast <- []byte(plainMsg)

	hours, minutes, _ := time.Now().Clock()
	msg := fmt.Sprintf("%d:%02d - %s", hours, minutes, plainMsg)
	AddCurrentMessages(RoomsMessages, c.hub.name, msg)
}

func (s *Server) QuitCurrentRoom(c *Client, arg string) {
	if c.hub != nil {
		c.hub.clients[c] = false

		plainMsg := "Someone has left the room"
		c.hub.broadcast <- []byte(plainMsg)

		hours, minutes, _ := time.Now().Clock()
		msg := fmt.Sprintf("%d:%02d - %s", hours, minutes, plainMsg)
		AddCurrentMessages(RoomsMessages, c.hub.name, string(msg))
	}
}

func (server *Server) Start() {
	router := loadRoutes(server)
	err := http.ListenAndServe(":"+config.AppConfig.Port, router)

	if err != nil {
		fmt.Println(err)
		log.Fatalf("Server crashed!. Error: %v", err)
	}
}
