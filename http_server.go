package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"regexp"
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
	StatusBadRequest		   = 400
	// Regexp
    EmilaRegExp = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

type email_info struct {
	Email string
	NickName string
}

func generate_response_bytes(code int)(status_line []byte){
	switch code {
	case StatusOK:
		status_line = []byte("Success")
	default:
		status_line = []byte("Invalid email")
	}
	return
}

func process_email(ei email_info)(status int){
	re := regexp.MustCompile(EmilaRegExp)
	if re.MatchString(ei.Email) {
		log.Println(ei.Email)
		log.Println(ei.NickName)
		status = StatusOK
	} else {
		status = StatusBadRequest
	}
	return
}

func receive_email(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var email email_info
	err := decoder.Decode(&email)
	if err != nil {
		panic(err)
	} else {
		status_code := process_email(email)
		status_line := generate_response_bytes(status_code)
		rw.WriteHeader(status_code)
		rw.Write(status_line)
	}
	defer req.Body.Close()
}

func main() {
	fmt.Println("Token manager starts on ", HttpLynxTokenManagerPort)

	http.HandleFunc("/email", receive_email)

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to token manager!")
	})
	http.ListenAndServe(HttpLynxTokenManagerPort, nil)

}
