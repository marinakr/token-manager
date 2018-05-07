package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"regexp"
	"log"
	"time"
	"encoding/base64"
)

const (
	// App parameters
	HttpLynxTokenManagerPort 	  = ":7665"
	// Regexps
    EmilaRegExp = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

type email_info struct {
	Email string
	NickName string
}

func process_email(ei email_info)(status int, status_line []byte) {
	re := regexp.MustCompile(EmilaRegExp)
	if re.MatchString(string(ei.Email)) {
		log.Println(ei.Email)
		log.Println(ei.NickName)
		exp := time.Now().Unix() + 48*60*60
		response_txt := map[string]interface{}{
			"email" : ei.Email,
			"exp" : exp}
		lnk, err  := json.Marshal(response_txt)
		if err == nil {
			log.Printf("%s", lnk)
			status_line = []byte(base64.StdEncoding.EncodeToString(lnk))
			status = http.StatusOK
		}
	} else {
		status_line = []byte("Invalid email")
		status = http.StatusBadRequest
	}
	return
}

func receive_email(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		decoder := json.NewDecoder(req.Body)
		var email email_info
		err := decoder.Decode(&email)
		if err != nil {
			panic(err)
		} else {
			status_code, status_line := process_email(email)
			rw.WriteHeader(status_code)
			rw.Write(status_line)
		}
	default:
		status_code, status_line := http.StatusMethodNotAllowed, []byte("Method not allowed")
		rw.WriteHeader(status_code)
		rw.Write(status_line)
	}
	defer req.Body.Close()
}

func main() {
	log.Println("Token manager starts on ", HttpLynxTokenManagerPort)

	http.HandleFunc("/email", receive_email)

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to token manager!")
	})
	http.ListenAndServe(HttpLynxTokenManagerPort, nil)

}
