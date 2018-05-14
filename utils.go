package main

import (
		"net/http"
		"encoding/json"
		"math/rand"
		"time"
		"os"
		"fmt"
       )

type resp_payload struct {
	code int
		mess string
}

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

func EncodeReqResp(rw http.ResponseWriter, Status int, mess string){
	fmt.Fprintf(rw, mess)
		rw.WriteHeader(http.StatusOK)
}

func Random(min, max int) int {
	rand.Seed(time.Now().Unix())
		return rand.Intn(max - min) + min
}

func ReadConfig(configuration *map[string]interface{}){
	cfgfile, _ := os.Open("config.json")
		defer cfgfile.Close()
		decoder := json.NewDecoder(cfgfile)
		err := decoder.Decode(configuration)
		if err != nil {
			panic("Can not read config file!")
		}
}

func GenResponsePayload(code int, mess string)(string){
payload := make(map[string]interface{})
		 payload["code"] = code
		 payload["mess"] = mess
		 resp, _ := json.Marshal(payload)
		 resp_payload := string(resp[:])
		 return resp_payload
}
