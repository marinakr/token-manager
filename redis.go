package main

import (
	"github.com/go-redis/redis"
	"time"
)

func InitRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil{
		panic("No connection to redis")
	} else {
		return client
	}
}

func store_regdata(client *redis.Client, info email_info, expTime int) {
	err_em := client.Set(info.Email, info.NickName, time.Duration(expTime) * time.Second).Err()
	err_nk := client.Set(info.NickName, info.Email, time.Duration(expTime) * time.Second).Err()

	if (err_em != nil) && (err_nk != nil) {
		panic("redis set failed")
	}
}

func get_keydata(client *redis.Client, key string)(nickname interface{}, err error) {
	val, err := client.Get(key).Result()
	if err == redis.Nil {
		nickname = nil
	} else if err != nil {
		panic(err)
	} else {
		nickname = val
	}
	return
}
