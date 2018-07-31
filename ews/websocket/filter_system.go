package websocket

const (
	// UnknownSubscription indicates an unknown subscription type
	UnknownSubscription Type = iota

	// 发现新token
	NewErc20

	// LastSubscription keeps track of the last index
	LastIndexSubscription
)

// Type determines the kind of filter and is used to put the filter in to
// the correct bucket when added.
type Type byte
