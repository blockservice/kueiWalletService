package main

import (
	"fmt"
	"github.com/ChungkueiBlock/kueiWalletService/cmd/utils"
	"github.com/ChungkueiBlock/kueiWalletService/internal/debug"
	"github.com/ChungkueiBlock/kueiWalletService/internal/log"
	"github.com/ChungkueiBlock/kueiWalletService/internal/node"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
	"os"
	"runtime"
	"sort"
)

const (
	clientIdentifier = "chungkueiwalletservice" // Client identifier to advertise over the network
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	// The app that holds all commands and flags.
	app = utils.NewApp(gitCommit, "the Chungkuei command line interface")
	// flags that configure the node
	nodeFlags = []cli.Flag{
		utils.IdentityFlag,
		utils.RPCCORSDomainFlag,
		utils.RPCVirtualHostsFlag,
		configFileFlag,
	}

	rpcFlags = []cli.Flag{
		utils.RPCEnabledFlag,
		utils.RPCListenAddrFlag,
		utils.RPCPortFlag,
		utils.RPCApiFlag,
		utils.WSEnabledFlag,
		utils.WSListenAddrFlag,
		utils.WSPortFlag,
		utils.WSApiFlag,
		utils.WSAllowedOriginsFlag,
		utils.Ethgasstation,
	}

	nsqFlags = []cli.Flag{
		utils.NSQNslookupHostFlag,
		utils.NSQNslookupIntervalFlag,
	}
)

func init() {
	// Initialize the CLI app and start Geth
	app.Action = ewsAction
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2018-2020 Chungkuei Team"
	app.Commands = []cli.Command{
		versionCommand,
		licenseCommand,
		// See config.go
		dumpConfigCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Flags = append(app.Flags, nodeFlags...)
	app.Flags = append(app.Flags, rpcFlags...)
	app.Flags = append(app.Flags, nsqFlags...)
	app.Flags = append(app.Flags, debug.Flags...)

	app.Before = func(ctx *cli.Context) error {
		altsrc.InitInputSourceWithContext(app.Flags, altsrc.NewTomlSourceFromFlagFunc("conf"))(ctx)
		runtime.GOMAXPROCS(runtime.NumCPU())
		if err := debug.Setup(ctx); err != nil {
			return err
		}
		// Start system runtime metrics collection
		//go metrics.CollectProcessMetrics(3 * time.Second)

		return nil
	}

	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// geth is the main entry point into the system if no special subcommand is ran.
// It creates a default node based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
func ewsAction(ctx *cli.Context) error {
	ewsNode := makeFullNode(ctx)
	startNode(ctx, ewsNode)
	log.Info("chungkueiwalletservice started...")
	ewsNode.Wait()
	return nil
}

// startNode boots up the system node and all registered protocols, after which
// it unlocks any requested accounts, and starts the RPC/IPC interfaces and the
// miner.
func startNode(ctx *cli.Context, stack *node.Node) {
	//debug.Memsize.Add("node", stack)

	// Start up the node itself
	utils.StartNode(stack)
}
