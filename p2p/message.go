package p2p

import "net"

// RPC holds any arbitrary data that is being sent over the
// each transpor between two nodes in the network
type RPC struct {
	From    net.Addr
	Payload []byte
}
