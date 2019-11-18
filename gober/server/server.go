package server

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
	"net"
	"net-http/gober"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// HandleFunc ..
type HandleFunc func(*bufio.ReadWriter)

// Endpoint ..
type Endpoint struct {
	listener net.Listener
	handler  map[string]HandleFunc

	m sync.RWMutex
}

// NewEndpoint ..
func NewEndpoint() *Endpoint {
	return &Endpoint{
		handler: map[string]HandleFunc{},
	}
}

// AddHandleFunc ..
func (e *Endpoint) AddHandleFunc(name string, f HandleFunc) {
	e.m.Lock()
	e.handler[name] = f
	e.m.Unlock()
}

// Listen ..
func (e *Endpoint) Listen() error {
	var err error
	e.listener, err = net.Listen("tcp", gober.Port)
	if err != nil {
		return errors.Wrapf(err, "Unable to listen on port %s\n", gober.Port)
	}
	log.Println("Listen on", e.listener.Addr().String())
	for {
		log.Println("Accept a connection request.")
		conn, err := e.listener.Accept()
		if err != nil {
			log.Println("Failed accepting a connection request: ", err)
			continue
		}
		log.Println("Handle incoming messages.")
		go e.handleMessages(conn)
	}
}

func (e *Endpoint) handleMessages(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	for {
		log.Println("Receive command '")
		cmd, err := rw.ReadString('\n')
		switch {
		case err == io.EOF:
			log.Println("Reached EOF - close this connectioon. \n ---")
			return
		case err != nil:
			log.Println("\nError reading command. Got : '"+cmd+"'\n", err)
			return
		}

		cmd = strings.Trim(cmd, "\n ")
		log.Println(cmd + "'")
		e.m.RLock()
		handleCommand, ok := e.handler[cmd]
		e.m.RUnlock()
		if !ok {
			log.Println("Command '" + cmd + "' is not registered.")
			return
		}
		handleCommand(rw)
	}
}

func handleStrings(rw *bufio.ReadWriter) {
	log.Print("Receive STRING message:")
	s, err := rw.ReadString('\n')
	if err != nil {
		log.Println("Cannot read from connection.\n", err)
	}
	s = strings.Trim(s, "\n ")
	log.Println(s)
	_, err = rw.WriteString("Thank you.\n")
	if err != nil {
		log.Println("Cannot write to connection.\n", err)
	}
	err = rw.Flush()
	if err != nil {
		log.Println("Flush failed.", err)
	}
}

func handleGob(rw *bufio.ReadWriter) {
	log.Print("Receive GOB data:")
	var data gober.ComplexData
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&data)
	if err != nil {
		log.Println("Error decoding GOB data:", err)
		return
	}
	log.Printf("Outer complexData struct: \n%#v\n", data)
	log.Printf("Inner complexData struct: \n%#v\n", data.C)
}

// Server ..
func Server() error {
	endpoint := NewEndpoint()
	endpoint.AddHandleFunc("STRING", handleStrings)
	endpoint.AddHandleFunc("GOB", handleGob)
	return endpoint.Listen()
}
