package p2p

import (
	"errors"
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

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

// RemoteAddr implements the Peer interface and will return the
// remote address of its underlying connection.
func (p *TCPPeer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

// Consume implements the Transport interface, wich will return read-only channel
// for reading the incoming messages received from another peer.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

// Dial implements the Transport interface.
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)
	return nil
}

func NewTCPTransport(ops TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: ops,
		rpcch:            make(chan RPC),
	}
}

// Close implements the Transport interface.
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	fmt.Printf("TCP transport listening on port: %s\n", t.ListenAddr)

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		fmt.Printf("new incoming connection %+v\n", conn)

		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outBound bool) {
	var err error
	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outBound)

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
