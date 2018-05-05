package main

import (
	"fmt"
	"net/http"
)

const (
	//
	HttpLynxTokenManagerPort 	  = ":7665"
	// HTTP Responces
	StatusOK                   = 200
	StatusCreated              = 201
	StatusAccepted             = 202
	StatusNonAuthoritativeInfo = 203
	StatusNoContent            = 204
	StatusResetContent         = 205
	StatusPartialContent       = 206
)

func main() {
	fmt.Println("Token manager started on ", HttpLynxTokenManagerPort)

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to token manager!")
	})
	http.ListenAndServe(HttpLynxTokenManagerPort, nil)

}
