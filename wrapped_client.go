package ethcli

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/RevelationOfTuring/ethcli/config"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/okex/exchain-ethereum-compatible/utils"
)

type WrappedClient struct {
	*ethclient.Client
	chainId           *big.Int
	abis              map[string]abi.ABI
	contractAddresses map[string]ethcmn.Address
	ecdsaKeys         []*ecdsa.PrivateKey
	addresses         []ethcmn.Address
	gasPrice          *big.Int
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

func (wc *WrappedClient) checkKeyIndex(keyIndex int) error {
	if keyIndex < len(wc.ecdsaKeys) {
		return nil
	}

	return errors.New("key index is out of range to the length of WrappedClient.ecdsaKeys")
}

func (wc *WrappedClient) getNonce(keyIndex int) (nonce uint64, err error) {
	if err = wc.checkKeyIndex(keyIndex); err != nil {
		return
	}

	for i := 0; i < 5; i++ {
		// query again with 5 times in case of timeout
		nonce, err = wc.PendingNonceAt(context.Background(), wc.addresses[keyIndex])
		if err != nil {
			continue
		}
		return // successfully
	}

	return nonce, fmt.Errorf("fail to get nonce of %s", wc.addresses[keyIndex])
}

func (wc *WrappedClient) CallContract(keyIndex int, contractName string, value *big.Int, methodName string, args ...interface{}) (
	txHash ethcmn.Hash, err error) {
	contractAddr, ok := wc.contractAddresses[contractName]
	if !ok {
		return txHash, fmt.Errorf("contract address of %s isn't provided in config file", contractName)
	}

	a, ok := wc.abis[contractName]
	if !ok {
		return txHash, fmt.Errorf("abi of contract %s is missed", contractName)
	}

	input, err := a.Pack(methodName, args...)
	if err != nil {
		return
	}

	gasEstimate, err := wc.EstimateGas(context.Background(), ethereum.CallMsg{
		From:     wc.addresses[keyIndex],
		To:       &contractAddr,
		GasPrice: wc.gasPrice,
		Data:     input,
	})
	if err != nil {
		gasEstimate = uint64(500000)
	}

	return wc.SendTx(keyIndex, contractAddr, value, gasEstimate, input)
}

func (wc *WrappedClient) SendTx(keyIndex int, to ethcmn.Address, value *big.Int, gasLimit uint64, input []byte) (
	txHash ethcmn.Hash, err error) {
	nonce, err := wc.getNonce(keyIndex)
	if err != nil {
		return
	}

	unsignedTx := types.NewTransaction(nonce, to, value, gasLimit, wc.gasPrice, input)
	signedTx, err := types.SignTx(unsignedTx, types.NewLondonSigner(wc.chainId), wc.ecdsaKeys[keyIndex])
	if err != nil {
		return
	}

	err = wc.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return
	}

	return utils.Hash(signedTx)
}

func (wc *WrappedClient) LoadPrivKeysFromFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	buffer := bufio.NewReader(f)
	log.Println("loading private keys ...")

	count := 1
	for {
		bytes, _, err := buffer.ReadLine()
		if err == io.EOF {
			break
		}

		privKey, err := crypto.HexToECDSA(string(bytes))
		if err != nil {
			return err
		}

		address := crypto.PubkeyToAddress(*(privKey.Public()).(*ecdsa.PublicKey))
		wc.ecdsaKeys = append(wc.ecdsaKeys, privKey)
		wc.addresses = append(wc.addresses, address)
		fmt.Printf("		    %d: %s\n", count, address)
		count++
	}

	return nil
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

	err = wrappedCli.LoadPrivKeysFromFile(cfg.PrivKeyPath)
	if err != nil {
		return
	}

	wrappedCli.gasPrice = big.NewInt(cfg.GasPrice)
	log.Printf("gas price is %s\n", wrappedCli.gasPrice)

	rpcClient, err := rpc.DialContext(context.Background(), cfg.RpcUrl)
	if err != nil {
		return
	}
	wrappedCli.Client = ethclient.NewClient(rpcClient)

	wrappedCli.chainId, err = wrappedCli.Client.ChainID(context.Background())
	return
}
