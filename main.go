package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"code.yarg.fr/kane/uchat/server"
)

func randString(length int) string {
	buf := make([]byte, (length+1)/2)

	if _, err := rand.Read(buf); err != nil {
		return ""
	}

	return hex.EncodeToString(buf)[:length]
}

func main() {
	var config Config
	config.Load()

	serverName := flag.String("name", randString(16), "The server name")
	serverAddress := flag.String("address", "0.0.0.0", "The server listen address")
	serverPort := flag.Int("port", 35498, "The server port")
	message := flag.String("msg", "", "The message to send")
	flag.Parse()

	config.SetName(*serverName)
	config.SetPort(uint16(*serverPort))
	config.SetAddress(*serverAddress)
	config.Save()

	if *message == "" {
		server := server.New(*serverName, *serverPort)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		go server.Serve(*serverAddress)

		<-ctx.Done()
		server.Close()
	} else {

		dst := "127.0.0.1"
		dstIP := net.ParseIP(dst)
		dstAddress := net.UDPAddr{IP: dstIP, Port: 35498}
		connection, err := net.DialUDP("udp", nil, &dstAddress)
		if err != nil {
			panic(err)
		}

		connection.Write([]byte(*message))

		buffer := make([]byte, 1200)
		_, err = connection.Read(buffer)
		if err != nil {
			panic(err)
		}

		fmt.Printf("received: %s\n", buffer)
	}
}
