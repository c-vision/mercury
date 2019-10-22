package main

import (
	"os"

	"github.com/renproject/kv"
	"github.com/renproject/mercury/api"
	"github.com/renproject/mercury/cache"
	"github.com/renproject/mercury/proxy"
	"github.com/renproject/mercury/rpc"
	"github.com/renproject/mercury/types/btctypes"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialise logger.
	logger := logrus.StandardLogger()

	// Initialise stores.
	btcStore := kv.NewJSON(kv.NewMemDB())
	btcCache := cache.New(btcStore, logger)
	zecStore := kv.NewJSON(kv.NewMemDB())
	zecCache := cache.New(zecStore, logger)
	bchStore := kv.NewJSON(kv.NewMemDB())
	bchCache := cache.New(bchStore, logger)

	// Initialise Bitcoin API.
	btcNodeClient := rpc.NewClient(os.Getenv("BTC_RPC_URL"), "user", "password")
	btcProxy := proxy.NewProxy(btcNodeClient)
	btcAPI := api.NewApi(btctypes.BtcLocalnet, btcProxy, btcCache, logger)

	// Initialise ZCash API.
	zecNodeClient := rpc.NewClient(os.Getenv("ZEC_RPC_URL"), "user", "password")
	zecProxy := proxy.NewProxy(zecNodeClient)
	zecAPI := api.NewApi(btctypes.ZecLocalnet, zecProxy, zecCache, logger)

	// Initialise BCash API.
	bchNodeClient := rpc.NewClient(os.Getenv("BCH_RPC_URL"), "user", "password")
	bchProxy := proxy.NewProxy(bchNodeClient)
	bchAPI := api.NewApi(btctypes.BchLocalnet, bchProxy, bchCache, logger)

	// Set-up and start the server.
	server := api.NewServer(logger, "5000", btcAPI, zecAPI, bchAPI)
	server.Run()
}