package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/smtp"
	"regexp"
	"strconv"
	"time"
)

const (
	//Regexp
	NickNameRegExp = "^([a-z0-9._-]){1-256}"
	EmilaRegExp    = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	//error status codes
	Ok              = 0
	InvalidEmail    = 1
	InvalidNickName = 2
	InvalidData     = 3
	EmailInUse      = 4
	NickNameInUse   = 5
	CodeExpired     = 6
	CodeWrong       = 7
	DBError         = -1
)

type EmailInfo struct {
	Email    string `json:"email"`
	NickName string `json:"nick"`
	Code     int    `json:"code"`
}

type smtp_data struct {
	username string
	password string
	host     string
	identity string
	auth     smtp.Auth
}

func (smtpdata *smtp_data) AuthEmailClient() (auth smtp.Auth) {
	auth = smtp.PlainAuth(
		smtpdata.identity,
		smtpdata.username,
		smtpdata.password,
		smtpdata.host)
	return
}

func InitEmailClient(smtpdata *smtp_data) {
	redisMap := config["email_creds"]
	data, err := json.Marshal(redisMap)
	json.Unmarshal(data, smtpdata)
	if err == nil {
		smtpdata.auth = smtpdata.AuthEmailClient()
	} else {
		fmt.Println("Error smtp connection")
		panic(err)
	}
}

func (ei *EmailInfo) SendEmail() {
	smtp.SendMail(
		smtp_authinfo.host,
		smtp_authinfo.auth,
		smtp_authinfo.username,
		[]string{ei.Email},
		[]byte(strconv.Itoa(ei.Code)))
}

func (ei *EmailInfo) ValidateRegdata() (code int, mess string) {
	re_email := regexp.MustCompile(EmilaRegExp)
	re_nodename := regexp.MustCompile(NickNameRegExp)
	is_email := re_email.MatchString(string(ei.Email))
	is_nick := re_nodename.MatchString(string(ei.NickName))
	switch {
	case is_nick == true, is_email == true:
		{
			code = Ok
		}
	case is_nick == true, is_email == false:
		{
			code = InvalidEmail
			mess = "Invalid email"
		}
	case is_email, is_nick == false, is_email == true:
		{
			code = InvalidNickName
			mess = "Invalid nickname"
		}
	default:
		code = InvalidData
		mess = "Invalid data"
	}
	return
}

func (ei *EmailInfo) CheckRegdataAvailabe() (err int, mess string) {
	data, _ := GetKeyData(ei.Email)
	email, _ := GetKeyData(ei.NickName)
	if (email == nil) && (data == nil) {
		err = Ok
	} else {
		if email != nil {
			err = NickNameInUse
			mess = "Nick already in use"
		} else {
			err = EmailInUse
			mess = "Email  already in use"
		}
	}

	return
}

func PrepareEmailCode(ei EmailInfo) (code int, mess string) {
	code, mess = ei.ValidateRegdata()
	if code == Ok {
		code, mess = ei.CheckRegdataAvailabe()
		if code == Ok {
			ei.Code = Random(1000, 9999)
			//48 hours
			StoreRegdata(ei, 48*60*60)
			ei.SendEmail()
		}
	}
	return
}

func GenJWT(nickname string) (jwttoken string) {
	//TODO: read from config
	mySigningKey := "secretkey"
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = nickname
	claims["exp"] = time.Now().Add(time.Hour * 24 * 31 * 12).Unix()
	claims["iot"] = time.Now().Unix()
	jwttoken, _ = token.SignedString(mySigningKey)
	return
}

func ConfirmEmail(ec *EmailInfo) (code int, mess string) {
	data, err := GetKeyData(ec.Email)
	if err != nil {
		code = DBError
		mess = "DB error"
	} else {
		if data == nil {
			code = CodeExpired
			mess = "Confirmation time expired"
		} else {
			regdata := &EmailInfo{}
			json.Unmarshal([]byte(data.(string)), regdata)
			if ec.Code == regdata.Code {
				StoreRegdata(*regdata, 0)
				code = Ok
				mess = GenJWT(regdata.NickName)
			} else {
				code = CodeWrong
				mess = "Confirmation code is not match"
			}
		}
	}
	return
}
