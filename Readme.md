# kueiWalletService
kueiWalletService给钱包提供冗余的历史数据接口,以及一些辅助功能:
+ Token列表
+ 交易历史
+ 当前全网平均gasprice(交易时的推荐gasprice)
+ 自定义gasprice时, 估算交易落块时间
+ fiat, 代币-法币兑换率

![walletService](https://raw.githubusercontent.com/ChungkueiBlock/kueiWalletService/docs/docs/images/walletService.png)

### Usage:
```
VERSION:
   1.8.11-unstable

COMMANDS:
   dumpconfig  Show configuration values
   license     Display license information
   version     Print version numbers
   help        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --identity value             Custom node name
   --rpccorsdomain value        Comma separated list of domains from which to accept cross origin requests (browser enforced)
   --rpcvhosts value            Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard. (default: "localhost")
   --conf value                 TOML configuration file
   --rpc                        Enable the HTTP-RPC server
   --rpcaddr value              HTTP-RPC server listening interface (default: "localhost")
   --rpcport value              HTTP-RPC server listening port (default: 8545)
   --rpcapi value               API's offered over the HTTP-RPC interface
   --ws                         Enable the WS-RPC server
   --wsaddr value               WS-RPC server listening interface (default: "localhost")
   --wsport value               WS-RPC server listening port (default: 8546)
   --wsapi value                API's offered over the WS-RPC interface
   --wsorigins value            Origins from which to accept websockets requests
   --ethgasstation              Enable crawling https://ethgasstation.info/json
   --nsqnslookup value          nsq nsqlookupd host (default: "127.0.0.1:4161")
   --nsqnslookupinterval value  nsqLookupInterval x seconds (default: 60)
   --verbosity value            Logging verbosity: 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=detail (default: 3)
   --vmodule value              Per-module verbosity: comma-separated list of <pattern>=<level> (e.g. eth/*=5,p2p=4)
   --backtrace value            Request a stack trace at a specific logging statement (e.g. "block.go:271")
   --debug                      Prepends log messages with call-site location (file and line number)
   --help, -h                   show help
```

#### APIs
1. ews_estimateGasprice {level} //预估gasprice, level取值"fast", "middel", "slow"
    ```bash
    curl -s -X POST \
      -H "Content-Type: application/json" \
      --data '{"jsonrpc": "2.0", "id": 1, "method": "ews_estimateGasprice", "params": ["middle"]}' \
      ${server}

    // results:
    {
      "jsonrpc": "2.0",
      "id": 1,
      "result": {
        "txMeanSecs": 234,
        "gasprice": 30 // Gwei
      }
    }
    ```

2. ews_estimateTxTime {gasPrice} {gasLimit} // 根据gasPrice估算交易打包时间
    ```bash
    curl -s -X POST \
      -H "Content-Type: application/json" \
      --data '{"jsonrpc": "2.0", "id": 1, "method": "ews_estimateTxTime", "params": [10, 21000]}' \
      ${server}

    // results:
    {
      "jsonrpc": "2.0",
      "id": 1,
      "result": {
        "txMeanSecs": 40.830126860792255,
        "minedProb": "Very High"
      }
    }
    ```
3. ews_transactionHistory {address} {page} // 交易历史
    > ews_transactionEthHistory {address} {page} // ETH 交易历史
    > ews_transactionContractHistory {address} {token_address} {page} // token 交易历史

    ```bash
    curl -s -X POST \
      -H "Content-Type: application/json" \
      --data '{"jsonrpc": "2.0", "id": 1, "method": "ews_transactionEthHistory", "params": ["'${address}'", 0, 1]}' \
      ${server}

    // results:
    {
      "jsonrpc": "2.0",
      "id": 1,
      "result": [
        {
          "id": 26,
          "hash": "0x6b32d600e27c64be45672c4f505fd1bd67431735e83a28c83a5c8add93b3010f",
          "blockNumber": "3009390",
          "blockHash": "0xe1472635fcb32113d8ce37768ae0ae03d48f4791123810ad9e94a4f6ab6677f4",
          "from": "0x81b7e08f65bdf5648606c89998a9cc8164397647",
          "to": "0xf527d95ca4537af3ce81302e3870e87a7eaca185",
          "receipt": "",
          "gas": "0x5208",
          "gasPrice": "0x55ae82600",
          "nonce": "0x4be45d",
          "transactionIndex": "0x19",
          "value": "0xde0b6b3a7640000",
          "v": "0x1b",
          "r": "0x2af763ef546cebcc66c560c0dd4356c60cd19067278074e92327806cc6ef8f8f",
          "s": "0x7011d0f09722408fd2f106d7467f91ed07973f18cead85a4f6840a17f0142a1b",
          "input": "0x",
          "is_contract_tx": 0
        },
        ...
      ]
    }
    ```

4. ews_tokens {address} //获取用户的token列表
    ```bash
    curl -s -X POST \
      -H "Content-Type: application/json" \
      --data '{"jsonrpc": "2.0", "id": 1, "method": "ews_tokens", "params": ["'${address}'"]}' \
      ${server}

    // results
    {
      "jsonrpc": "2.0",
      "id": 1,
      "result": [
        {
          "address": "0xa63bc6ac2fbc95b83650464314ea37858f0d6944",
          "name": "XGuppy",
          "symbol": "XGUP",
          "decimals": 3
        }
      ]
    }
    ```

5. websockets推送(experimental)