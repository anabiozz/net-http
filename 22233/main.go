package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
)

const (
	virtualHost = "example.com"
)

func checkResponse(ts *httptest.Server, client *http.Client) {
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		panic(err)
	}

	req.Host = virtualHost

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	_, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf(">> ERROR: Got unexpected response code: %d\n", resp.StatusCode)
	} else {
		fmt.Println(">> SUCCESS")
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	mux.Handle(virtualHost+"/", http.RedirectHandler("/landing", http.StatusFound))
	mux.HandleFunc(virtualHost+"/landing", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// This will fail due to the missing Host header
	checkResponse(ts, http.DefaultClient)
	// This hack succeeds
	checkResponse(ts, &http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}

			// request.Host won't be preserved following redirection, in order to
			// workaround this, when Location header in response is just a relative
			// path, just copy the Host to the current http.Request object.
			lastResp := r.Response
			loc := lastResp.Header.Get("Location")
			if u, err := url.Parse(loc); err == nil && u.Host == "" {
				r.Host = lastResp.Request.Host
			}

			return nil
		},
	})
}
