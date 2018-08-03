package mysql

import (
	ISql "database/sql"
	"errors"
	"fmt"
	"github.com/ChungkueiBlock/kueiWalletService/internal/log"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"strings"
)

var (
	TableNameBlock              = "block"
	TableNameTransaction        = "transaction"
	TableNameTransactionReceipt = "transaction_receipt"
	TableNameContract           = "erc20_contract"
	TableNameUserTokens         = "user_tokens"
)

var (
	// mysql全局唯一句柄
	db  *ISql.DB
	err error
)

type Connection struct {
	err error
	db  *ISql.DB
	tx  *ISql.Tx
}

// sql.Open函数实际上是返回一个连接池对象，不是单个连接。在open的时候并没有去连接数据库，只有在执行query、exce方法的时候才会去实际连接数据库
func connect(dsn string) *ISql.DB {
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	db, err = ISql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)

	log.Info("initiate mysql successfully!")
	return db
}

func NewConnection(dsn string) *Connection {
	if dsn == "" {
		log.Error("Mysql dsn is empty!\n")
		os.Exit(1)
	}

	conn := &Connection{}
	if db != nil {
		conn.db = db
	} else {
		conn.db = connect(dsn)
		_, err := conn.MaxBlockNumber()
		if err != nil {
			panic(err)
		}
	}

	conn.err = nil
	conn.tx = nil
	return conn
}
func (conn *Connection) Close() {
	if conn != nil && conn.db != nil {
		conn.db.Close()
	}
}
func (conn *Connection) Error() error {
	return conn.err
}
func (conn *Connection) hasError() bool {
	if conn.err != nil {
		return true
	}
	if conn.tx == nil {
		conn.err = errors.New("conn.tx is empty")
		return true
	}
	return false
}

func logSql(format string, a ...interface{}) {
	for _, v := range a {
		switch v.(type) {
		case string:
			format = strings.Replace(format, "?", "\"%v\"", 1)
		default:
			format = strings.Replace(format, "?", "%v", 1)
		}
	}

	log.Trace(fmt.Sprintf(format, a...))
}

/*************
 *  helpers
 ************/
func SanitizeUint32(hex string) uint32 {
	hex = hex[2:]
	value, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		panic("invalid uint32 hexstring:" + hex + "\n" + err.Error())
	}
	return uint32(value)
}

func sanitizeNonce(hex string) string {
	if len(hex) != 18 {
		panic(fmt.Sprintf("invalid nonce with length:%d:%s", len(hex), hex))
	}
	return hex
}

func sanitizeAddress(hex string) string {
	if len(hex) != 42 {
		panic(fmt.Sprintf("invalid address with length:%d:%s", len(hex), hex))
	}
	return hex
}

func sanitizeHash(hex string) string {
	if len(hex) != 66 {
		panic(fmt.Sprintf("invalid hash with length:%d:%s", len(hex), hex))
	}
	return hex
}
