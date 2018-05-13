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

var smtp_authinfo = &smtp_data{}
var rediscli *redis.Client
var config map[string]interface{}

func ReceiveEmail(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		var ei email_info
		code, mess := DecodeReqBody(req, &ei)
		if code != 0 {
			http.Error(rw, mess, code)
		} else {
			code, mess := ProcessEmail(ei)
			EncodeReqResp(&rw, code, mess)
		}
	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
	defer req.Body.Close()
}

func ConfirmEmail(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		ec := &email_confirm{}
		code, mess := DecodeReqBody(req, ec)
		if code != 0 {
			http.Error(rw, mess, code)
		} else {
			code, mess := GenEmailJWT(ec)
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
	InitEmailClient(smtp_authinfo)

	log.Println("Token manager starts on ", HttpLynxTokenManagerPort)
	http.HandleFunc("/reg-email", ReceiveEmail)
	http.HandleFunc("/confirm-email", ConfirmEmail)
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to token manager!")
	})
	http.ListenAndServe(HttpLynxTokenManagerPort, nil)
}
