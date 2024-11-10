package main

import "net"
import "fmt"

type Server struct {
	address string
	listener net.Listener
	quitChannel     chan struct{}
}

func NewServer(address string) *Server {
	return &Server{
		address:    address,
		quitChannel: make(chan struct{}),
	}
}

func (server *Server) acceptLoop() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go server.readLoop(conn)
	}
}

func (server *Server) readLoop(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			continue
		}

		message := buf[:n]
		fmt.Println(string(message))
	}
}

func (server *Server) Start() error {
	listener, err := net.Listen("tcp", server.address)
	if err != nil {
		return err
	}

	defer listener.Close()
	server.listener = listener

	// blocks return until quitChannel is closed
	<-server.quitChannel

	return nil
}

func main() {
}
