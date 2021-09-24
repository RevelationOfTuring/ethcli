package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	RpcUrl            string            `json:"rpc_url"`
	WsUrl             string            `json:"ws_url"`
	AbisPath          string            `json:"abis_path"`
	ContractAddresses map[string]string `json:"contract_addresses"`
}

func ParseConfig(configPath string) (config *Config, err error) {
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &config)
	return
}
