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
	err_em := rediscli.Set(info.Email, map[string]interface{}{"nick": info.NickName, "code": code}, time.Duration(expTime) * time.Second).Err()
	err_nk := rediscli.Set(info.NickName, info.Email, time.Duration(expTime) * time.Second).Err()
	if (err_em != nil) && (err_nk != nil){
		panic("redis set failed")
	}
}

func GetKeyData(key string)(value interface{}, err error) {
	val, err := rediscli.Get(key).Result()
	if err == redis.Nil {
		value = nil
	} else {
		value = val
	}
	return
}
