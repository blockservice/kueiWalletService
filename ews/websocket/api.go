package websocket

import (
	"context"
	"github.com/ChungkueiBlock/kueiWalletService/internal/log"
	"github.com/ChungkueiBlock/kueiWalletService/internal/rpc"
	"github.com/nsqio/go-nsq"
	"sync"
	"time"
)

var (
	deadline = 5 * time.Minute // consider a filter inactive if it has not been polled for within deadline
)

type subscription struct {
	id       rpc.ID
	typ      Type
	topic    string
	deadline *time.Timer   // filter is inactiv when deadline triggers
	erc20s   chan string   // new token json string
	consumer *nsq.Consumer // nsq consumer
}

func (sub *subscription) Unsubscribe() {
	if sub != nil && sub.consumer != nil {
		sub.consumer.Stop()
	}
}

type PubSubAPI struct {
	log           *log.Logger
	filtersMu     sync.Mutex
	nsqHost       string
	nsqInterval   time.Duration
	subscriptions map[rpc.ID]*subscription
}

func NewPubSubAPI(nsqHost string, nsqInterval time.Duration, log *log.Logger) *PubSubAPI {
	api := &PubSubAPI{
		log:           log,
		nsqHost:       nsqHost,
		nsqInterval:   nsqInterval,
		subscriptions: make(map[rpc.ID]*subscription),
	}
	go api.timeoutLoop()

	return api
}

// timeoutLoop runs every 5 minutes and deletes filters that have not been recently used.
// Tt is started when the api is created.
func (pubsub *PubSubAPI) timeoutLoop() {
	ticker := time.NewTicker(deadline)
	for {
		<-ticker.C
		pubsub.filtersMu.Lock()
		for id, f := range pubsub.subscriptions {
			select {
			case <-f.deadline.C:
				log.Info("subscription closed", "rpcid", f.id, "typ", f.typ)
				f.Unsubscribe()
				delete(pubsub.subscriptions, id)
			default:
				continue
			}
		}
		pubsub.filtersMu.Unlock()
	}
}

type logger2 struct {
	logger *log.Logger
}

func (l *logger2) Output(calldepth int, s string) error {
	(*(l.logger)).Info(s)
	return nil
}

type ConsumerT struct {
	erc20s chan<- string
}

func (c *ConsumerT) HandleMessage(msg *nsq.Message) error {
	log.Info("NSQ RECV", "receive", msg.NSQDAddress, "message:", string(msg.Body))
	c.erc20s <- string(msg.Body)
	return nil
}

func (pubsub *PubSubAPI) createNsqConsumer(rpcid rpc.ID, erc20s chan string, address string) (*nsq.Consumer, error) {
	topic := address
	channel := "test-channel"
	host := pubsub.nsqHost
	nsqLookupInterval := time.Second * pubsub.nsqInterval
	log.Info("new nsq", "host", host, "interval", nsqLookupInterval)

	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = nsqLookupInterval
	c, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		return nil, err
	}

	log2 := &logger2{logger: pubsub.log}
	c.SetLogger(log2, 0)

	c.AddHandler(&ConsumerT{erc20s})

	if err := c.ConnectToNSQLookupd(host); err != nil {
		return nil, err
	}

	return c, nil
}

func (pubsub *PubSubAPI) NewToken(ctx context.Context, address string) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	erc20s := make(chan string)
	rpcSub := notifier.CreateSubscription()
	c, err := pubsub.createNsqConsumer(rpcSub.ID, erc20s, address)
	if err != nil {
		return &rpc.Subscription{}, err
	}

	go func() {
		subid := rpcSub.ID
		pubsub.subscriptions[subid] = &subscription{
			id:       subid,
			typ:      NewErc20,
			deadline: time.NewTimer(deadline),
			topic:    address,
			erc20s:   erc20s,
		}
		pubsub.subscriptions[subid].consumer = c

		for {
			select {
			case h := <-erc20s:
				notifier.Notify(rpcSub.ID, h)
			case <-rpcSub.Err():
				pubsub.subscriptions[subid].Unsubscribe()
				delete(pubsub.subscriptions, subid)
				return
			case <-notifier.Closed():
				pubsub.subscriptions[subid].Unsubscribe()
				delete(pubsub.subscriptions, subid)
				return
			}
		}
	}()

	return rpcSub, nil
}

func (pubsub *PubSubAPI) PingPong(ctx context.Context, targetRpcId rpc.ID) (*rpc.Subscription, error) {
	var rpcid rpc.ID
	for id, f := range pubsub.subscriptions {
		if targetRpcId == id {
			f.deadline.Reset(deadline)
			rpcid = id
			break
		}
	}

	if rpcid == "" {
		log.Debug("heartbeat notFound", "no match", targetRpcId)
	} else {
		log.Debug("heartbeat received", "rpcid", targetRpcId)
	}
	return &rpc.Subscription{ID: rpcid}, nil
}
