package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	RpcUrl string `json:"rpc_url"`
	WsUrl  string `json:"ws_url"`
}

func ParseConfig(configPath string) (config *Config, err error) {
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return
	}

	if err = json.Unmarshal(bytes, &config); err != nil {
		return
	}

	return
}
