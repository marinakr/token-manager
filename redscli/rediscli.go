package redscli

import (
	//
	"github.com/go-redis/redis"
	"encoding/json"
	"time"
	//
)

type RedisENV interface {
	GetKeyData(key string) (value interface{}, err error)
	StoreData(key string, value interface{}, exp int) error
}

type DBCli struct {
	*redis.Client
}

func New(redisCreds interface{}) *DBCli {
	data, err_json := json.Marshal(redisCreds)
	if err_json == nil {
		var result redis.Options
		json.Unmarshal(data, &result)
		client := redis.NewClient(&result)
		_, err_ping := client.Ping().Result()
		if err_ping != nil {
			panic("No connection to redis")
		} else {
			return &DBCli{client}
		}
	} else {
		panic("Error in config file: redis")
	}
}

func (cli *DBCli) GetKeyData(key string) (value interface{}, err error) {
	val, err := cli.Get(key).Result()
	if err == redis.Nil {
		value = nil
	} else {
		value = val
	}
	return
}

func (cli *DBCli) StoreData(key string, value interface{}, exp int) error {
	return cli.Set(key, value, time.Duration(exp)*time.Second).Err()
}

