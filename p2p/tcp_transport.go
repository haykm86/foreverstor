package p2p

import (
	"fmt"
	"net"
)

// TCPPeer represents the remote node ocer a TCP established connection
type TCPPeer struct {
	// conn is the underlying connection of the peer
	conn net.Conn
	// if we dial and retrive a conn => oubound == true
	// if we accept and retrive a conn => oubound == false
	outbound bool
}

type TCPTransportOps struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOps
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func NewTCPTransport(ops TCPTransportOps) *TCPTransport {
	return &TCPTransport{
		TCPTransportOps: ops,
		rpcch:           make(chan RPC),
	}
}

// Consume implements the Transport interface, wich will return read-only channel
// for reading the incoming messages received from another peer.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		fmt.Printf("new incoming connection %+v\n", conn)

		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error
	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read Loop
	rpc := RPC{}
	for {
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			if err == net.ErrClosed {
				// remote closed cleanly; stop handling this connection
				return
			}

			fmt.Printf("TCP error: %s\n", err)
			continue
		}
		// TODO: handle decoded message
		// fmt.Printf("message: %v from:%v\n ", string(rpc.Payload), rpc.From)
		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}
}
