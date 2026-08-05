package server

import "net"

type Peer struct {
	publicKey string
	name      string
	ip        net.IP
	port      int
	status    int
}
