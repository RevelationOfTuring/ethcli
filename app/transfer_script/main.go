package main

import (
	"encoding/json"
	"github.com/RevelationOfTuring/ethcli"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"log"
	"math/big"
	"time"
)

type Config struct {
	To    string `json:"to"`
	Value string `json:"value"`
}

func main() {
	bytes, err := ioutil.ReadFile("./transfer.json")
	if err != nil {
		log.Fatalln(err)
	}

	var transferCfg Config
	if err = json.Unmarshal(bytes, &transferCfg); err != nil {
		log.Fatalln(err)
	}

	transferAmount, ok := new(big.Int).SetString(transferCfg.Value, 10)
	if !ok {
		panic("convert transfer amount failed")
	}

	toAddr := ethcmn.HexToAddress(transferCfg.To)

	cli, err := ethcli.NewEthClient("./config_network.json")
	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < cli.GetEcdsaKeysNum(); i++ {

		fromAddr, err := cli.GetAddress(i)
		if err != nil {
			log.Fatalln(err)
		}

		txHash, err := cli.SendTx(i, toAddr, transferAmount, 21000, nil)
		if err != nil {
			log.Printf("key index [%d] Address [%s]: send Tx error: %s \n", i, fromAddr, err)
		}

		log.Printf("[%d] Address [%s] transfers successfully, tx hash [%s]\n", i+1, fromAddr, txHash)

		time.Sleep(time.Second)
	}

}
