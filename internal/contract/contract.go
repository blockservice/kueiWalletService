package contract

import (
	"errors"
	"fmt"
	"github.com/ChungkueiBlock/kueiWalletService/internal/amount"
	"github.com/ChungkueiBlock/kueiWalletService/internal/log"
	"github.com/ChungkueiBlock/kueiWalletService/internal/rpc"
	EthAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"strings"
)

const abijson = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"version","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"},{"name":"_extraData","type":"bytes"}],"name":"approveAndCall","outputs":[{"name":"success","type":"bool"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"remaining","type":"uint256"}],"payable":false,"type":"function"},{"inputs":[{"name":"_initialAmount","type":"uint256"},{"name":"_tokenName","type":"string"},{"name":"_decimalUnits","type":"uint8"},{"name":"_tokenSymbol","type":"string"}],"type":"constructor"},{"payable":false,"type":"fallback"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_from","type":"address"},{"indexed":true,"name":"_to","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_owner","type":"address"},{"indexed":true,"name":"_spender","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Approval","type":"event"}]`
const (
	MethodIdSymbol      = "symbol"
	MethodIdName        = "name"
	MethodIdDecimals    = "decimals"
	MethodIdBalanceOf   = "balanceOf"
	MethodIdTotalSupply = "totalSupply"
)

type CoinwallEthContract struct {
	address   string
	abi       *CoinwallEthAbi
	rpcClient *rpc.Client
}
type EthCallTx struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Value    string `json:"value"`
	Data     string `json:"data"`
}

func NewContractFromAddress(addr string) (*CoinwallEthContract, error) {
	//address := common.HexToAddress(addr)

	contract := &CoinwallEthContract{}
	contract.address = addr

	err := contract.initiate()
	if err != nil {
		return nil, err
	}

	return contract, nil
}

func (contract *CoinwallEthContract) Symbol() (string, error) {
	methodId := MethodIdSymbol
	b, err := contract.tokenInfo(methodId)
	if err != nil {
		return "", nil
	}
	var v string
	err = contract.abi.Unpack(&v, methodId, b)
	if err != nil {
		return "", err
	}
	return v, nil
}
func (contract *CoinwallEthContract) Name() (string, error) {
	methodId := MethodIdName
	b, err := contract.tokenInfo(methodId)
	if err != nil {
		return "", nil
	}
	var v string
	err = contract.abi.Unpack(&v, methodId, b)
	if err != nil {
		return "", err
	}
	return v, nil
}
func (contract *CoinwallEthContract) Decimals() (int, error) {
	methodId := MethodIdDecimals
	b, err := contract.tokenInfo(methodId)
	if err != nil {
		return 0, nil
	}
	var v uint8
	err = contract.abi.Unpack(&v, methodId, b)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}
func (contract *CoinwallEthContract) TotalSupply() (*amount.Amount, error) {
	methodId := MethodIdTotalSupply
	b, err := contract.tokenInfo(methodId)
	if err != nil {
		return amount.Zero, nil
	}
	var v *big.Int
	err = contract.abi.Unpack(&v, methodId, b)
	if err != nil {
		return amount.Zero, err
	}
	value := amount.NewAmount(v)
	return value, nil
}

func (contract *CoinwallEthContract) initiate() error {
	if contract.address == "" {
		return errors.New("contract address must specified")
	}

	if contract.abi == nil {
		contract.ABI()
	}

	if contract.abi == nil {
		return errors.New("contract abi must instantiated")
	}

	return nil
}

func (contract *CoinwallEthContract) SetRpc(client *rpc.Client) {
	contract.rpcClient = client
}

func (contract *CoinwallEthContract) tokenInfo(infoType string) ([]byte, error) {
	var methodId string
	switch infoType {
	case MethodIdSymbol:
		methodId = MethodIdSymbol
	case MethodIdName:
		methodId = MethodIdName
	case MethodIdDecimals:
		methodId = MethodIdDecimals
	case MethodIdTotalSupply:
		methodId = MethodIdTotalSupply
	default:
		return nil, errors.New("querying unsupported token info")
	}

	err := contract.initiate()
	if err != nil {
		return nil, err
	}

	dataBytes, err := contract.abi.Pack(methodId)
	if err != nil {
		return nil, err
	}
	log.Trace("method id:", hexutil.Encode(dataBytes))
	if err != nil {
		return nil, err
	}
	tx := &EthCallTx{
		From: contract.address,
		To:   contract.address,
		Data: hexutil.Encode(dataBytes),
	}

	wsRpcClient := contract.rpcClient
	if err != nil {
		return nil, err
	}

	var results string
	log.Trace("tx:", *tx)
	err = wsRpcClient.Call(&results, "eth_call", tx, "latest")
	if err != nil {
		log.Error("eth_call failed:", err.Error())
	}
	log.Trace("rpc Call:eth_call", results, *tx)

	b, err := hexutil.Decode(results)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (contract *CoinwallEthContract) GetAddressBalance(addr string) (*amount.Amount, error) {
	err := contract.initiate()
	if err != nil {
		return amount.Zero, err
	}

	methodId := MethodIdBalanceOf
	dataBytes, err := contract.abi.Pack(methodId, common.HexToAddress(addr))
	if err != nil {
		return nil, err
	}
	log.Trace("method id:", hexutil.Encode(dataBytes))
	if err != nil {
		return amount.Zero, err
	}
	tx := &EthCallTx{
		From: addr,
		To:   contract.address,
		Data: hexutil.Encode(dataBytes),
	}

	wsRpcClient := contract.rpcClient
	if err != nil {
		return amount.Zero, err
	}

	var results string
	log.Trace("tx:", *tx)
	err = wsRpcClient.Call(&results, "eth_call", tx, "latest")
	if err != nil {
		log.Error("eth_call failed:", err.Error())
	}
	log.Trace("rpc Call:eth_call", results, *tx)

	b, err := hexutil.Decode(results)
	if err != nil {
		return amount.Zero, err
	}
	if len(b) == 0 {
		return amount.Zero, nil
	}

	var v *big.Int
	err = contract.abi.Unpack(&v, methodId, b)
	if err != nil {
		return amount.Zero, err
	}
	c := amount.NewAmount(v)
	return c, nil
}

func (contract *CoinwallEthContract) ABI() *CoinwallEthAbi {
	ethAbi, err := EthAbi.JSON(strings.NewReader(abijson))
	if err != nil {
		panic(err)
	}
	abi := CoinwallEthAbi{ethAbi}

	contract.abi = &abi
	return contract.abi
}

// https://theethereum.wiki/w/index.php/ERC20_Token_Standard
// github.com/ethereum/go-ethereum/accounts/abi#ABI
type CoinwallEthAbi struct {
	EthAbi.ABI
}

func (abi *CoinwallEthAbi) UnpackInputs(v interface{}, name string, output []byte) (err error) {
	if len(output) == 0 {
		return fmt.Errorf("abi: unmarshalling empty output")
	}
	// since there can't be naming collisions with contracts and events,
	// we need to decide whether we're calling a method or an event
	if method, ok := abi.Methods[name]; ok {
		if len(output)%32 != 0 {
			return fmt.Errorf("abi: improperly formatted output")
		}
		return method.Inputs.Unpack(v, output)
	} else if event, ok := abi.Events[name]; ok {
		return event.Inputs.Unpack(v, output)
	}
	return fmt.Errorf("abi: could not locate named method or event")
}
