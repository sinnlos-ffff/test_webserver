package main

import "net"
import "fmt"
import "log"

type Message struct {
	source string
	payload []byte
}

type Server struct {
	address string
	listener net.Listener
	quitChannel     chan struct{}
	messageChannel chan Message
}

func NewServer(address string) *Server {
	return &Server{
		address:    address,
		quitChannel:    make(chan struct{}),
		// without buffer, the channel will block the readLoop
		messageChannel: make(chan Message, 10),
	}
}

func (server *Server) acceptLoop() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("Accepted connection from: ", conn.RemoteAddr())

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

		server.messageChannel <- Message{
			source: conn.RemoteAddr().String(),
			payload: buf[:n],
		}
	}
}

func (server *Server) Start() error {
	listener, err := net.Listen("tcp", server.address)
	if err != nil {
		return err
	}

	defer listener.Close()
	server.listener = listener

	go server.acceptLoop()

	// blocks return until quitChannel is closed
	<-server.quitChannel
	close(server.messageChannel)

	return nil
}

func main() {
	server := NewServer(":8000")

	go func() {
		for message := range server.messageChannel {
			fmt.Printf("Received message from %s: %s\n", message.source, string(message.payload))
		}
	}()

	log.Fatal(server.Start())
}
