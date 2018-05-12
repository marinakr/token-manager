package main

import (
	"net/http"
	"fmt"
	"log"
	"github.com/go-redis/redis"
)

const (
	HttpLynxTokenManagerPort 	  = ":7665"
)

var smtpauthinfo = &smtp_data{}
var rediscli *redis.Client
var config map[string]interface{}

func receive_email(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		var email email_info
		code, mess := DecodeReqBody(req, &email)
		if code != 0 {
			http.Error(rw, mess, code)
		} else {
			code, mess := ProcessEmail(email)
			EncodeReqResp(&rw, code, mess)
		}
	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
	defer req.Body.Close()
}

func main() {
	ReadConfig(&config)
	rediscli = InitRedisClient()
	InitEmailClient(smtpauthinfo)

	log.Println("Token manager starts on ", HttpLynxTokenManagerPort)
	http.HandleFunc("/email", receive_email)
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to token manager!")
	})
	http.ListenAndServe(HttpLynxTokenManagerPort, nil)
}
