package main

import (
	"net/http"
	"fmt"
	"github.com/go-redis/redis"
)

var smtp_authinfo = &smtp_data{}
var rediscli *redis.Client
var config map[string]interface{}

func AppPort()(port string){
	tm := config["token_manager"]
	token_manager := tm.(map[string]interface{})
	port = token_manager["port"].(string)
	return
}

func ReceiveEmail(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		var ei EmailInfo
		code, mess := DecodeReqBody(req, &ei)
		if code != Ok {
			http.Error(rw, mess, code)
		} else {
			code, mess := PrepareEmailCode(ei)
			EncodeReqResp(rw, code, mess)
		}
	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
	defer req.Body.Close()
}

func ReceiveEmailCode(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		ec := &EmailInfo{}
		code, mess := DecodeReqBody(req, ec)
		if code != Ok {
			http.Error(rw, mess, code)
		} else {
			code, mess := ConfirmEmail(ec)
			payload := GenResponsePayload(code, mess)
			EncodeReqResp(rw, http.StatusOK, payload)
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

	fmt.Println("Token manager starts on ",  AppPort())
	http.HandleFunc("/reg-email", ReceiveEmail)
	http.HandleFunc("/confirm-email", ReceiveEmailCode)
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to token manager!")
	})
	http.ListenAndServe( AppPort(), nil)
}
