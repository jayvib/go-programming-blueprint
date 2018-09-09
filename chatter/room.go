package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const (
	SocketBufferSize = 1024
	MessageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize: SocketBufferSize,
	WriteBufferSize: SocketBufferSize,
}

func NewRoom() *room {
	return &room{
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
	}
}

type room struct {
	forward chan []byte
	join chan *client
	leave chan *client
	clients map[*client]bool
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Panicln("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send: make(chan []byte, MessageBufferSize),
		room: r,
	}
	r.join <- client
	defer func() {
		r.leave <- client
	}()

	go client.write()
	client.read()
}
