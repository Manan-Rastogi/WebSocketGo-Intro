package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

// Server with multiple connections. Eg: chat Appln can have multiple connections active
type Server struct {
	conns map[*websocket.Conn]bool
	mu    sync.Mutex
}

// NewServer function Creates a new server and Return this server.
func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

// Will Read the incoming messages till the connection exists.
func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)

	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF { // EOF means that the connection has been closed. So no Incoming message
				break
			}

			log.Println("Error reading from connection:", err)
			continue
		}

		msg := buf[:n]
		
		s.broadcast(msg)
	}
}

// A handler for WebSocket
func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New incoming connection from client - ", ws.RemoteAddr())

	s.mu.Lock()
	s.conns[ws] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.conns, ws)
		s.mu.Unlock()
		ws.Close()
		fmt.Println("Connection closed from client - ", ws.RemoteAddr())
	}()

	s.readLoop(ws)
}

func (s *Server) broadcast(b []byte){
	for ws := range s.conns{
		go func(ws *websocket.Conn){
			if _, err := ws.Write(b); err != nil{
				fmt.Println("Error - ", err.Error())
			}
		}(ws)
	}
}


func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	fmt.Println("Starting server on :3100") // Log to confirm the server is starting
	if err := http.ListenAndServe(":3100", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
