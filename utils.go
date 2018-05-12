package main

import (
	"net/http"
	"encoding/json"
)

func DecodeReqBody(req *http.Request, v interface{}) (code int, message string){
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(v)
	if err != nil{
		code = http.StatusBadRequest
		message = "Invalid json"
	} else {
		code = 0
	}
	return
}

