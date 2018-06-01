package main

import (
	//
	"fmt"
	"net/http"
	"token-manager/conf"
	"token-manager/redscli"
	"token-manager/sender"
	"token-manager/reg"
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
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		} else {
			err = ei.CheckDBRegData(env.dbcli)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusConflict)
			} else {
				err = ei.RegisterEmail(env.dbcli, env.smtpcli)
				if err != nil{
					http.Error(rw, err.Error(), http.StatusBadRequest)
				} else {
					rw.WriteHeader(http.StatusOK)
				}
			}
		}
	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
	defer req.Body.Close()
}

func (env *Env)ReceiveEmailCode(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		ec, err := reg.ValidateEmailConfirm(req)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		} else {
			jwt, err := ec.CheckEmailConfirm(env.dbcli)
			if err != nil{
				http.Error(rw, err.Error(), http.StatusBadRequest)
			} else {
				fmt.Fprint(rw, jwt)
				rw.WriteHeader(http.StatusOK)
			}
		}
	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
	defer req.Body.Close()
}

func main() {
	config := conf.NewConig()
	dbclient := redscli.New(config.RedisConf())
	smtpsender := sender.New(config.EmailConf())
	env := &Env{dbcli: dbclient, smtpcli: smtpsender}

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
