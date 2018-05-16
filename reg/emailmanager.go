package reg

import (
	"net/http"
	"encoding/json"
	"errors"
	"math/rand"
	"time"
	"strconv"
	"regexp"
	"github.com/dgrijalva/jwt-go"
	"fmt"
)

const (
	//Regexp
	NickNameRegExp = "^([a-z0-9._-]){1-256}"
	EmilaRegExp    = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	//Code range
	MinCODE = 1000
	MaxCODE = 9999
)

type EmailReg struct {
	Email    string `json:"email"`
	NickName string `json:"nick"`
}

type EmailConf struct {
	Email    string `json:"email"`
	Code int `json:"code"`
}

func VailidateRegData(req *http.Request) (ei EmailReg, err error) {
	v := &EmailReg{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(v)
	if err != nil {
		err = errors.New("Invalid json")
	} else {
		re_email := regexp.MustCompile(EmilaRegExp)
		re_nodename := regexp.MustCompile(NickNameRegExp)
		is_email := re_email.MatchString(string(ei.Email))
		is_nick := re_nodename.MatchString(string(ei.NickName))
		if !is_email {
			err = errors.New("Invalid email")
		} else if !is_nick {
			err = errors.New("Invalid nickname")
		}
	}
	return
}

func ValidateEmailConfirm(req *http.Request) (ei EmailConf, err error) {
	v := &EmailConf{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(v)
	if err != nil {
		err = errors.New("Invalid json")
	} else {
		re_email := regexp.MustCompile(EmilaRegExp)
		is_email := re_email.MatchString(string(ei.Email))
		if !is_email || v.Code > MaxCODE || v.Code < MinCODE {
			err = errors.New("Invalid data")
		}
	}
	return
}

func (ei *EmailReg) CheckDBRegData(
	dbcli interface{GetKeyData(key string) (interface{}, error)}) (err error){
	nickaval, _ := dbcli.GetKeyData(ei.NickName)
	if nickaval != nil {
		emeilaval, _ := dbcli.GetKeyData(ei.Email)
		if emeilaval != nil {
			errors.New("Email already in use")
		}
	} else {
		err = errors.New("NickName already in use")
	}
	return
}

func Random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
}

func (ei EmailReg)RegisterEmail(
	dbcli interface{StoreData(string, string, int) error},
	smtp interface{SendEmail(string, int) error})(err error){
	code := Random(MinCODE, MaxCODE)
	//store code 60 seconds
	err = dbcli.StoreData(ei.Email, ei.NickName, 60)
	if err != nil {
		errors.New("DB write nick error")
	} else {
		//book nickname 60 seconds for confirmation
		err = dbcli.StoreData(ei.NickName, strconv.Itoa(code), 60)
		if err != nil {
			errors.New("DB write code error")
		} else {
			err = smtp.SendEmail(ei.Email, code)
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

func (ec *EmailConf) CheckEmailConfirm(
	dbcli interface{
		StoreData(string, string, int) error
		GetKeyData(key string) (interface{}, error)
	})(jwtoken string, err error){
	nickname, _ := dbcli.GetKeyData(ec.Email)
	if nickname == nil {
		err = errors.New("Confirmation time expired / Email not found")
	} else {
		nick := nickname.(string)
		code, _ := dbcli.GetKeyData(nick)
		if code != ec.Code {
			errors.New("Confirmation code is not match")
		} else {
			jwtoken = GenJWT(nick)
		}
	}
	return
}
