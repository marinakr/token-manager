package main

import (
	"regexp"
)

const(
	//Regexp
	NickNameRegExp = "^([a-z0-9\\.\\-_\\+]+){1-256}"
	EmilaRegExp = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	//error status codes
	Ok = 0
	InvalidEmail = 1
	InvalidNickName = 2
	InvalidData = 3
	EmailInUse = 4
	NickNameInUse = 5
	)

type email_info struct {
	Email string
	NickName string
}

func validate_regdata(email email_info) (code int, mess string) {
	re_email := regexp.MustCompile(EmilaRegExp)
	re_nodename := regexp.MustCompile(NickNameRegExp)
	is_email := re_email.MatchString(string(email.Email))
	is_nick := re_nodename.MatchString(string(email.NickName))
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

func check_regdata_availabe() int {

}