package ews

import "time"

const (
	DefaultMysqlDsn            = "root:123456@tcp(localhost:3306)/ethereum_blockchain_v1?charset=utf8"
	DefaultNSQNslookupHost     = "127.0.0.1:4161"
	DefaultNSQNslookupInterval = 60
	DefaultPageSize            = 10
)

var DefaultConfig = Config{
	MysqlDSN:   DefaultMysqlDsn,
	TxPageSize: DefaultPageSize,
}

type Config struct {
	MysqlDSN            string
	TxPageSize          int
	EnableGasStation    bool
	NSQNslookupHost     string
	NSQNslookupInterval time.Duration
}
