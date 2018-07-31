package ews

import (
	"fmt"
	"github.com/ChungkueiBlock/ecoinWalletService/internal/log"
	"github.com/ChungkueiBlock/ecoinWalletService/internal/storage/mysql"
)

type WalletAPI struct {
	cfg     *Config
	rpchost string
}

var db *mysql.Connection

func NewWalletAPI(ews *ChungkueiWalletService) *WalletAPI {
	db = mysql.NewConnection(ews.config.MysqlDSN)

	if ews.config.EnableGasStation {
		log.Info("gasstation enabled")
		go updatePredictTableJson()
	}

	return &WalletAPI{
		cfg: ews.config,
	}
}

func (api *WalletAPI) Foo(p1 string) string {
	return p1 + " response successfully!"
}

func (api *WalletAPI) parsePage(page int) (pageFrom, pageTo int) {
	if page <= 0 {
		page = 0
	}
	pageSize := api.cfg.TxPageSize
	pageFrom = page * pageSize
	pageTo = (page + 1) * pageSize
	return
}
func (api *WalletAPI) TransactionHistory(address string, page int, withReceipt int) []*mysql.Transaction {
	pageFrom, pageTo := api.parsePage(page)
	b, err := db.GetTransactionHistory(address, pageFrom, pageTo)
	if err != nil {
		log.Error("TransactionHistory fails.", "address", address, "err", err)
		return []*mysql.Transaction{}
	}

	for i := 0; i < len(b); i++ {
		if withReceipt == 1 {
			receipt, _ := db.GetReceipt(b[i].Hash)
			receipt.BlockNumber = fmt.Sprintf("%d", mysql.SanitizeUint32(receipt.BlockNumber))
			b[i].Receipt = string(receipt.String())
		}
		b[i].BlockNumber = fmt.Sprintf("%d", mysql.SanitizeUint32(b[i].BlockNumber))
	}

	return b
}

func (api *WalletAPI) TransactionEthHistory(address string, page int, withReceipt int) []*mysql.Transaction {
	pageFrom, pageTo := api.parsePage(page)
	log.Info("xxx", "from", pageFrom, "to", pageTo)
	b, err := db.GetTransactionEthHistory(address, pageFrom, pageTo)
	if err != nil {
		log.Error("TransactionEthHistory fails.", "address", address, "err", err)
		return []*mysql.Transaction{}
	}

	for i := 0; i < len(b); i++ {
		if withReceipt == 1 {
			receipt, _ := db.GetReceipt(b[i].Hash)
			receipt.BlockNumber = fmt.Sprintf("%d", mysql.SanitizeUint32(receipt.BlockNumber))
			b[i].Receipt = string(receipt.String())
		}
		b[i].BlockNumber = fmt.Sprintf("%d", mysql.SanitizeUint32(b[i].BlockNumber))
	}

	return b
}

func (api *WalletAPI) TransactionContractHistory(address, contractAddress string, page int, withReceipt int) []*mysql.Transaction {
	pageFrom, pageTo := api.parsePage(page)
	b, err := db.GetTransactionErc20History(address, contractAddress, pageFrom, pageTo)
	if err != nil {
		log.Error("TransactionContractHistory fails.", "address", address, "contract", contractAddress, "err", err)
		return []*mysql.Transaction{}
	}

	for i := 0; i < len(b); i++ {
		if withReceipt == 1 {
			receipt, _ := db.GetReceipt(b[i].Hash)
			receipt.BlockNumber = fmt.Sprintf("%d", mysql.SanitizeUint32(receipt.BlockNumber))
			b[i].Receipt = string(receipt.String())
		}
		b[i].BlockNumber = fmt.Sprintf("%d", mysql.SanitizeUint32(b[i].BlockNumber))
	}

	return b
}

func (api *WalletAPI) Tokens(address string) []*mysql.Contract {
	b, err := db.GetUserTokens(address)
	if err != nil {
		log.Error("get user Tokens fails.", "address", address, "err", err)
		return []*mysql.Contract{}
	}
	return b
}
