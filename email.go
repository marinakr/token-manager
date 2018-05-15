package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/smtp"
	"regexp"
	"strconv"
	"time"
	"jwt-go"
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
)

type EmailInfo struct {
	Email    string `json:"email"`
	NickName string `json:"nick"`
	Code     int    `json:"code"`
}

type SmtpData struct {
	Username string
	Password string
	Host     string
	Smtphost string
	Identity string
	Auth     smtp.Auth
}

func InitEmailClient(smtpdata *SmtpData) {
	emailCreds := config["email_creds"]
	data, err := json.Marshal(emailCreds)
	json.Unmarshal(data, smtpdata)
	if err == nil {
		smtpdata.Auth = smtp.PlainAuth(smtpdata.Identity, smtpdata.Username, smtpdata.Password, smtpdata.Smtphost)
	} else {
		fmt.Println("Error smtp auth")
		panic(err)
	}
}

func (smtpdata *SmtpData) SendEmail(ei EmailInfo) {
	err := smtp.SendMail(
		smtpdata.Host,
		smtpdata.Auth,
		smtpdata.Username,
		[]string{ei.Email},
		[]byte(strconv.Itoa(ei.Code)))
	if err != nil {
		fmt.Println(err)
	}
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
			//48 hours  48*60*60
			StoreRegdata(ei, 60)
			smtpdata.SendEmail(ei)
		}
	}
	return
}

func GenJWT(nickname string) string {
	//TODO: read from config
	mySigningKey := []byte("secretkey")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = nickname
	claims["exp"] = time.Now().Add(time.Hour * 24 * 31 * 12).Unix()
	claims["iot"] = time.Now().Unix()
	jwtres, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("JWT generation error: ", err)
		panic(err)
	}
	return jwtres
}

func ConfirmEmail(ec *EmailInfo) (code int, mess string) {
	data, _ := GetKeyData(ec.Email)
	if data == nil {
		code = CodeExpired
		mess = "Confirmation time expired / Email not found"
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

	return
}
