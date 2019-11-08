package main

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type timeoutHandler struct{}

func (h timeoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	timer := time.AfterFunc(5*time.Second, func() {
		r.Body.Close()
	})
	bodyBytes := make([]byte, 0)
	for {
		//We reset the timer, for the variable time
		timer.Reset(1 * time.Second)

		_, err := io.CopyN(bytes.NewBuffer(bodyBytes), r.Body, 256)
		if err == io.EOF {
			// This is not an error in the common sense
			// io.EOF tells us, that we did read the complete body
			break
		} else if err != nil {
			//You should do error handling here
			break
		}
	}
}

func main() {

	h := timeoutHandler{}
	server := &http.Server{
		// defines how long you allow a connection to be open during a client sends data
		ReadTimeout: 1 * time.Minute,
		// it is in the other direction
		WriteTimeout: 2 * time.Minute,
		// represents the time until the full request header (send by a client) should be read
		ReadHeaderTimeout: 20 * time.Second,
		Handler:           h,
		Addr:              ":8080",
	}

	server.ListenAndServe()
}
