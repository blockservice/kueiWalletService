package mysql

import (
	ISql "database/sql"
)

type Contract struct {
	Bytecode    string `json:"bytecode,omitempty"`
	Address     string `json:"address"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimals    int    `json:"decimals"`
	TotalSupply string `json:"total_supply,omitempty"`
}

func (conn *Connection) GetUserTokens(address string) ([]*Contract, error) {
	if address == "" {
		return []*Contract{}, nil
	}

	db := conn.db
	sql := "SELECT t.`address`,t.`name`,t.`symbol`,t.`decimals` FROM `" + TableNameContract + "` as t,`" + TableNameUserTokens + "` as c "
	sql = sql + " WHERE t.`address`=c.contract_address AND c.`user_address`=?"
	stm, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer func(stm *ISql.Stmt) {
		if stm != nil {
			stm.Close()
		}
	}(stm)
	rows, err := stm.Query(address)
	if err != nil {
		return nil, err
	}

	results := make([]*Contract, 0)
	for rows.Next() {
		d := &Contract{}
		rows.Scan(&d.Address, &d.Name, &d.Symbol, &d.Decimals)
		results = append(results, d)
	}

	logSql(sql, address)
	return results, nil
}
