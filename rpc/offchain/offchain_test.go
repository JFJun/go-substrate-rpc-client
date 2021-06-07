package offchain

import (
	"os"
	"testing"

	"github.com/JFJun/go-substrate-rpc-client/client"
	"github.com/JFJun/go-substrate-rpc-client/config"
)

var offchain *Offchain

func TestMain(m *testing.M) {
	cl, err := client.Connect(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}
	offchain = NewOffchain(cl)
	os.Exit(m.Run())
}
