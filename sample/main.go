package main

import (
	"ethcli"
	"fmt"
	"log"
)

func main() {
	cli, err := ethcli.NewEthClient("./sample/config_oec_testnet.json")
	if err!=nil{
		log.Fatalln(err)
	}

	fmt.Println(cli.GetChainId())
}
