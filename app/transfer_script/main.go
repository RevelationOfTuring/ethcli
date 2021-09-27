package main

import (
	"log"

	"github.com/RevelationOfTuring/ethcli"
)

func main() {
	cli, err := ethcli.NewEthClient("./sample/config_oec_testnet.json")
	if err != nil {
		log.Fatalln(err)
	}




}
