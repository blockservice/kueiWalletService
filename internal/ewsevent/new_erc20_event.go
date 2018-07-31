package ewsevent

import "github.com/ethereum/go-ethereum/common"

type NewErc20 struct {
	Address  *common.Address `json:"address"`
	Name     string          `json:"name"`
	Symbol   string          `json:"symbol"`
	Decimals int             `json:"decimals"`
}
