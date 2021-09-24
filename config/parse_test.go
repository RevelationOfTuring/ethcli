package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	config, err := ParseConfig("./config_test.json")
	require.NoError(t, err)
	require.Equal(t, "https://exchaintestrpc.okex.org", config.RpcUrl)
	require.Equal(t, "wss://exchaintestws.okex.org:8443", config.WsUrl)
	require.Equal(t, "./abis_test", config.AbisPath)
	require.Equal(t, "./.priv_key", config.PrivKeyPath)
	require.Equal(t, int64(1000000000), config.GasPrice)
	require.Equal(t, 2, len(config.ContractAddresses))
	require.Equal(t, "0x0000000000000000000000000000000000000000", config.ContractAddresses["contract1"])
	require.Equal(t, "0x0000000000000000000000000000000000000001", config.ContractAddresses["contract2"])
}
