package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
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
	serverName := flag.String("name", randString(32), "The server name")
	serverAddress := flag.String("address", "0.0.0.0", "The server listen address")
	serverPort := flag.Int("port", 35498, "The server port")
	flag.Parse()
	server := server.New(*serverName, *serverPort)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go server.Serve(*serverAddress)

	<-ctx.Done()
	server.Close()
}
