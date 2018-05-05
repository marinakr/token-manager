package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
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

type email_info struct {
	Email string
	NickName string
}

func receive_email(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var email email_info
	err := decoder.Decode(&email)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	log.Println(email.Email)
	log.Println(email.NickName)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Success"))
}

func main() {
	fmt.Println("Token manager starts on ", HttpLynxTokenManagerPort)

	http.HandleFunc("/email", receive_email)

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to token manager!")
	})
	http.ListenAndServe(HttpLynxTokenManagerPort, nil)

}
