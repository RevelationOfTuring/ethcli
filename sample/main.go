package main

import (
	"ethcli"
	"fmt"
	"log"
)

func main() {
	cli, err := ethcli.NewEthClient("./sample/config_oec_testnet.json")
	if err != nil {
		log.Fatalln(err)
	}

	//err = cli.CallContract("IERC20", "balanceOf", ethcmn.HexToAddress("0xF202E4e0EB3C10c4F0ace15aF6B6EA3AFAe777AC"))
	//if err != nil {
	//	log.Fatalln(err)
	//}
	fmt.Println(cli.GetChainId())
}
