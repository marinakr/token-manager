package main

import (
	//
	"fmt"
	"net/http"
	"./redscli"
	//
	"github.com/go-redis/redis"
	"token-manager/conf"
)

var rediscli *redis.Client
var config map[string]interface{}
var smtpdata = &SmtpData{}

type Env struct {
	dbcli redscli.RedisENV
}

func (env *Env)ReceiveEmail(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		var ei EmailInfo
		code, mess := DecodeReqBody(req, &ei)
		if code != Ok {
			http.Error(rw, mess, code)
		} else {
			code, mess := PrepareEmailCode(ei)
			payload := GenResponsePayload(code, mess)
			EncodeReqResp(rw, http.StatusOK, payload)
		}
	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
	defer req.Body.Close()
}

func (env *Env)ReceiveEmailCode(rw http.ResponseWriter, req *http.Request) {
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
	//read config//
	config := conf.NewConig()

	//cleints
	InitEmailClient(smtpdata)

	dbclient := redscli.New(config.RedisConf())
	env := &Env{dbcli: dbclient}

	//main app
	port := config.PortConf()
	fmt.Println("Token manager starts on ", port)
	http.HandleFunc("/reg-email", env.ReceiveEmail)
	http.HandleFunc("/confirm-email", env.ReceiveEmailCode)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to token manager!")
	})
	http.ListenAndServe(port, nil)
}
