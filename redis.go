package main

import (
	"github.com/go-redis/redis"
	"time"
	"encoding/json"
)

func InitRedisClient() *redis.Client {
	redisMap := config["redis"]
	data, err_json := json.Marshal(redisMap)
	if err_json != nil{
		var result redis.Options
		json.Unmarshal(data, &result)
		client := redis.NewClient(&result)
		_, err_ping := client.Ping().Result()
		if err_ping != nil {
			panic("No connection to redis")
		} else {
			return client
		}
	} else {
		panic("Error in config: redis")
	}
}

func StoreRegdata(ei EmailInfo, expTime int) {
	JsonEmail, _ := json.Marshal(ei)
	err_em := rediscli.Set(
		ei.Email,
		JsonEmail,
		time.Duration(expTime) * time.Second).Err()
	err_nk := rediscli.Set(
		ei.NickName,
		ei.Email,
		time.Duration(expTime) * time.Second).Err()
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
