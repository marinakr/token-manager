package main

import (
	"net/http"
	"fmt"
	"log"
)

const (
	HttpLynxTokenManagerPort 	  = ":7665"
)


func receive_email(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		var email email_info
		code, mess := DecodeReqBody(req, &email)
		if code != 0 {
			http.Error(rw, mess, code)
		} else {
			err = process_email(email)
			rw.WriteHeader(status_code)
			rw.Write(status_line)
		}
	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
	defer req.Body.Close()
}

func main() {

	InitRedisClient()

	log.Println("Token manager starts on ", HttpLynxTokenManagerPort)
	http.HandleFunc("/email", receive_email)
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to token manager!")
	})
	http.ListenAndServe(HttpLynxTokenManagerPort, nil)
}
