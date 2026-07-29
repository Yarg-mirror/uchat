package server

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type Server struct {
	name string
	ip   net.IP
	port int
	conn *net.UDPConn
}

func New(name string, port int) *Server {
	s := Server{name: name, port: port}

	return &s
}

func (s *Server) Serve(address string) {
	s.ip = net.ParseIP(address)
	localeAddress := net.UDPAddr{IP: s.ip, Port: s.port}

	connection, err := net.ListenUDP("udp", &localeAddress)
	if err != nil {
		log.Fatalf("Unable to list on port %d: %v", s.port, err)
	}
	log.Printf("UDP server listening on %s:%d", s.ip, s.port)
	s.conn = connection

	buffer := make([]byte, 1200)
	for {
		n, src, err := connection.ReadFromUDP(buffer)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			log.Printf("Reading error : %v", err)
			continue
		}
		s.process(buffer[:n], src)
	}
}

func (s *Server) Close() {
	s.conn.Close()
}

func (s *Server) process(message []byte, src *net.UDPAddr) {
	msg := string(message)

	switch msg {
	case "hello":
		response := fmt.Sprintf("Hello, my name is %s", s.name)
		_, err := s.conn.WriteToUDP([]byte(response), src)
		if err != nil {
			log.Printf("Error while sending to %s : %v", src, err)
		}

	case "bye":
		response := fmt.Sprintln("Good bye")
		_, err := s.conn.WriteToUDP([]byte(response), src)
		if err != nil {
			log.Printf("Error while sending to %s : %v", src, err)
		}

	default:
		_, err := s.conn.WriteToUDP(message, src)
		if err != nil {
			log.Printf("Error while sending to %s : %v", src, err)
		}
	}
}
