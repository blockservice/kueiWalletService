package mysql

import (
	ISql "database/sql"
	"fmt"
)

type Transaction struct {
	Id               uint32 `json:"id"`
	Hash             string `json:"hash"`
	BlockNumber      string `json:"blockNumber"`
	BlockHash        string `json:"blockHash"`
	From             string `json:"from"`
	To               string `json:"to"`
	Receipt          string `json:"receipt"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Nonce            string `json:"nonce"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
	Input            string `json:"input"`
	IsContract       int    `json:"is_contract_tx"`
	TimeStamp        string `json:"time_stamp"`
}

//  所有交易历史
func (conn *Connection) GetTransactionHistory(address string, pageFrom, pageTo int) ([]*Transaction, error) {
	sql := "SELECT `id`,`hash`,`blockNumber`,`blockHash`,`from`,`to`,`receipt`,`gas`,`gasPrice`,`nonce`,`transactionIndex`,`value`,`v`,`r`,`s`,`input`,UNIX_TIMESTAMP(`createdAt`) as 'timeStamp' FROM `" + TableNameTransaction + "` WHERE (`from`=?) OR (`to`=?) OR (`receipt`=?)"
	sql = sql + fmt.Sprintf(" ORDER BY `id` DESC LIMIT %d,%d", pageFrom, pageTo)
	db := conn.db
	stm, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer func(stm *ISql.Stmt) {
		if stm != nil {
			stm.Close()
		}
	}(stm)
	rows, err := stm.Query(address, address, address)
	if err != nil {
		return nil, err
	}

	results := make([]*Transaction, 0)
	for rows.Next() {
		d := &Transaction{}
		rows.Scan(&d.Id, &d.Hash, &d.BlockNumber, &d.BlockHash, &d.From, &d.To, &d.Receipt, &d.Gas, &d.GasPrice, &d.Nonce, &d.TransactionIndex, &d.Value, &d.V, &d.R, &d.S, &d.Input, &d.TimeStamp)
		results = append(results, d)
	}

	logSql(sql, address, address, address)
	return results, nil
}

// 仅EOA交易历史
func (conn *Connection) GetTransactionEthHistory(address string, pageFrom, pageTo int) ([]*Transaction, error) {
	db := conn.db
	sql := "SELECT `id`,`hash`,`blockNumber`,`blockHash`,`from`,`to`,`receipt`,`gas`,`gasPrice`,`nonce`,`transactionIndex`,`value`,`v`,`r`,`s`,`input`,UNIX_TIMESTAMP(`createdAt`) as 'timeStamp' FROM `" + TableNameTransaction + "` WHERE `is_contract_tx`=0 AND (`from`=? OR `to`=?)"
	sql = sql + fmt.Sprintf(" ORDER BY `id` DESC LIMIT %d,%d", pageFrom, pageTo)
	stm, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer func(stm *ISql.Stmt) {
		if stm != nil {
			stm.Close()
		}
	}(stm)
	rows, err := stm.Query(address, address)
	if err != nil {
		return nil, err
	}

	results := make([]*Transaction, 0)
	for rows.Next() {
		d := &Transaction{}
		rows.Scan(&d.Id, &d.Hash, &d.BlockNumber, &d.BlockHash, &d.From, &d.To, &d.Receipt, &d.Gas, &d.GasPrice, &d.Nonce, &d.TransactionIndex, &d.Value, &d.V, &d.R, &d.S, &d.Input, &d.TimeStamp)
		results = append(results, d)
	}

	logSql(sql, address, address)
	return results, nil
}

// 仅ERC20交易历史
func (conn *Connection) GetTransactionErc20History(address string, contractAddress string, pageFrom, pageTo int) ([]*Transaction, error) {
	db := conn.db
	sql := "SELECT `id`, `hash`,`blockNumber`,`blockHash`,`from`,`to`,`receipt`,`gas`,`gasPrice`,`nonce`,`transactionIndex`,`value`,`v`,`r`,`s`,`input`,UNIX_TIMESTAMP(`createdAt`) as 'timeStamp' FROM `" + TableNameTransaction + "` "
	sql = sql + " WHERE `is_contract_tx`=1 AND (`from`=? AND `to`=?) OR (`receipt`=? AND `to`=?)"
	sql = sql + fmt.Sprintf(" ORDER BY `id` DESC LIMIT %d,%d", pageFrom, pageTo)
	stm, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer func(stm *ISql.Stmt) {
		if stm != nil {
			stm.Close()
		}
	}(stm)
	rows, err := stm.Query(address, contractAddress, address, contractAddress)
	if err != nil {
		return nil, err
	}

	results := make([]*Transaction, 0)
	for rows.Next() {
		d := &Transaction{}
		rows.Scan(&d.Id, &d.Hash, &d.BlockNumber, &d.BlockHash, &d.From, &d.To, &d.Receipt, &d.Gas, &d.GasPrice, &d.Nonce, &d.TransactionIndex, &d.Value, &d.V, &d.R, &d.S, &d.Input, &d.TimeStamp)
		results = append(results, d)
	}

	logSql(sql, address, address, contractAddress)
	return results, nil
}
