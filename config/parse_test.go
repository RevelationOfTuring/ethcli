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
}
