package ethcli

import (
	"context"
	"crypto/ecdsa"
	"ethcli/config"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type WrappedClient struct {
	*ethclient.Client
	chainId           *big.Int
	abis              map[string]abi.ABI
	contractAddresses map[string]ethcmn.Address
	ecdsakey          *ecdsa.PrivateKey
}

func (wc *WrappedClient) GetChainId() *big.Int {
	return wc.chainId
}

// ONLY invoke at the initialization of WrappedClient
func (wc *WrappedClient) loadContractInfos(cfg *config.Config) error {
	wc.abis = make(map[string]abi.ABI)
	wc.contractAddresses = make(map[string]ethcmn.Address)
	for contractName, contracrAddr := range cfg.ContractAddresses {
		wc.contractAddresses[contractName] = ethcmn.HexToAddress(contracrAddr)
	}

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
	if err != nil {
		return
	}

	wrappedCli.ecdsakey, err = crypto.LoadECDSA(cfg.PrivKeyPath)
	if err != nil {
		return
	}

	rpcClient, err := rpc.DialContext(context.Background(), cfg.RpcUrl)
	if err != nil {
		return
	}
	wrappedCli.Client = ethclient.NewClient(rpcClient)

	wrappedCli.chainId, err = wrappedCli.Client.ChainID(context.Background())
	return
}
