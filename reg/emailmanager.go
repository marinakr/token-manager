package reg

import (
	//
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
	//
)

const (
	//Regexp
	NickNameRegExp = "^[a-zA-Z0-9._-]{1,256}$"
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
	Email string `json:"email"`
	Code  int    `json:"code"`
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
		is_email := re_email.MatchString(v.Email)
		is_nick := re_nodename.MatchString(v.NickName)
		if !is_email {
			err = errors.New("Invalid email")
		} else if !is_nick {
			err = errors.New("Invalid nickname")
		} else {
			ei = *v
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
		is_email := re_email.MatchString(string(v.Email))
		if !is_email || v.Code > MaxCODE || v.Code < MinCODE {
			err = errors.New("Invalid data")
		} else {
			ei = *v
		}
	}
	return
}

func (ei *EmailReg) CheckDBRegData(
	dbcli interface {
		GetKeyData(key string) (interface{}, error)
	}) (err error) {
	emeilaval, _ := dbcli.GetKeyData(ei.NickName)
	nickaval, _ := dbcli.GetKeyData(ei.Email)
	if nickaval != nil {
		err = errors.New("NickName already in use")
	}
	if emeilaval != nil {
		err = errors.New("Email already in use")

	}
	return
}

func Random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func (ei EmailReg) RegisterEmail(
	dbcli interface {
		StoreData(string, interface{}, int) error
	},
	smtp interface {
		SendEmail(string, int) error
	}) (err error) {
	code := Random(MinCODE, MaxCODE)
	//store code 60 seconds
	err = dbcli.StoreData(ei.Email, ei.NickName, 60)
	if err != nil {
		errors.New("DB write nick error")
	} else {
		//book nickname 60 seconds for confirmation
		err = dbcli.StoreData(ei.NickName, code, 60)
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
	dbcli interface {
		StoreData(string, interface{}, int) error
		GetKeyData(key string) (interface{}, error)
	}) (jwtoken string, err error) {
	nickname, _ := dbcli.GetKeyData(ec.Email)
	if nickname == nil {
		err = errors.New("Confirmation time expired / Email not found")
	} else {
		nick := nickname.(string)
		codestr, _ := dbcli.GetKeyData(nick)
		code, _ := strconv.Atoi(codestr.(string))
		if code != ec.Code {
			err = errors.New("Confirmation code is not match")
		} else {
			jwtoken = GenJWT(nick)
		}
	}
	return
}
