package p2p

import "net"

// Peeer is an interface that represents the remote node.
type Peer interface {
	net.Conn
	Send([]byte) error
	CloseStream()
}

// Transport is anyting that handles the communication
// between the nodes in the network. This cna be of the
// form (TCP, UDP, websockets, ...)
type Transport interface {
	Addr() string
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
