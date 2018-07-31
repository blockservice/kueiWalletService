package mysql

import (
	ISql "database/sql"
	"encoding/json"
)

// logs https://codeburst.io/deep-dive-into-ethereum-logs-a8d2047c7371
type TransactionReceipt struct {
	Id                uint32      `json:"id"`
	BlockHash         string      `json:"block_hash"`
	BlockNumber       string      `json:"block_number"`
	ContractAddress   string      `json:"contract_address"`
	CumulativeGasUsed string      `json:"-"`
	From              string      `json:"from"`
	To                string      `json:"to"`
	GasUsed           string      `json:"gas_used"`
	TransactionHash   string      `json:"transaction_hash"`
	TransactionIndex  string      `json:"transaction_index"`
	Status            string      `json:"status"`
	LogsBloom         string      `json:"-"`
	Logs              interface{} `json:"-"`
}

func (receipt *TransactionReceipt) String() []byte {
	if receipt == nil {
		return nil
	}
	res, _ := json.Marshal(receipt)
	return res
}

func (conn *Connection) GetReceipt(txhash string) (*TransactionReceipt, error) {
	fields := "`id`,`blockHash`, `blockNumber`,`contractAddress`,`from`,`to`,`gasUsed`,`transactionHash`,`transactionIndex`,`status`"
	sql := "SELECT " + fields + " FROM `" + TableNameTransactionReceipt + "` WHERE `transactionHash`=? LIMIT 1"
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
	rows, err := stm.Query(txhash)
	if err != nil {
		return nil, err
	}

	results := make([]*TransactionReceipt, 0)
	for rows.Next() {
		d := &TransactionReceipt{}
		rows.Scan(&d.Id, &d.BlockHash, &d.BlockNumber, &d.ContractAddress, &d.From, &d.To, &d.GasUsed, &d.TransactionHash, &d.TransactionIndex, &d.Status)
		results = append(results, d)
	}

	logSql(sql, txhash)
	if len(results) > 0 {
		return results[0], nil
	} else {
		return nil, nil
	}
}
