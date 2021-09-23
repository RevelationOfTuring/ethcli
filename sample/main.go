package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	rpcUrlTestnet = "https://exchaintestrpc.okex.org"
)

func main() {
	rpcClient, err := rpc.DialContext(context.Background(), rpcUrlTestnet)
	if err != nil {
		panic(err)
	}

	ethcli := ethclient.NewClient(rpcClient)
	chainId, err := ethcli.ChainID(context.Background())
	if err!=nil{
		panic(err)
	}

	fmt.Println(chainId)
}
