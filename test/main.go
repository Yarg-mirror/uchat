package main

import (
	"context"
	"fmt"
	"net"
	"net/netip"
)

type Bootstrap struct {
	Address string
	Port    uint16
}

func main() {
	bootstrap := []Bootstrap{
		{Address: "yarg.fr", Port: 3000},
		{Address: "yarh.fr", Port: 3000},
		{Address: "192.168.1.114", Port: 3000},
	}

	var peers []netip.AddrPort

	for _, item := range bootstrap {
		ips, _ := net.DefaultResolver.LookupNetIP(
			context.Background(),
			"ip",
			item.Address,
		)

		for _, ip := range ips {
			addr := netip.AddrPortFrom(ip.Unmap(), item.Port)
			peers = append(peers, addr)

			fmt.Println(item.Address, addr)
		}
	}
}
