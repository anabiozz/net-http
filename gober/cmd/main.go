package main

import (
	"flag"
	"log"
	"net"
	"net-http/gober/client"
	"net-http/gober/server"
	"time"

	"github.com/pkg/errors"
)

func main() {
	connect := flag.String("connect", "", "IP address of process to join. If empty, go into listen mode.")
	flag.Parse()

	dialer := &net.Dialer{
		Timeout:   3 * time.Minute,
		KeepAlive: 3 * time.Minute,
		LocalAddr: &net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 0,
		},
	}

	if *connect != "" {
		err := client.Client(*connect, dialer)
		if err != nil {
			log.Println("Error:", errors.WithStack(err))
		}
		log.Println("Client done.")
		return
	}

	err := server.Server()
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}

	log.Println("Server done.")
}

func init() {
	log.SetFlags(log.Lshortfile)
}
