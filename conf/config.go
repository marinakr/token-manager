package conf

import (
	"os"
	"encoding/json"
)

type conf map[string]interface{}

func ReadConfigFile(configuration *map[string]interface{}){
	cfgfile, _ := os.Open("config.json")
	defer cfgfile.Close()
	decoder := json.NewDecoder(cfgfile)
	err := decoder.Decode(configuration)
	if err != nil {
		panic("Can not read config file!")
	}
}

func NewConig() conf {
	config := make(map[string]interface{})
	ReadConfigFile(&config)
	return  config
}

func (c *conf)RedisConf() interface{}{
	return (*c)["redis"]
}

func (c *conf)EmailConf() interface{}{
	return (*c)["email_creds"]
}

func (c *conf)PortConf() string{
	tm := (*c)["token_manager"].(map[string]interface{})
	return  tm["port"].(string)
}