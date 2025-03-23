package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}
func (s *Server) handleWs(ws *websocket.Conn) {
	fmt.Println("New connection", ws.RemoteAddr().String())

	s.conns[ws] = true

	s.readLoop(ws)
}
func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 10000)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read err", err)
			continue
		}
		msg := buf[:n]
		s.broadcast(msg, n)
	}
}
func (s *Server) broadcast(b []byte, n int) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b[0:n]); err != nil {
				fmt.Println("write err", err)
			}
			println(string(b[0:n]))
		}(ws)
	}

}

func main() {
	server := NewServer()

	http.Handle("/ws", websocket.Handler(server.handleWs))
	http.ListenAndServe(":7777", nil)
}
