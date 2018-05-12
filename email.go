package main

import (
	"regexp"
)

const(
	NickNameRegExp = "^([a-z0-9\\.\\-_\\+]+){1-256}"
	EmilaRegExp = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

type email_info struct {
	Email string
	NickName string
}

func process_email(email email_info) error {
	re_email := regexp.MustCompile(EmilaRegExp)
	re_nodename := regexp.MustCompile(NickNameRegExp)
	is_email := re_email.MatchString(string(email.Email))
	is_phone := re_nodename.MatchString(string(email.NickName))

	return
}
