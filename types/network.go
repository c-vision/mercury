package types

import (
	"fmt"
)

type Chain uint8 

const (
	Bitcoin Chain = 0
	Ethereum Chain = 1
	ZCash Chain = 2
)

// String implements the `Stringer` interface.
func (chain Chain) String() string {
	switch chain {
	case Bitcoin:
		return "bitcoin"
	case Ethereum:
		return "ethereum"
	case ZCash:
		return "zcash"
	default:
		panic(ErrUnknownChain)
	}
}

// Network of the blockchain.
type Network interface {
	fmt.Stringer
	Chain() Chain
}