package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type Message struct {
	msg string `json:"msg"`
}

type WSHub struct {
	connections      map[string]*websocket.Conn
	addClientChan    chan *websocket.Conn
	removeClientChan chan *websocket.Conn
	broadcastChan    chan Message
}

func newHub() *WSHub {
	return &WSHub{
		connections:      make(map[string]*websocket.Conn, 0),
		addClientChan:    make(chan *websocket.Conn),
		removeClientChan: make(chan *websocket.Conn),
		broadcastChan:    make(chan Message),
	}
}

func (h *WSHub) run() {
	for {
		select {
		case conn := <-h.addClientChan:
			h.addClient(conn)
		case conn := <-h.removeClientChan:
			h.removeClient(conn)
		case msg := <-h.broadcastChan:
			h.broadcastMessage(msg)
		}
	}
}

func (h *WSHub) addClient(conn *websocket.Conn) {
	h.connections[conn.RemoteAddr().String()] = conn
	// fmt.Println(h.connections)
}

func (h *WSHub) removeClient(ws *websocket.Conn) {
	delete(h.connections, ws.LocalAddr().String())
}

func (h *WSHub) broadcastMessage(msg Message) {
	fmt.Println("broadcasting message : ", msg)
	for _, conn := range h.connections {
		err := websocket.JSON.Send(conn, msg)
		if err != nil {
			fmt.Println("Error broadcasting message: ", err)
			return
		}
	}
}

func handle(ws *websocket.Conn, h *WSHub) {
	go h.run()
	h.addClientChan <- ws
	for {
		var m Message
		err := websocket.JSON.Receive(ws, &m)
		if err != nil {
			h.broadcastChan <- Message{msg: err.Error()}
			h.removeClientChan <- ws
			break
		}
		h.broadcastChan <- m
	}
}

func server(port string) error {
	h := newHub()
	mux := http.NewServeMux()
	mux.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		handle(ws, h)
	}))
	s := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	return s.ListenAndServe()
}

var port *string

func init() {
	port = flag.String("port", "9000", "specify port")
}

func main() {
	flag.Parse()
	log.Fatal(server(*port))
}
