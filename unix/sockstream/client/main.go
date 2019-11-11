package main

import (
	"net"
	"os"
)

func main() {
	laddr := net.UnixAddr{
		Name: "/tmp/unixdomaincli",
		Net:  "unix",
	}
	conn, err := net.DialUnix("unix", &laddr, &net.UnixAddr{
		Name: "/tmp/unixdomain",
		Net:  "unix",
	})
	if err != nil {
		panic(err)
	}
	defer os.Remove("/tmp/unixdomaincli")

	_, err = conn.Write([]byte("hello"))
	if err != nil {
		panic(err)
	}
	conn.Close()
}
