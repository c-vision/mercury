package ethclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/renproject/mercury/types"
)

// EthClient is a client which is used to talking with certain bitcoin network. It can interacting with the blockchain
// through Mercury server.
type EthClient struct {
	url string
}

// NewEthClient returns a new EthClient of given bitcoin network.
func NewEthClient(network types.EthNetwork) *EthClient {
	switch network {
	case types.EthMainnet:
		return &EthClient{
			url: "https://ren-mercury.herokuapp.com/eth",
		}
	case types.EthKovan:
		return &EthClient{
			url: "https://ren-mercury.herokuapp.com/eth-kovan",
		}
	default:
		panic("unknown eth network")
	}
}

// Balance returns the balance of the given bitcoin address. It filters the utxos which have less confirmations than
// required. It times out if the context exceeded.
func (client *EthClient) Balance(ctx context.Context, address types.EthAddr) (types.Amount, error) {
	b := []string{address.Hex(), "latest"}
	data, err := json.Marshal(b)
	if err != nil {
		return types.Amount{}, err
	}
	request := types.JSONRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "eth_getBalance",
		Params:  data,
	}
	var response types.JSONResponse
	resp, err := client.sendRequest(request)
	if err != nil {
		return types.Amount{}, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return types.Amount{}, err
	}
	var res string
	if err := json.Unmarshal(response.Result, &res); err != nil {
		return types.Amount{}, err
	}
	value, err := hexutil.DecodeBig(res)
	if err != nil {
		return types.Amount{}, err
	}
	return types.NewAmount(value), nil
}

// sendRequest sends the JSON-2.0 request to the target url and returns the response and any error.
func (client *EthClient) sendRequest(request types.JSONRequest) (*http.Response, error) {
	var url string
	if !strings.HasPrefix(client.url, "http") {
		url = "http://" + client.url
	} else {
		url = client.url
	}

	httpclient := &http.Client{
		Timeout: 10 * time.Second,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(data)

	return httpclient.Post(url, "application/json", buff)
}
