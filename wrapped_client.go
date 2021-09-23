package ethcli

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type WrappedClient struct {
	*ethclient.Client
	chainID *big.Int
}

//func NewEthClient(configPath string) *WrappedClient {
//
//}
