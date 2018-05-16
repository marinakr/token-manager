package main

import (
	//
	"fmt"
	"net/http"
	"./conf"
	"./redscli"
	"./sender"
	"./utils"
	"./reg"
	//
)

type Env struct {
	dbcli redscli.RedisENV
	smtpcli sender.SmtpENV
}

func (env *Env)ReceiveEmail(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		ei, err := reg.VailidateRegData(req)
		if err == nil {
			err = ei.CheckDBRegData(env.dbcli)
			if err == nil {

			} else {
				http.Error(rw, err.Error(), http.StatusConflict)
			}
		} else {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		}
	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
	defer req.Body.Close()
}

func (env *Env)ReceiveEmailCode(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		ec := &EmailConfirm{}
		code, mess := utils.DecodeReqBody(req, ec)
		if code != Ok {
			http.Error(rw, mess, code)
		} else {
			code, mess := ConfirmEmail(ec)
			payload := utils.GenResponsePayload(code, mess)
			EncodeReqResp(rw, http.StatusOK, payload)
		}
	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
	defer req.Body.Close()
}

func main() {
	config := conf.NewConig()
	dbclient := redscli.New(config.RedisConf())
	smtpsender := sender.NewEmailSender(config.EmailConf())
	env := &Env{dbcli: &dbclient, smtpcli: smtpsender}

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
