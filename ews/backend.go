package ews

import (
	"github.com/ChungkueiBlock/ecoinWalletService/ews/websocket"
	"github.com/ChungkueiBlock/ecoinWalletService/internal/log"
	"github.com/ChungkueiBlock/ecoinWalletService/internal/node"
	"github.com/ChungkueiBlock/ecoinWalletService/internal/rpc"
)

type ChungkueiWalletService struct {
	ctx    *node.ServiceContext
	config *Config
	log    *log.Logger
}

// New creates a new Ethereum object (including the
// initialisation of the common Ethereum object)
func New(ctx *node.ServiceContext, config *Config, log *log.Logger) (*ChungkueiWalletService, error) {

	instance := &ChungkueiWalletService{
		ctx:    ctx,
		config: config,
		log:    log,
	}
	return instance, nil
}

// APIs returns the collection of RPC services the ethereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *ChungkueiWalletService) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "ews",
			Version:   "1.0",
			Service:   NewWalletAPI(s),
			Public:    true,
		},
		{
			Namespace: "ews",
			Version:   "1.0",
			Service:   websocket.NewPubSubAPI(s.config.NSQNslookupHost, s.config.NSQNslookupInterval, s.log),
			Public:    true,
		},
	}
}

func (s *ChungkueiWalletService) Stop() error {
	return nil
}
