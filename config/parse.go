package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	RpcUrl            string            `json:"rpc_url"`
	WsUrl             string            `json:"ws_url"`
	AbisPath          string            `json:"abis_path"`
	PrivKeyPath       string            `json:"priv_key_path"`
	GasPrice          int64             `json:"gas_price"`
	IsOECKind         bool              `json:"is_oec_kind"`
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
