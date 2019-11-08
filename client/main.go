package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{
		Transport: &http.Transport{
			// Defines the setup timeout for encripted HTTP connection
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			// Cares about the setup timout for upgrading the unencripted connection to an encripted one HTTPS
			TLSHandshakeTimeout: 10 * time.Second,
			//  This configures how long you want to wait after you send your payload for the beginning of an answer
			// (in form of the beginning of the header)
			ExpectContinueTimeout: 4 * time.Second,
			// And with this parameter you set how long the complete transfer of the header is allowed to last
			// So you want to have the complete header information ExpectContinueTimeout + ResponseHeaderTimeout
			// after your did send you complete request
			ResponseHeaderTimeout: 10 * time.Second,
		},
	}
	resp, err := client.Get("https://blog.simon-frey.eu/")
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
	}

	timer := time.AfterFunc(5*time.Second, func() {
		resp.Body.Close()
	})

	bodyBytes := make([]byte, 0)

	for {
		timer.Reset(1 * time.Second)

		_, err := io.CopyN(bytes.NewBuffer(bodyBytes), resp.Body, 256)
		if err != io.EOF {
			// This is not an error in the common sense
			// io.EOF tells us, that we did read the complete body
			break
		} else if err != nil {
			//You should do error handling here
			break
		}
	}

	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
	}

	fmt.Println(string(bodyBytes))
}
