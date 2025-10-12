package p2p

// Peeer is an interface that represents the remote node.
type Peer interface {
	Close() error
}

// Transport is anyting that handles the communication
// between the nodes in the network. This cna be of the
// form (TCP, UDP, websockets, ...)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
