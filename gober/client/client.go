package client

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net-http/gober"
	"strconv"

	"github.com/pkg/errors"
)

// Open ..
func Open(addr string, dialer *net.Dialer) (*bufio.ReadWriter, error) {
	fmt.Println("Dial " + addr)
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+"failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

// Client ..
func Client(ip string, dialer *net.Dialer) error {
	testStruct := gober.ComplexData{
		N: 23,
		S: "string data",
		M: map[string]int{"one": 1, "two": 2, "three": 3},
		P: []byte("abc"),
		C: &gober.ComplexData{
			N: 256,
			S: "Recursive structs? Piece of cake!",
			M: map[string]int{"01": 1, "10": 2, "11": 3},
		},
	}

	rw, err := Open(ip+gober.Port, dialer)
	if err != nil {
		return errors.Wrap(err, "Client: Failed to open connection to "+ip+gober.Port)
	}

	log.Println("Send the string request.")
	n, err := rw.WriteString("STRING\n")
	if err != nil {
		return errors.Wrap(err, "Could not send the STRING request ("+strconv.Itoa(n)+" bytes written)")
	}
	n, err = rw.WriteString("Additional data.\n")
	if err != nil {
		return errors.Wrap(err, "Could not send additional STRING data ("+strconv.Itoa(n)+" bytes written)")
	}
	log.Println("Flush the buffer.")
	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed.")
	}
	log.Println("Read the reply.")
	response, err := rw.ReadString('\n')
	if err != nil {
		return errors.Wrap(err, "Client: Failed to read the reply: '"+response+"'")
	}
	log.Println("STRING request: got a response:", response)
	log.Println("Send a struct as GOB:")
	log.Printf("Outer complexData struct: \n%#v\n", testStruct)
	log.Printf("Inner complexData struct: \n%#v\n", testStruct.C)
	enc := gob.NewEncoder(rw)
	n, err = rw.WriteString("GOB\n")
	if err != nil {
		return errors.Wrap(err, "Could not write GOB data ("+strconv.Itoa(n)+" bytes written)")
	}
	err = enc.Encode(testStruct)
	if err != nil {
		return errors.Wrapf(err, "Encode failed for struct: %#v", testStruct)
	}
	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed.")
	}
	return nil
}

// func main() {

// 	dialer := &net.Dialer{
// 		Timeout:   3 * time.Minute,
// 		KeepAlive: 3 * time.Minute,
// 		LocalAddr: &net.TCPAddr{
// 			IP:   net.ParseIP("127.0.0.1"),
// 			Port: 0,
// 		},
// 	}

// 	rw, err := Open("127.0.0.1:9595", dialer)

// 	rw.

// 	// fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")

// 	// for i := 0; i < 65535; i++ {
// 	// 	time.Sleep(1 * time.Second)
// 	// 	fmt.Println("connection run")

// 	// 	fmt.Fprintf(conn, "Host: www.nowhere123.com\r\n\r\n")

// 	// 	status, err := bufio.NewReader(conn).ReadString('\n')
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// 	fmt.Println("status", status)
// 	// }
// }
