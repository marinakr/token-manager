package sender

import (
	"net/smtp"
	"encoding/json"
	"fmt"
	"strconv"
)

type SmtpENV interface {
	SendEmail(email string, code int) error
}

type SmtpData struct {
	Username string
	Password string
	Host     string
	Smtphost string
	Identity string
	Subject  string
	Auth     smtp.Auth
}

func New(emailCreds interface{}) *SmtpData {
	smtpdata := &SmtpData{}
	data, err := json.Marshal(emailCreds)
	json.Unmarshal(data, smtpdata)
	if err == nil {
		smtpdata.Auth = smtp.PlainAuth(smtpdata.Identity, smtpdata.Username, smtpdata.Password, smtpdata.Smtphost)
		return  smtpdata
	} else {
		fmt.Println("Error smtp auth")
		panic(err)
	}
}

func (sm *SmtpData) SendEmail(email string, code int)(err error){
	err = smtp.SendMail(
		sm.Host,
		sm.Auth,
		sm.Username,
		[]string{email},
		[]byte("Subject: " + sm.Subject + "\r\n" + strconv.Itoa(code)))
	if err != nil {
		fmt.Println(err)
	}
	return
}
