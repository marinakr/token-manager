package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)



func StoreRegdata(ei EmailInfo, expTime int) {
	JsonEmail, _ := json.Marshal(ei)
	err_em := rediscli.Set(
		ei.Email,
		JsonEmail,
		time.Duration(expTime)*time.Second).Err()
	err_nk := rediscli.Set(
		ei.NickName,
		ei.Email,
		time.Duration(expTime)*time.Second).Err()
	if (err_em != nil) && (err_nk != nil) {
		panic("redis set failed")
	}
}

func GetKeyData(key string) (value interface{}, err error) {
	val, err := rediscli.Get(key).Result()
	if err == redis.Nil {
		value = nil
	} else {
		value = val
	}
	return
}
