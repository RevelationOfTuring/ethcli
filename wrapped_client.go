package ethcli

import (
	"context"
	"ethcli/config"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

type WrappedClient struct {
	*ethclient.Client
	chainId           *big.Int
	abis              map[string]abi.ABI
	contractAddresses map[string]string
}

func (wc *WrappedClient) GetChainId() *big.Int {
	return wc.chainId
}

// ONLY invoke at the initialization of WrappedClient
func (wc *WrappedClient) loadContractInfos(cfg *config.Config) error {
	wc.abis = make(map[string]abi.ABI)
	wc.contractAddresses = cfg.ContractAddresses

	filesInfo, err := ioutil.ReadDir(cfg.AbisPath)
	if err != nil {
		return err
	}

	const abiSuffix = ".json"
	for _, fileInfo := range filesInfo {
		fileName := fileInfo.Name()
		if !strings.HasSuffix(fileName, abiSuffix) {
			continue
		}

		// read abi files
		f, err := os.Open(filepath.Join(cfg.AbisPath, fileName))
		if err != nil {
			return err
		}

		a, err := abi.JSON(f)
		if err != nil {
			return err
		}

		wc.abis[strings.TrimSuffix(fileName, abiSuffix)] = a
	}

	return nil
}

func (wc *WrappedClient) buildInput(contractName, methodName string, args ...interface{}) ([]byte, error) {
	a, ok := wc.abis[contractName]
	if !ok {
		return nil, fmt.Errorf("abi of contract %s is missed", contractName)
	}

	return a.Pack(methodName, args)
}

func NewEthClient(configPath string) (wrappedCli *WrappedClient, err error) {
	cfg, err := config.ParseConfig(configPath)
	if err != nil {
		return
	}

	wrappedCli = new(WrappedClient)
	err = wrappedCli.loadContractInfos(cfg)

	rpcClient, err := rpc.DialContext(context.Background(), cfg.RpcUrl)
	if err != nil {
		return
	}
	wrappedCli.Client = ethclient.NewClient(rpcClient)

	wrappedCli.chainId, err = wrappedCli.Client.ChainID(context.Background())
	return
}
