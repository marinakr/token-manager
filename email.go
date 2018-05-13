package main

import (
	"regexp"
	"net/smtp"
	"encoding/json"
	"strconv"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const(
	//Regexp
	NickNameRegExp = "^([a-z0-9._-]){1-256}"
	EmilaRegExp = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	//error status codes
	Ok = 0
	InvalidEmail = 1
	InvalidNickName = 2
	InvalidData = 3
	EmailInUse = 4
	NickNameInUse = 5
	CodeExpired = 6
	CodeWrong = 7
	DBError = -1
	)

type email_info struct {
	Email string
	NickName string
}

type email_confirm struct {
	Email string
	Code int
}

type smtp_data struct {
	username string
	password string
	host string
	identity string
	auth smtp.Auth
}

func (smtpdata *smtp_data)AuthEmailClient() (auth smtp.Auth) {
	auth = smtp.PlainAuth(
		smtpdata.identity,
		smtpdata.username,
		smtpdata.password,
		smtpdata.host)
	return
}

func InitEmailClient(smtpdata *smtp_data) (auth smtp.Auth) {
	redisMap := config["email_creds"]
	data, err := json.Marshal(redisMap)
	json.Unmarshal(data, smtpdata)
	if err != nil {
		smtpdata.auth = smtpdata.AuthEmailClient()
		return  smtpdata.auth
	} else {
		panic("Error smtp connection")
	}
}

func (ei *email_info) SendEmail() {
	code := Random(1000, 9999)
	smtp.SendMail(
		smtp_authinfo.host,
		smtp_authinfo.auth,
		smtp_authinfo.username,
		[]string{ei.Email},
		[]byte(strconv.Itoa(code)))
	//48 hours
	StoreRegdata(ei, code, 48*60*60)
}

func (ei *email_info) ValidateRegdata() (code int, mess string) {
	re_email := regexp.MustCompile(EmilaRegExp)
	re_nodename := regexp.MustCompile(NickNameRegExp)
	is_email := re_email.MatchString(string(ei.Email))
	is_nick := re_nodename.MatchString(string(ei.NickName))
	switch {
	case is_nick  == true, is_email == true: {
		code = Ok
	}
	case is_nick == true, is_email == false: {
		code = InvalidEmail
		mess = "Invalid email"
	}
	case is_email, is_nick == false, is_email == true: {
		code = InvalidNickName
		mess = "Invalid nickname"
	}
	default:
		code = InvalidData
		mess = "Invalid data"
	}
	return
}

func (ei *email_info)CheckRegdataAvailabe() (err int, mess string) {
	nick_code, err_nick := GetKeyData(ei.Email)
	email, err_email := GetKeyData(ei.NickName)
	if (err_nick != nil) && (err_email != nil ) {
		err = DBError
		mess = "DB error"
	} else {
		if (email == nil) && (nick_code == nil){
			err = Ok
		} else {
			if email != nil {
				err = EmailInUse
				mess = "Email already in use"
			} else {
				err = NickNameInUse
				mess = "Nick  already in use"
			}
		}
	}
	return
}

func ProcessEmail(ei email_info)(code int, mess string){
	code, mess = ei.ValidateRegdata()
	if code == Ok {
		code, mess = ei.CheckRegdataAvailabe()
		if code == Ok {
			ei.SendEmail()
		}
	}
	return
}

func GenJWT(nickname string)(jwttoken string){
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

func GenEmailJWT(ec *email_confirm)(code int, mess string){
	reg_info, err := GetKeyData(ec.Email)
	if err != nil{
		code = DBError
		mess = "DB error"
	} else {
		if reg_info == nil {
			code = CodeExpired
			mess = "Confirmation time expired"
		} else {
			ri := reg_info.(map[string]interface{})
			if ec.Code == ri["code"].(int){
				code = Ok
				mess = GenJWT(ri["nick"].(string))
			} else {
				code = CodeWrong
				mess = "Confirmation code is not match"
			}
		}
	}
	return
}