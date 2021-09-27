package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/RevelationOfTuring/ethcli"
	ethcmn "github.com/ethereum/go-ethereum/common"
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

	txHash, err := cli.CallContract(
		0,
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

	//receipt, err := cli.TransactionReceipt(context.Background(), ethcmn.HexToHash("0x230022E4820DD1F3896BD90D7DD1147433C097AB79EDAE13E3FBA6B51D9F8EF8"))
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println(receipt.Status)
}
