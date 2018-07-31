package rpcclient

import (
	"github.com/ChungkueiBlock/ecoinWalletService/internal/amount"
	"github.com/ChungkueiBlock/ecoinWalletService/internal/log"
	"github.com/ChungkueiBlock/ecoinWalletService/internal/rpc"
)

type Client struct {
	endpoint string
	client   *rpc.Client
}

func NewRpc(endpoint string) (*Client, error) {

	server := &Client{
		endpoint: endpoint,
	}

	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	client, err := rpc.Dial(server.endpoint)
	if err != nil {
		return nil, err
	}
	log.Info("connected to rpc server ", "host", server.endpoint)
	server.client = client
	return server, nil
}

func (rpc *Client) Client() *rpc.Client {
	return rpc.client
}

// 似乎沒有任何作用?
func (rpc *Client) Close() {
	if rpc.client != nil {
		rpc.client.Close()
	}
}

// `{"addresses":["addr1","addr2"]}` return [sum up] all addresses balance
func (rpc *Client) GetAddressBalance(addr string) (*amount.Amount, error) {
	var results string
	err := rpc.client.Call(&results, "eth_getBalance", addr, "latest")
	log.Debug("rpc Call:eth_getBalance:", "address", addr)
	if err != nil {
		return amount.Zero, err
	}

	v, _ := amount.FromHex(results)
	return v, nil
}

func (rpc *Client) SendRawTransaction(rawtx string) (string, error) {
	var results string
	err := rpc.client.Call(&results, "eth_sendRawTransaction", rawtx)
	log.Debug("rpc Call:eth_sendRawTransaction:", "rawtx", rawtx)
	if err != nil {
		return "", err
	}

	return results, nil
}
