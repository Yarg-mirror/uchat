package server

import (
	"errors"
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
		log.Fatalf("Impossible d'écouter sur le port %d: %v", s.port, err)
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
			log.Printf("Erreur de lecture : %v", err)
			continue
		}
		log.Printf("Réception: %s", buffer[:n])

		_, err = connection.WriteToUDP(buffer[:n], src)
		if err != nil {
			log.Printf("Erreur d'envoi vers %s : %v", src, err)
		}
	}
}

func (s *Server) Close() {
	s.conn.Close()
}
