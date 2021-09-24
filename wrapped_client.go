package ethcli

import (
	"context"
	"ethcli/config"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

type WrappedClient struct {
	*ethclient.Client
	chainId *big.Int
}

func (wc *WrappedClient) GetChainId() *big.Int {
	return wc.chainId
}

func NewEthClient(configPath string) (wrappedCli *WrappedClient, err error) {
	cfg, err := config.ParseConfig(configPath)
	if err != nil {
		return
	}

	wrappedCli = &WrappedClient{}
	rpcClient, err := rpc.DialContext(context.Background(), cfg.RpcUrl)
	if err != nil {
		return
	}
	wrappedCli.Client = ethclient.NewClient(rpcClient)

	wrappedCli.chainId, err = wrappedCli.Client.ChainID(context.Background())
	return
}
