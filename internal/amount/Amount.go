package amount

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"strconv"
)

type Amount struct {
	Number *big.Int //wei
}

var Zero = &Amount{Number: big.NewInt(0)}

func NewAmount(i interface{}) *Amount {
	var c Amount
	switch i.(type) {
	case uint64:
		ii := new(big.Int).SetUint64(i.(uint64))
		c = Amount{Number: ii}
	case *big.Int:
		c = Amount{Number: i.(*big.Int)}
	}

	return &c
}

func FromHex(hex string) (*Amount, error) {
	if hex == "" {
		return NewAmount(0), nil
	}
	v, err := hexutil.DecodeBig(hex)
	if err != nil {
		return nil, err
	}

	c := NewAmount(v)
	return c, nil
}

//BUG:ether浮点数不能超过20位
func FromEther(ether float64) *Amount {
	a := new(big.Int).SetUint64(uint64(ether * 1e18))
	e := NewAmount(a)
	return e
}

func (a *Amount) ToWei() *big.Int {
	return a.Number
}

func (a *Amount) ToHex() string {
	c := hexutil.EncodeBig(a.Number)
	return c
}

func (a *Amount) Mul(b uint64) *Amount {
	d := big.NewInt(0)
	c := a.Number
	d.Mul(c, new(big.Int).SetUint64(b))
	e := Amount{Number: d}
	return &e
}

// ParseEther 返回整数 单位:Wei
func (a *Amount) ParseEther() *big.Int {
	return a.Number
}

// FormatEther 返回浮点数 单位:ETH
func (a *Amount) FormatEther(decimals float64) float64 {
	b := a.Number
	if b == nil {
		return 0.0
	}

	whole := big.NewInt(0)
	fraction := big.NewInt(0)
	decimalsBig := big.NewInt(int64(decimals))
	whole.Div(b, decimalsBig)

	t := big.NewInt(0)
	fraction.Sub(b, t.Mul(whole, decimalsBig))

	wholeInt := whole.Uint64()
	fractionInt := fraction.Uint64()
	result := float64(wholeInt) + float64(fractionInt)/float64(decimals)

	return result
}

// ToEtherStr 返回基于eth的小数字符串 单位:ETH
func (a Amount) ToEtherStr(decimals float64) string {
	b := a.FormatEther(decimals)
	return strconv.FormatFloat(b, 'f', 8, 64)
}

//func (a Amount) Float64() float64 {
//	return a.FormatEther(a.Decimals)
//}
