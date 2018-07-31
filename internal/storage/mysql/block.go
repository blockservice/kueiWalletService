package mysql

import (
	ISql "database/sql"
)

type Block struct {
	Number           string
	Hash             string
	Size             string
	Timestamp        string
	GasLimit         string
	GasUsed          string
	Difficulty       string
	TotalDifficulty  string
	Miner            string
	ParentHash       string
	MixHash          string
	ReceiptsRoot     string
	Sha3Uncles       string
	StateRoot        string
	TransactionsRoot string
	Nonce            string
	ExtraData        string
	LogsBloom        string
	Transactions     []*Transaction
	Uncles           interface{}
}

func (conn *Connection) MaxBlockNumber() (uint64, error) {
	db := conn.db
	sql := "SELECT max(`block_number`) as 'number' FROM `" + TableNameBlock + "`"
	stm, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer func(stm *ISql.Stmt) {
		if stm != nil {
			stm.Close()
		}
	}(stm)
	rows, err := stm.Query()
	if err != nil {
		return 0, err
	}

	var d uint64
	for rows.Next() {
		rows.Scan(&d)
	}

	logSql(sql)
	return d, nil
}
