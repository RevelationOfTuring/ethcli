package main

import (
	"ethcli"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
)

func main() {
	cli, err := ethcli.NewEthClient("./sample/config_oec_testnet.json")
	if err != nil {
		log.Fatalln(err)
	}

	amount, ok := new(big.Int).SetString("1000000000000000000000000000000000", 10)
	if !ok {
		panic("convert failed")
	}

	fmt.Println(amount)
	txHash, err := cli.CallContract(
		"ERC20Token",
		big.NewInt(0),
		"mintDirectly",
		ethcmn.HexToAddress("0xF202E4e0EB3C10c4F0ace15aF6B6EA3AFAe777AC"),
		amount,
	)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(txHash)
}
