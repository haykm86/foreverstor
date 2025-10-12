package main

import (
	"fmt"
	"log"

	"github.com/haykm86/foreverstor/p2p"
)

func OnPeer(p2p.Peer) error {
	fmt.Println("doing some logic with the peer outside the tcp transport")
	return nil
}

func main() {
	tcpOpts := p2p.TCPTransportOps{
		ListenAddr:    ":44040",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("Messag: %+v\n", msg)
			//fmt.Printf("Message:%v\n", string(msg.Payload))
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
