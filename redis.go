package main

import (
	"github.com/go-redis/redis"
	"time"
	"encoding/json"
)

func InitRedisClient() *redis.Client {
	redisMap := config["redis"]
	data, err_json := json.Marshal(redisMap)
	var result redis.Options
	json.Unmarshal(data, &result)
	client := redis.NewClient(&result)
	_, err := client.Ping().Result()
	if err != nil || err_json != nil {
		panic("No connection to redis")
	} else {
		return client
	}
}

func StoreRegdata(info *email_info, code int, expTime int) {
	err_em := rediscli.Set(info.Email, info.NickName, time.Duration(expTime) * time.Second).Err()
	err_nk := rediscli.Set(info.NickName, info.Email, time.Duration(expTime) * time.Second).Err()
	err_code := rediscli.Set(info.Email, code, time.Duration(expTime) * time.Second).Err()
	if (err_em != nil) && (err_nk != nil) && (err_code != nil){
		panic("redis set failed")
	}
}

func GetKeyData(key string)(nickname interface{}, err error) {
	val, err := rediscli.Get(key).Result()
	if err == redis.Nil {
		nickname = nil
	} else {
		nickname = val
	}
	return
}
